package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
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

	return newReleases, nil
}


type TagRelease struct {
	Name string
	PublishedAt time.Time
}

func FetchTagReleaseAfter(ctx context.Context, client *github.Client, owner, repo string, after time.Time) ([]TagRelease, error) {
	/* 与えられた日付以降に作成されたタグを取得します */

	// タグの一覧を取得
	tags, _, err := client.Repositories.ListTags(ctx, owner, repo, nil)
	if err != nil {
		return []TagRelease{}, fmt.Errorf("%s/%s のタグ情報の取得に失敗しました: %v", owner, repo, err)
	}
	//　コミットログから作成日を取得
	var tagReleases []TagRelease
	for _, tag := range tags {
		commit, _, err := client.Git.GetCommit(ctx, owner, repo, *tag.Commit.SHA)
		if err != nil {
			return nil, fmt.Errorf("%s/%s のコミット情報の取得に失敗しました: %v", owner, repo, err)
		}
		if commit.Author.Date.After(after) {
			tagReleases = append(tagReleases, TagRelease{
				Name: *tag.Name,
				PublishedAt: *commit.Author.Date,
			})
		}
	}

	return tagReleases, nil
}
