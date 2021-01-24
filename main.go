package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/zorkian/go-datadog-api"
	"golang.org/x/oauth2"
)

type PR struct {
	Number             *int
	Labels             []*github.Label
	User               *string
	RequestedReviewers []*github.User
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apikey, appkey, err := readDatadogConfig()
	if err != nil {
		return fmt.Errorf("failed to read Datadog Config: %w", err)
	}

	githubToken, err := readGithubConfig()
	if err != nil {
		return fmt.Errorf("failed to read Datadog Config: %w", err)
	}

	repositories, err := getRepositories()
	if err != nil {
		return fmt.Errorf("failed to get GitHub repository name: %w", err)
	}

	repositoryList := parseRepositories(repositories)

	prs, err := getPullRequests(githubToken, repositoryList)
	if err != nil {
		return fmt.Errorf("failed to get PullRequests: %w", err)
	}

	prInfos := getPRInfos(prs)

	ddClient := datadog.NewClient(apikey, appkey)

	customMetrics, err := generateCustomMetrics(prInfos)
	if err != nil {
		return fmt.Errorf("failed to generate CustomMetrics: %w", err)
	}

	if err := sendCustomMetric(ddClient, customMetrics); err != nil {
		return fmt.Errorf("failed to send custom metrics: %w", err)
	}

	return nil
}

func readDatadogConfig() (string, string, error) {
	apikey := os.Getenv("DATADOG_API_KEY")
	if len(apikey) == 0 {
		return "", "", fmt.Errorf("missing environment variable: DATADOG_API_KEY")
	}

	appkey := os.Getenv("DATADOG_APP_KEY")
	if len(appkey) == 0 {
		return "", "", fmt.Errorf("missing environment variable: DATADOG_APP_KEY")
	}

	return apikey, appkey, nil
}

func sendCustomMetric(ddClient *datadog.Client, customMetrics []datadog.Metric) error {
	if err := ddClient.PostMetrics(customMetrics); err != nil {
		return fmt.Errorf("failed to post metrics(%v): %w", customMetrics, err)
	}
	log.Printf("[Info] sent custom metrics. Count: %v", len(customMetrics))
	return nil
}

func readGithubConfig() (string, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if len(githubToken) == 0 {
		return "", fmt.Errorf("missing environment variable: GITHUB_TOKEN")
	}

	return githubToken, nil
}

func getPullRequests(githubToken string, githubRepositories []string) ([]*github.PullRequest, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	prs := []*github.PullRequest{}

	for _, githubRepository := range githubRepositories {
		repo := strings.Split(githubRepository, "/")
		org := repo[0]
		name := repo[1]
		prsInRepo, _, err := client.PullRequests.List(ctx, org, name, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitHub Pull Requests: %w", err)
		}

		prs = append(prs, prsInRepo...)
	}

	return prs, nil
}

func getRepositories() (string, error) {
	githubRepositories := os.Getenv("GITHUB_REPOSITORIES")
	if len(githubRepositories) == 0 {
		return "", fmt.Errorf("missing environment variable: GITHUB_REPOSITORIES")
	}

	return githubRepositories, nil
}

func parseRepositories(repositories string) []string {
	return strings.Split(repositories, ",")
}

func getPRInfos(prs []*github.PullRequest) []PR {
	prInfos := []PR{}

	for _, pr := range prs {
		prInfos = append(prInfos, PR{
			Number:             pr.Number,
			Labels:             pr.Labels,
			User:               pr.User.Login,
			RequestedReviewers: pr.RequestedReviewers,
		})
	}

	return prInfos
}

func generateCustomMetrics(prInfos []PR) ([]datadog.Metric, error) {
	const timeDifferencesToJapan = +9 * 60 * 60
	tz := time.FixedZone("JST", timeDifferencesToJapan)

	customMetrics := []datadog.Metric{}
	now := time.Now().In(tz)
	nowF := float64(now.Unix())

	countf, err := strconv.ParseFloat("1", 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse. error %w", err)
	}

	var labelsTag []string
	var reviewersTag []string
	for _, prInfo := range prInfos {
		labelsTag = []string{}
		reviewersTag = []string{}

		for _, label := range prInfo.Labels {
			labelsTag = append(labelsTag, "label:"+*label.Name)
		}

		for _, reviewer := range prInfo.RequestedReviewers {
			reviewersTag = append(reviewersTag, "reviewer:"+*reviewer.Login)
		}

		labelAndReviewer := append(labelsTag, reviewersTag...)

		customMetrics = append(customMetrics, datadog.Metric{
			Metric: datadog.String("datadog.custom.github.pr.count"),
			Points: []datadog.DataPoint{
				{datadog.Float64(nowF), datadog.Float64(countf)},
			},
			Type: datadog.String("gauge"),
			Tags: append([]string{"number:" + strconv.Itoa(*prInfo.Number), "author:" + *prInfo.User, "repo:" + "quipper/kubernetes-clusters"}, labelAndReviewer...),
		})
	}

	return customMetrics, nil
}
