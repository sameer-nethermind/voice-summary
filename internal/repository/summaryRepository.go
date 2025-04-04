package repository

import (
	"database/sql"
	"fmt"
)

type RecordingRepository struct {
	DB *sql.DB
}

func NewRecordingRepository(db *sql.DB) *RecordingRepository {
	return &RecordingRepository{DB: db}
}

func (sr *RecordingRepository) CreateRecording(userId string) int {
	rows, err := sr.DB.Query("INSERT INTO recording (user_id) VALUES ($1) RETURNING id", userId)
	fmt.Println("error db = ", err, userId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
	}
	return id
}

func (sr *RecordingRepository) UpdateRecordingUploaded(id int) {
	_, err := sr.DB.Query("UPDATE recording SET uploaded = true WHERE id = $1", id)
	if err != nil {
		panic(err)
	}
}

func (sr *RecordingRepository) GetRecordingById(id int) bool {
	rows, err := sr.DB.Query("SELECT 1 FROM recording WHERE id = $1 AND uploaded = true", id)
	if err != nil {
		return false
	}

	defer rows.Close()

	return rows.Next()
}
