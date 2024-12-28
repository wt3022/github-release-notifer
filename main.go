package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wt3022/github-release-notifier/handlers"
	"github.com/wt3022/github-release-notifier/internal/db"
	"github.com/wt3022/github-release-notifier/internal/env"
	"github.com/wt3022/github-release-notifier/internal/github"
)

func main() {
	config := env.LoadConfig()
	dbClient := db.OpenDB()

	githubClient, err := github.OpenGitHubClient(context.Background(), config.Token)
	if err != nil {
		log.Fatalf("GitHubクライアントの初期化に失敗しました: %v", err)
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		handlers.HomeHandler(c, dbClient, githubClient)
	})

	router.Run()
}
