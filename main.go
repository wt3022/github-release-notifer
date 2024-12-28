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
	projectRouter.GET("/:id", func(ctx *gin.Context) {
		handlers.DetailProject(ctx, dbClient)
	})
	projectRouter.GET("/", func(ctx *gin.Context) {
		handlers.ListProjects(ctx, dbClient)
	})
	projectRouter.POST("/", func(c *gin.Context) {
		handlers.CreateProjects(c, dbClient)
	})
	projectRouter.PATCH("/", func(c *gin.Context) {
		handlers.UpdateProject(c, dbClient)
	})
	projectRouter.DELETE("/:id", func(c *gin.Context) {
		handlers.DeleteProject(c, dbClient)
	})

	/* 監視リポジトリ */
	repositoriesRouter := router.Group("/repositories")
	repositoriesRouter.GET("/:id", func(ctx *gin.Context) {
		handlers.DetailRepository(ctx, dbClient)
	})
	repositoriesRouter.GET("/", func(ctx *gin.Context) {
		handlers.ListRepositories(ctx, dbClient)
	})
	repositoriesRouter.POST("/", func(ctx *gin.Context) {
		handlers.CreateRepository(ctx, dbClient, githubClient)
	})
	repositoriesRouter.DELETE("/:id", func(ctx *gin.Context) {
		handlers.DeleteRepository(ctx, dbClient)
	})

	/* 通知先 */
	notificationRouter := router.Group("/notifications")
	notificationRouter.GET("/:id", func(ctx *gin.Context) {
		handlers.DetailNotification(ctx, dbClient)
	})
	notificationRouter.GET("/", func(ctx *gin.Context) {
		handlers.ListNotifications(ctx, dbClient)
	})
	notificationRouter.POST("/", func(ctx *gin.Context) {
		handlers.CreateNotification(ctx, dbClient)
	})
	notificationRouter.POST("/:id/test_notification", todo)
	notificationRouter.PATCH("/", func(ctx *gin.Context) {
		handlers.UpdateNotification(ctx, dbClient)
	})
	notificationRouter.DELETE("/:id", func(ctx *gin.Context) {
		handlers.DeleteNotification(ctx, dbClient)
	})

	/* 定期タスク実行
	* リポジトリの更新通知
	 */

	router.Run()
}
