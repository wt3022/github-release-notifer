package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() Config {
	if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))); err != nil {
		log.Fatalf("環境ファイルの読み込みに失敗しました: %v", err)
	}
	return Config{
		Env:            os.Getenv("GO_ENV"),
		GithubToken:    os.Getenv("GITHUB_TOKEN"),
		ProductionMode: os.Getenv("PRODUCTION_MODE") == "true",
		EmailConfig: EmailConfig{
			SMTPHost:  os.Getenv("SMTP_HOST"),
			SMTPPort:  os.Getenv("SMTP_PORT"),
			FromEmail: os.Getenv("FROM_EMAIL"),
		},
	}
}
