package tasks

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/google/go-github/github"

	"github.com/wt3022/github-release-notifier/internal/db"
	mygithub "github.com/wt3022/github-release-notifier/internal/github"
	"gorm.io/gorm"
)

func WatchRepositoryRelease(dbClient *gorm.DB, githubClient *github.Client) {
	log.Printf("リポジトリの処理を開始します\n")

	var watchRepos []db.WatchRepository

	if err := dbClient.Find(&watchRepos).Error; err != nil {
		log.Fatalf("データベースからのリポジトリの取得に失敗しました: %v", err)
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("タイムゾーンのロードに失敗しました: %v", err)
	}

	for _, repo := range watchRepos {
		newReleases, err := mygithub.FetchReleasesAfter(context.Background(), githubClient, repo.Owner, repo.Name, repo.LastNotificationDate)
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
			err := dbClient.Model(&repo).Update("LastNotificationDate", latestPublishedAt).Error
			if err != nil {
				log.Printf("データベースの更新に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
			}
			

			// 通知を送信
			// プロジェクトの通知設定を取得
			notification := db.Notification{}
			err = dbClient.Model(&repo).Association("Notification").Find(&notification)
			if err != nil {
				log.Printf("通知設定の取得に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
			} else {
				// 一旦プロジェクトIDとリポジトリID、通知設定を出力
				log.Printf("プロジェクトID: %d, リポジトリID: %d, 通知設定: %s", repo.ProjectID, repo.ID, notification.Type)
			}			

		} else {
			fmt.Printf("リポジトリ %s/%s に新しいリリースはありません\n", repo.Owner, repo.Name)
		}
	}

	fmt.Printf("リポジトリの処理が完了しました\n")
}
