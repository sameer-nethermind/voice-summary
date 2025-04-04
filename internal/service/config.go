package service

import "github.com/spf13/viper"

type Config struct {
	LlamaAPIKey        string `mapstructure:"llama_api_key"`
	Region             string `mapstructure:"region"`
	Bucket             string `mapstructure:"bucket"`
	Acl                string `mapstructure:"acl"`
	AWSAccessKeyID     string `mapstructure:"aws_access_key_id"`
	AWSSecretAccessKey string `mapstructure:"aws_secret_access_key"`
	LemonFoxAPIKey     string `mapstructure:"lemonfox_api_key"`
}

func LoadConfig() (config *Config, err error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return &Config{}, err
	}
	var c Config
	c.Region = v.GetString("S3_REGION")
	c.Bucket = v.GetString("S3_BUCKET")
	c.Acl = v.GetString("S3_ACL")
	c.AWSAccessKeyID = v.GetString("AWS_ACCESS_KEY_ID")
	c.AWSSecretAccessKey = v.GetString("AWS_SECRET_ACCESS_KEY")
	c.LlamaAPIKey = v.GetString("LLAMA_API_KEY")
	c.LemonFoxAPIKey = v.GetString("LEMONFOX_API_KEY")

	return &c, nil
}
