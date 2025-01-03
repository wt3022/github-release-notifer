package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/wt3022/github-release-notifier/internal/db"
	"gorm.io/gorm"
)

func ListRepositories(c *gin.Context, dbClient *gorm.DB) {
	/* リポジトリ一覧を取得します */
	var repositories []db.WatchRepository

	if err := dbClient.Find(&repositories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repositories)
}

func DetailRepository(c *gin.Context, dbClient *gorm.DB) {
	/* リポジトリの詳細を取得します */
	var repository db.WatchRepository

	id := c.Param("id")
	if err := dbClient.First(&repository, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repository)
}

func CreateRepository(c *gin.Context, dbClient *gorm.DB, githubClient *github.Client) {
	/* リポジトリを作成します */
	var repository db.WatchRepository

	if err := c.ShouldBindJSON(&repository); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーが存在するか確認
	_, _, err := githubClient.Users.Get(c, repository.Owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザーが存在しません"})
		return
	}

	// リポジトリが存在するか確認
	_, _, err = githubClient.Repositories.Get(c, repository.Owner, repository.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リポジトリが存在しません"})
		return
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 最新のリリース日を取得
	// 取得できなかった場合はデフォルトの現在のタイムスタンプが入る
	if repository.WatchType == db.WatchTypeRelease {
		release, _, _ := githubClient.Repositories.GetLatestRelease(c, repository.Owner, repository.Name)
		if release != nil && release.PublishedAt != nil {
			repository.LastPublishedAt = release.PublishedAt.In(jst)
		}
	} else {
		tags, _, _ := githubClient.Repositories.ListTags(c, repository.Owner, repository.Name, nil)
		commit, _, _ := githubClient.Git.GetCommit(c, repository.Owner, repository.Name, *tags[0].Commit.SHA)
		if commit != nil && commit.Author.Date != nil {
			repository.LastPublishedAt = commit.Author.Date.In(jst)
		}
	}

	if err := dbClient.Create(&repository).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, repository)
}

func DeleteRepository(c *gin.Context, dbClient *gorm.DB) {
	/* リポジトリを削除します */
	var repository db.WatchRepository

	id := c.Param("id")
	if err := dbClient.Delete(&repository, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
