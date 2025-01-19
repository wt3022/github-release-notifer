package github

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/go-github/github"
	"github.com/wt3022/github-release-notifier/internal/db"
	"golang.org/x/oauth2"
)

func OpenGitHubClient(ctx context.Context, token string) (*github.Client, error) {
	/* GitHubとの通信Clientの作成 */
	if token == "" {
		return nil, fmt.Errorf("GitHubトークンが設定されていません")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return github.NewClient(oauth2.NewClient(ctx, ts)), nil
}

func FetchReleasesAfter(ctx context.Context, client *github.Client, owner, repo string, after time.Time) ([]*github.RepositoryRelease, error) {
	/* DBに登録されている最終取得日以降のリリースを全て取得します */
	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("%s/%s のリリース情報の取得に失敗しました: %v", owner, repo, err)
	}

	// 与えられた日付以降のリリースのみを取得
	var newReleases []*github.RepositoryRelease
	for _, release := range releases {
		if release.PublishedAt != nil && release.PublishedAt.Time.After(after) {
			newReleases = append(newReleases, release)
		}
	}

	// PublishedAtを降順にソート
	sort.Slice(newReleases, func(i, j int) bool {
		return newReleases[i].PublishedAt.Time.After(newReleases[j].PublishedAt.Time)
	})

	return newReleases, nil
}

type TagRelease struct {
	Name        string
	PublishedAt time.Time
}

func saveTagCommit(commit *github.Commit) error {
	dbClient := db.OpenDB()
	dbClient.Create(&db.Commit{
		SHA:        *commit.SHA,
		AuthorDate: *commit.Author.Date,
	})
	return nil
}

func getDBTagCommit(sha string) *github.Commit {
	dbClient := db.OpenDB()
	var commit db.Commit

	// レコードがない場合は nil を返す
	result := dbClient.Where("sha = ?", sha).First(&commit)
	if result.Error != nil {
		return nil
	}
	return &github.Commit{
		SHA: &commit.SHA,
		Author: &github.CommitAuthor{
			Date: &commit.AuthorDate,
		},
	}
}

func FetchTagReleaseAfter(ctx context.Context, client *github.Client, owner, repo string, after time.Time) ([]TagRelease, error) {
	/* 与えられた日付以降に作成されたタグを取得します */

	// タグの一覧を取得
	tags, _, err := client.Repositories.ListTags(ctx, owner, repo, &github.ListOptions{PerPage: 10})
	if err != nil {
		return []TagRelease{}, fmt.Errorf("%s/%s のタグ情報の取得に失敗しました: %v", owner, repo, err)
	}

	//　コミットログから作成日を取得
	var tagReleases []TagRelease

	for _, tag := range tags {
		commit := getDBTagCommit(*tag.Commit.SHA)
		if commit == nil {
			fmt.Println("DBに保存されていないコミットを取得します")
			newCommit, _, err := client.Git.GetCommit(ctx, owner, repo, *tag.Commit.SHA)
			if err != nil {
				return nil, fmt.Errorf("%s/%s のコミット情報の取得に失敗しました: %v", owner, repo, err)
			}
			if err := saveTagCommit(newCommit); err != nil {
				return nil, fmt.Errorf("コミット情報の保存に失敗しました: %v", err)
			}
			commit = newCommit
		}

		if commit != nil && commit.Author != nil && commit.Author.Date != nil && commit.Author.Date.After(after) {
			fmt.Printf("commit: %+v\n", commit)
			tagReleases = append(tagReleases, TagRelease{
				Name:        *tag.Name,
				PublishedAt: *commit.Author.Date,
			})
		}

		// PublishedAtを降順にソート
		sort.Slice(tagReleases, func(i, j int) bool {
			return tagReleases[i].PublishedAt.After(tagReleases[j].PublishedAt)
		})
	}

	// PublishedAtを降順にソート
	sort.Slice(tagReleases, func(i, j int) bool {
		return tagReleases[i].PublishedAt.After(tagReleases[j].PublishedAt)
	})

	return tagReleases, nil
}
