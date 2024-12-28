package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wt3022/github-release-notifier/handlers"
	"github.com/wt3022/github-release-notifier/internal/db"
	"github.com/wt3022/github-release-notifier/internal/env"
	"github.com/wt3022/github-release-notifier/internal/github"
)

func todo(c *gin.Context) {
	c.String(http.StatusOK, "todo")
}

func main() {
	config := env.LoadConfig()
	dbClient := db.OpenDB()

	githubClient, err := github.OpenGitHubClient(context.Background(), config.Token)
	if err != nil {
		log.Fatalf("GitHubクライアントの初期化に失敗しました: %v", err)
	}

	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		handlers.HomeHandler(c, dbClient, githubClient)
	})

	/* TODO: ユーザー周りのAPI定義 */

	/* プロジェクト */
	projectRouter := router.Group("/projects")
	projectRouter.GET("/:id", todo)
	projectRouter.GET("/", todo)
	projectRouter.POST("/", todo)
	projectRouter.PATCH("/", todo)
	projectRouter.DELETE("/", todo)

	/* 監視リポジトリ */
	repositoriesRouter := router.Group("/repositories")
	repositoriesRouter.GET("/:id", todo)
	repositoriesRouter.GET("/", todo)
	repositoriesRouter.POST("/", todo)
	repositoriesRouter.DELETE("/", todo)

	/* 通知 */
	notificationRouter := router.Group("/notifications")
	notificationRouter.GET("/:id", todo)
	notificationRouter.GET("/", todo)
	notificationRouter.POST("/:id/test_notification", todo)
	notificationRouter.POST("/", todo)
	notificationRouter.PATCH("/", todo)
	notificationRouter.DELETE("/", todo)


	/* 定期タスク実行
	* リポジトリの更新通知
	*/


	router.Run()
}
