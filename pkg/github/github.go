package github

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-github/v87/github"
	"github.com/shurcooL/githubv4"
	"github.com/suzuki-shunsuke/tfaction-go/pkg/util"
	"golang.org/x/oauth2"
)

type (
	ListOptions = github.ListOptions
	Response    = github.Response

	GitHubIssue  = github.Issue //nolint:revive
	IssueRequest = github.IssueRequest
)

type ParamNew struct {
	Token              string
	GHEBaseURL         string
	GHEGraphQLEndpoint string
}

type ClientImpl struct {
	v4Client v4Client
	issue    IssueClient
}

func New(ctx context.Context, param *ParamNew) (*ClientImpl, error) {
	httpClient := getHTTPClientForGitHub(ctx, param.Token)
	client := &ClientImpl{}

	opts := []github.ClientOptionsFunc{github.WithHTTPClient(httpClient)}
	if param.GHEBaseURL != "" {
		opts = append(opts, github.WithEnterpriseURLs(param.GHEBaseURL, param.GHEBaseURL))
	}
	gh, err := github.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("initialize GitHub API Client: %w", err)
	}
	client.issue = gh.Issues

	if param.GHEGraphQLEndpoint == "" {
		client.v4Client = githubv4.NewClient(httpClient)
	} else {
		client.v4Client = githubv4.NewEnterpriseClient(param.GHEGraphQLEndpoint, httpClient)
	}

	return client, nil
}

const (
	keySearchQuery  = "searchQuery"
	keySearchType   = "searchType"
	keyIssuesCursor = "issuesCursor"
)

type v4Client interface {
	Query(ctx context.Context, q any, variables map[string]any) error
}

func getHTTPClientForGitHub(ctx context.Context, token string) *http.Client {
	if token == "" {
		return http.DefaultClient
	}
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
}

//go:generate mockery --name Client --testonly=false
type Client interface {
	ListIssues(ctx context.Context, repoOwner, repoName string) ([]*Issue, error)
	ListLeastRecentlyUpdatedIssues(ctx context.Context, repoOwner, repoName string, numOfIssues int, deadline string) ([]*Issue, error)
	CreateIssue(ctx context.Context, repoOwner, repoName string, param *IssueRequest) (*GitHubIssue, error)
	CloseIssue(ctx context.Context, repoOwner, repoName string, issueNumber int) (*GitHubIssue, error)
	GetIssue(ctx context.Context, repoOwner, repoName, title string) (*Issue, error)
	ArchiveIssue(ctx context.Context, repoOwner, repoName string, issueNumber int, title string) (*GitHubIssue, error)
}

type IssueClient interface {
	Create(ctx context.Context, owner string, repo string, issue *IssueRequest) (*GitHubIssue, *Response, error)
	Edit(ctx context.Context, owner string, repo string, issueNumber int, issue *IssueRequest) (*GitHubIssue, *Response, error)
}

type Issue struct {
	Number int    `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
	Target string `json:"target,omitempty"`
	State  string `json:"state,omitempty"`
	RunsOn string `json:"runs_on,omitempty"`
}

var titlePattern = regexp.MustCompile(`^Terraform Drift \((\S+)\)$`)

func (cl *ClientImpl) ListIssues(ctx context.Context, repoOwner, repoName string) ([]*Issue, error) {
	var q struct {
		Search struct {
			Nodes []struct {
				Issue struct {
					Number githubv4.Int
					Title  githubv4.String
				} `graphql:"... on Issue"`
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"search(first: 100, after: $issuesCursor, query: $searchQuery, type: $searchType)"`
	}
	variables := map[string]any{
		keySearchQuery:  githubv4.String(fmt.Sprintf(`repo:%s/%s "Terraform Drift" in:title`, repoOwner, repoName)),
		keySearchType:   githubv4.SearchTypeIssue,
		keyIssuesCursor: (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var allIssues []*Issue
	for {
		if err := cl.v4Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("list issue comments by GitHub API: %w", err)
		}
		for _, issue := range q.Search.Nodes {
			title := string(issue.Issue.Title)
			a := titlePattern.FindStringSubmatch(title)
			if a == nil {
				continue
			}
			allIssues = append(allIssues, &Issue{
				Number: int(issue.Issue.Number),
				Title:  title,
				Target: a[1],
			})
		}
		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables[keyIssuesCursor] = githubv4.NewString(q.Search.PageInfo.EndCursor)
	}
	return allIssues, nil
}

func (cl *ClientImpl) CreateIssue(ctx context.Context, repoOwner, repoName string, issue *IssueRequest) (*GitHubIssue, error) {
	ret, _, err := cl.issue.Create(ctx, repoOwner, repoName, issue)
	if err != nil {
		return nil, fmt.Errorf("create an issue by GitHub API v3: %w", err)
	}
	return ret, nil
}

func (cl *ClientImpl) CloseIssue(ctx context.Context, repoOwner, repoName string, issueNumber int) (*GitHubIssue, error) {
	ret, _, err := cl.issue.Edit(ctx, repoOwner, repoName, issueNumber, &IssueRequest{
		State: util.StrP("closed"),
	})
	if err != nil {
		return nil, fmt.Errorf("close an issue by GitHub API v3: %w", err)
	}
	return ret, nil
}

func (cl *ClientImpl) ArchiveIssue(ctx context.Context, repoOwner, repoName string, issueNumber int, title string) (*GitHubIssue, error) {
	ret, _, err := cl.issue.Edit(ctx, repoOwner, repoName, issueNumber, &IssueRequest{
		State: util.StrP("closed"),
		Title: util.StrP(title),
	})
	if err != nil {
		return nil, fmt.Errorf("edit an issue by GitHub API v3: %w", err)
	}
	return ret, nil
}

func (cl *ClientImpl) ListLeastRecentlyUpdatedIssues(ctx context.Context, repoOwner, repoName string, numOfIssues int, deadline string) ([]*Issue, error) {
	var q struct {
		Search struct {
			Nodes []struct {
				Issue struct {
					Number githubv4.Int
					Title  githubv4.String
					State  githubv4.String
				} `graphql:"... on Issue"`
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"search(first: 100, after: $issuesCursor, query: $searchQuery, type: $searchType)"`
	}
	variables := map[string]any{
		keySearchQuery:  githubv4.String(fmt.Sprintf(`repo:%s/%s "Terraform Drift" in:title sort:updated-asc updated:<%s`, repoOwner, repoName, deadline)),
		keySearchType:   githubv4.SearchTypeIssue,
		keyIssuesCursor: (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var allIssues []*Issue
	for {
		if err := cl.v4Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("list issue comments by GitHub API: %w", err)
		}
		for _, issue := range q.Search.Nodes {
			title := string(issue.Issue.Title)
			a := titlePattern.FindStringSubmatch(title)
			if a == nil {
				continue
			}
			allIssues = append(allIssues, &Issue{
				Number: int(issue.Issue.Number),
				Title:  title,
				Target: a[1],
				State:  strings.ToLower(string(issue.Issue.State)),
			})
			if len(allIssues) == numOfIssues {
				return allIssues, nil
			}
		}
		if !q.Search.PageInfo.HasNextPage {
			return allIssues, nil
		}
		variables[keyIssuesCursor] = githubv4.NewString(q.Search.PageInfo.EndCursor)
	}
}

func (cl *ClientImpl) GetIssue(ctx context.Context, repoOwner, repoName, title string) (*Issue, error) {
	var q struct {
		Search struct {
			Nodes []struct {
				Issue struct {
					Number githubv4.Int
					Title  githubv4.String
					State  githubv4.String
				} `graphql:"... on Issue"`
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"search(first: 100, after: $issuesCursor, query: $searchQuery, type: $searchType)"`
	}
	variables := map[string]any{
		keySearchQuery:  githubv4.String(fmt.Sprintf(`repo:%s/%s "%s" in:title`, repoOwner, repoName, title)),
		keySearchType:   githubv4.SearchTypeIssue,
		keyIssuesCursor: (*githubv4.String)(nil), // Null after argument to get first page.
	}

	for {
		if err := cl.v4Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("list issue comments by GitHub API: %w", err)
		}
		for _, issue := range q.Search.Nodes {
			if string(issue.Issue.Title) != title {
				continue
			}
			return &Issue{
				Number: int(issue.Issue.Number),
				State:  strings.ToLower(string(issue.Issue.State)),
			}, nil
		}
		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables[keyIssuesCursor] = githubv4.NewString(q.Search.PageInfo.EndCursor)
	}
	return nil, nil //nolint:nilnil
}
