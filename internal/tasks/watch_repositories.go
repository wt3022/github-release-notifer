package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/github"
	"github.com/wt3022/github-release-notifier/internal/db"
	"github.com/wt3022/github-release-notifier/internal/env"
	mygithub "github.com/wt3022/github-release-notifier/internal/github"
	"github.com/wt3022/github-release-notifier/internal/utils"
	"gorm.io/gorm"
)

type Release struct {
	ReleaseName *string
	TagName     *string
	PublishedAt time.Time
}

func WatchRepositoryRelease(dbClient *gorm.DB, githubClient *github.Client) {
	log.Printf("リポジトリの処理を開始します\n")

	env := env.LoadConfig()

	var watchRepos []db.WatchRepository
	var allNotifications []struct {
		Repository db.WatchRepository
		Releases   []Release
	}

	if err := dbClient.Find(&watchRepos).Error; err != nil {
		log.Fatalf("データベースからのリポジトリの取得に失敗しました: %v", err)
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("タイムゾーンのロードに失敗しました: %v", err)
	}

	for _, repo := range watchRepos {
		releases := processRepository(dbClient, githubClient, repo, jst, env)
		if len(releases) > 0 {
			allNotifications = append(allNotifications, struct {
				Repository db.WatchRepository
				Releases   []Release
			}{
				Repository: repo,
				Releases:   releases,
			})
		}
	}

	if len(allNotifications) > 0 {
		sendBatchNotification(dbClient, allNotifications, env, jst)
	}

	fmt.Printf("リポジトリの処理が完了しました\n")
}

func processRepository(dbClient *gorm.DB, githubClient *github.Client, repo db.WatchRepository, jst *time.Location, env env.Config) []Release {
	var releases []Release

	if repo.WatchType == db.WatchTypeRelease {
		newReleases, err := mygithub.FetchReleasesAfter(context.Background(), githubClient, repo.Owner, repo.Name, repo.LastPublishedAt)
		if err != nil {
			log.Println(err)
			return nil
		}

		for _, release := range newReleases {
			releases = append(releases, Release{
				ReleaseName: release.Name,
				TagName:     nil,
				PublishedAt: release.PublishedAt.Time,
			})
		}

		if len(releases) > 0 {
			updateLastPublishedAt(dbClient, releases[0], repo, jst, "release")
		} else {
			fmt.Printf("リポジトリ %s/%s に新しいリリースはありません\n", repo.Owner, repo.Name)
		}
	} else if repo.WatchType == db.WatchTypeTag {
		tagRelease, err := mygithub.FetchTagReleaseAfter(context.Background(), githubClient, repo.Owner, repo.Name, repo.LastPublishedAt)
		if err != nil {
			log.Println(err)
			return nil
		}

		for _, tag := range tagRelease {
			releases = append(releases, Release{
				ReleaseName: nil,
				TagName:     &tag.Name,
				PublishedAt: tag.PublishedAt,
			})
		}

		if len(releases) > 0 {
			updateLastPublishedAt(dbClient, releases[0], repo, jst, "tag")
		} else {
			fmt.Printf("リポジトリ %s/%s に新しいタグリリースはありません\n", repo.Owner, repo.Name)
		}
	}

	updateLastNotificationDate(dbClient, repo)
	return releases
}

func updateLastNotificationDate(dbClient *gorm.DB, repo db.WatchRepository) {
	currentTime := time.Now()
	if err := dbClient.Model(&repo).Update("LastNotificationDate", currentTime).Error; err != nil {
		log.Printf("データベースの更新に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
	}
}

func updateLastPublishedAt(dbClient *gorm.DB, newRelease Release, repo db.WatchRepository, jst *time.Location, updateType string) {
	if updateType == "release" {
		publishedAtJST := newRelease.PublishedAt.In(jst)
		if err := dbClient.Model(&repo).Update("LastPublishedAt", publishedAtJST).Error; err != nil {
			log.Printf("データベースの更新に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
		}
	} else if updateType == "tag" {
		publishedAtJST := newRelease.PublishedAt.In(jst)
		if err := dbClient.Model(&repo).Update("LastPublishedAt", publishedAtJST).Error; err != nil {
			log.Printf("データベースの更新に失敗しました (%s/%s): %v", repo.Owner, repo.Name, err)
		}
	}
}

func sendBatchNotification(dbClient *gorm.DB, notifications []struct {
	Repository db.WatchRepository
	Releases   []Release
}, env env.Config, jst *time.Location) {
	fmt.Println("run sendBatchNotification")

	var emailBody string
	emailBody = "監視中の以下のリポジトリに更新がありました：\n\n"

	for _, notification := range notifications {
		repo := notification.Repository
		releases := notification.Releases

		emailBody += fmt.Sprintf("リポジトリ: %s/%s\n", repo.Owner, repo.Name)
		for _, release := range releases {
			publishedAtJST := release.PublishedAt.In(jst)
			if release.ReleaseName != nil {
				emailBody += fmt.Sprintf("リリース名: %s\n公開日: %s\n\n",
					*release.ReleaseName,
					publishedAtJST.Format("2006-01-02 15:04:05"))
			} else if release.TagName != nil {
				emailBody += fmt.Sprintf("タグ: %s\n公開日: %s\n\n",
					*release.TagName,
					publishedAtJST.Format("2006-01-02 15:04:05"))
			}
		}
		emailBody += "-------------------\n\n"
	}

	notificationSettings := db.Notification{}
	if err := dbClient.First(&notificationSettings).Error; err != nil {
		log.Printf("通知設定の取得に失敗しました: %v", err)
		return
	}

	if notificationSettings.Type == db.Email {
		emailRequest := utils.EmailRequest{
			To:      "test@example.com",
			Subject: "監視中のリポジトリに更新がありました",
			Body:    emailBody,
		}
		utils.SendEmail(emailRequest, env)
	}
}
