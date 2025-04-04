package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
)

// LemonFoxResponse models the JSON response from LemonFox
type LemonFoxResponse struct {
	Text string `json:"text"`
}

// UploadAudio handles the "receive file → simultaneously upload to S3 and
// chunk-transcribe with LemonFox → pass the combined transcription to Llama → return result."
func (us *UserService) UploadAudio(c *gin.Context) {
	//-------------------------------------------------------------------
	// 1. Receive file from client
	//-------------------------------------------------------------------
	fileHeader, err := c.FormFile("file")
	// Read the userId from the body
	userId := c.PostForm("user_id")
	// request header should be "Content-Type: multipart/form-data"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found in the request"})
		return
	}

	// Additional check: file sanitization - size check, format check etc.
	if !sanitationChecks(fileHeader) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File did not pass sanity checks"})
		return
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open uploaded file"})
		return
	}
	defer file.Close()

	//-------------------------------------------------------------------
	// Create a TeeReader (fork the stream)
	//-------------------------------------------------------------------
	// We'll simultaneously:
	//  - Upload the file to S3 (reading from teeReader)
	//  - Chunk-transcribe the file (reading from pr)
	//-------------------------------------------------------------------
	pr, pw := io.Pipe()
	teeReader := io.TeeReader(file, pw)

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine #1: Upload to S3
	go func() {
		defer wg.Done()
		defer pw.Close() // ensure the pipe is closed when S3 upload finishes

		// initiate uploading to S3
		id := us.recordingRepo.CreateRecording(userId)
		if _, err := us.uploadToS3(teeReader, userId); err != nil {
			fmt.Printf("Error uploading to S3: %v\n", err)
		}

		// Update the recording as uploaded
		us.recordingRepo.UpdateRecordingUploaded(id)
	}()

	// Goroutine #2: Chunk-based transcription with LemonFox - Transcription would take longer than uploading to S3
	var transcriptionResult string
	var transcriptionErr error
	go func() {
		defer wg.Done()
		defer pr.Close()

		transcriptionResult, transcriptionErr = chunkedTranscription(pr)
	}()

	//-------------------------------------------------------------------
	// 2. Wait for concurrency (S3 + chunked transcription) to finish
	//-------------------------------------------------------------------
	wg.Wait()

	// @TODO: In case there is an error and we want to retry then we would need to download from S3 and do it
	if transcriptionErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Transcription failed: %v", transcriptionErr),
		})
		return
	}

	//-------------------------------------------------------------------
	// 3. Pass the combined transcription to the Llama API
	//-------------------------------------------------------------------
	llamaRespBody, statusCode, err := us.callLlamaAPI(transcriptionResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//-------------------------------------------------------------------
	// 4. Return the Llama API response to the client
	//-------------------------------------------------------------------
	c.Data(statusCode, "application/json", llamaRespBody)
}

// sanitationChecks checks the size of the uploaded file, etc.
func sanitationChecks(fileHeader *multipart.FileHeader) bool {
	// Example: limit file size to 100MB
	return fileHeader.Size <= 1024*1024*100
}

func (us *UserService) uploadToS3(r io.Reader, userId string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(us.config.Region),
	)
	if err != nil {
		log.Println("error:", err)
		return "", err
	}

	rows, err := us.DB.Query("SELECT id FROM recording ORDER BY created_at DESC LIMIT 1")
	if err != nil {
		return "", err
	}

	defer rows.Close()
	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return "", err
		}
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)
	req, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(us.config.Bucket),
		Key:    aws.String(createFileName(id+1, userId)),
	})

	fmt.Println("Presigned URL:", req)

	if err != nil {
		log.Println("error:", err)
		return "", err
	}

	return req.URL, nil
}

func (us *UserService) RetryTranscription(recordingId int, userId string) {
	exists := us.recordingRepo.GetRecordingById(recordingId)
	// Create filename
	if !exists {
		fmt.Println("Recording does not exist")
		return
	}
	filename := fmt.Sprintf("%s-%d", userId, recordingId)
	// download the file from S3 and call transcribe on it
	us.downloadFromS3(filename)
	// transcribe
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file")
		return
	}
	defer file.Close()
	// Transcribe the whole file
	transcription, err := chunkedTranscription(file)
	if err != nil {
		fmt.Println("Error transcribing file")
		return
	}
	// Call Llama API
	llamaRespBody, statusCode, err := us.callLlamaAPI(transcription)
	if err != nil {
		fmt.Println("Error calling Llama API")
		return
	}
	fmt.Println("Llama API response: ", llamaRespBody, statusCode)

}

func chunkedTranscription(r io.Reader) (string, error) {
	const chunkSize = 15 * 1024 * 1024 // 15 MB chunks
	var transcript string

	buffer := make([]byte, chunkSize)
	chunkIndex := 0

	for {
		n, err := io.ReadFull(r, buffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if n > 0 {
				text, txErr := transcribeChunk(buffer[:n], chunkIndex)
				if txErr != nil {
					return "", txErr
				}
				transcript += text + " "
			}
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading chunk: %v", err)
		}

		// We have a full chunk
		text, txErr := transcribeChunk(buffer[:n], chunkIndex)
		if txErr != nil {
			return "", txErr
		}
		transcript += text + " "
		chunkIndex++
	}

	return transcript, nil
}

func transcribeChunk(chunk []byte, chunkIndex int) (string, error) {
	// fmt.Println("chunk", chunkIndex, chunk)
	// TODO: Call your LemonFox transcription API with this chunk:
	// Example pseudo-call:
	//   return callLemonFoxTranscription(chunk)
	//
	// For demonstration, we'll just return placeholder text:
	return fmt.Sprintf("[transcribed-chunk-%d]", chunkIndex), nil
}

func (us *UserService) downloadFromS3(filename string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(us.config.Region),
	)
	if err != nil {
		log.Println("error:", err)
		return
	}

	client := s3.NewFromConfig(cfg)

	downloader := manager.NewDownloader(client)

	// Create a file to store in the local system
	file, err := os.Create(filename)
	if err != nil {
		log.Println("error:", err)
		return
	}

	_, err = downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(us.config.Bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		log.Println("error:", err)
		return
	}
}

func createFileName(id int, userId string) string {
	return fmt.Sprintf("%d-%s", id, userId)
}
