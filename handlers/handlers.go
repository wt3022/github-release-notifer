package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"

	"github.com/gin-gonic/gin"
	"github.com/wt3022/github-release-notifier/internal/db"
	mygithub "github.com/wt3022/github-release-notifier/internal/github"
	"gorm.io/gorm"
)

func HomeHandler(c *gin.Context, dbClient *gorm.DB, githubClient *github.Client) {
	watchRepos, err := db.GetWatchRepositories(dbClient)
	if err != nil {
		log.Printf("データベースからのリポジトリ取得に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"エラー": "リポジトリ情報の取得に失敗しました"})
		return
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("タイムゾーンのロードに失敗しました: %v", err)
	}

	for _, repo := range watchRepos {
		newReleases, err := mygithub.FetchReleasesAfter(context.Background(), githubClient, repo.Owner, repo.Name, repo.PublishedAt)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(newReleases) > 0 {
			fmt.Printf("リポジトリ %s/%s の新しいリリース:\n", repo.Owner, repo.Name)
			for _, release := range newReleases {
				publishedAtJST := release.PublishedAt.Time.In(jst)
				fmt.Printf("-----------------------------\n")
				fmt.Printf("リリース名: %s\nタグ: %s\n公開日: %s\n", *release.Name, *release.TagName, publishedAtJST.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("-----------------------------\n")

			// データベースの更新
			latestPublishedAt := newReleases[0].PublishedAt.Time
			for _, release := range newReleases {
				if release.PublishedAt.Time.After(latestPublishedAt) {
					latestPublishedAt = release.PublishedAt.Time
				}
			}
			err := dbClient.Model(&repo).Update("PublishedAt", latestPublishedAt).Error
			if err != nil {
				log.Printf("データベースの更新に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
			}
		} else {
			fmt.Printf("リポジトリ %s/%s に新しいリリースはありません\n", repo.Owner, repo.Name)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"メッセージ": "リポジトリの処理が完了しました",
		"件数":    len(watchRepos),
	})
}
