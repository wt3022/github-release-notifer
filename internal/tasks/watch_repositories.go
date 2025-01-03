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
		processRepository(dbClient, githubClient, repo, jst)
	}

	fmt.Printf("リポジトリの処理が完了しました\n")
}

func processRepository(dbClient *gorm.DB, githubClient *github.Client, repo db.WatchRepository, jst *time.Location) {
	// リリース情報の取得

	if repo.WatchType == db.WatchTypeRelease {
		newReleases, err := mygithub.FetchReleasesAfter(context.Background(), githubClient, repo.Owner, repo.Name, repo.LastNotificationDate)
		if err != nil {
			log.Println(err)
			return
		}

		if len(newReleases) > 0 {
			printNewReleases(repo, newReleases, jst)
			updateLastPublishedAt(dbClient, newReleases[0], repo, jst)
			sendNotification(dbClient, repo)
		} else {
			fmt.Printf("リポジトリ %s/%s に新しいリリースはありません\n", repo.Owner, repo.Name)
		}
	} else if repo.WatchType == db.WatchTypeTag {
		// タグ情報の取得
		tagRelease, err := mygithub.FetchTagReleaseAfter(context.Background(), githubClient, repo.Owner, repo.Name, repo.LastNotificationDate)
		if err != nil {
			log.Println(err)
			return
		}

		if len(tagRelease) > 0 {
			printNewTagRelease(repo, tagRelease, jst)
			updateLastPublishedAtForTag(dbClient, tagRelease[0], repo, jst)
			sendNotification(dbClient, repo)
		} else {
			fmt.Printf("リポジトリ %s/%s に新しいタグリリースはありません\n", repo.Owner, repo.Name)
		}
	}

	// 最終確認日時の更新
	updateLastNotificationDate(dbClient, repo)
}

func printNewReleases(repo db.WatchRepository, newReleases []*github.RepositoryRelease, jst *time.Location) {
	fmt.Printf("リポジトリ %s/%s の新しいリリース:\n", repo.Owner, repo.Name)
	for _, release := range newReleases {
		publishedAtJST := release.PublishedAt.Time.In(jst)
		fmt.Printf("-----------------------------\n")
		fmt.Printf("リリース名: %s\nタグ: %s\n公開日: %s\n", *release.Name, *release.TagName, publishedAtJST.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("-----------------------------\n")
}

func printNewTagRelease(repo db.WatchRepository, tagRelease []mygithub.TagRelease, jst *time.Location) {
	fmt.Printf("リポジトリ %s/%s の新しいタグリリース:\n", repo.Owner, repo.Name)
	for _, tag := range tagRelease {
		publishedAtJST := tag.PublishedAt.In(jst)
		fmt.Printf("-----------------------------\n")
		fmt.Printf("リリース名: %s\nタグ: %s\n公開日: %s\n", tag.Name, tag.Name, publishedAtJST.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("-----------------------------\n")
}

func updateLastNotificationDate(dbClient *gorm.DB, repo db.WatchRepository) {
	currentTime := time.Now()
	if err := dbClient.Model(&repo).Update("LastNotificationDate", currentTime).Error; err != nil {
		log.Printf("データベースの更新に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
	}
}

func updateLastPublishedAt(dbClient *gorm.DB, newRelease *github.RepositoryRelease, repo db.WatchRepository, jst *time.Location) {
	publishedAtJST := newRelease.PublishedAt.Time.In(jst)
	if err := dbClient.Model(&repo).Update("LastPublishedAt", publishedAtJST).Error; err != nil {
		log.Printf("データベースの更新に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
	}
}

func updateLastPublishedAtForTag(dbClient *gorm.DB, tagRelease mygithub.TagRelease, repo db.WatchRepository, jst *time.Location) {
	publishedAtJST := tagRelease.PublishedAt.In(jst)
	if err := dbClient.Model(&repo).Update("LastPublishedAt", publishedAtJST).Error; err != nil {
		log.Printf("データベースの更新に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
	}
}

func sendNotification(dbClient *gorm.DB, repo db.WatchRepository) {
	notification := db.Notification{}
	if err := dbClient.Where("project_id = ?", repo.ProjectID).First(&notification).Error; err != nil {
		log.Printf("通知設定の取得に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
		return
	}
	
	// 一旦プロジェクトIDとリポジトリID、通知設定を出力
	log.Printf("プロジェクトID: %d, リポジトリID: %d, 通知設定: %s", repo.ProjectID, repo.ID, notification.Type)
}
