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

	// after 以降のリリースのみにフィルタリング
	var newReleases []*github.RepositoryRelease
	for _, release := range releases {
		if release.PublishedAt != nil && release.PublishedAt.Time.After(after) {
			newReleases = append(newReleases, release)
		}
	}

	return newReleases, nil
}
