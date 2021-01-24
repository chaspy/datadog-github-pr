package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/zorkian/go-datadog-api"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	type PR struct {
		Number              *int
		Labels              []*github.Label
		User                *string
		Assignee            *string
		Assignees           []*string
		RequestedReviewers  []*github.User
	}

	const timeoutsecond = 10

	apikey, appkey, err := readDatadogConfig()
	if err != nil {
		return fmt.Errorf("failed to read Datadog Config: %w", err)
	}

	githubToken, err := readGithubConfig()
	if err != nil {
		return fmt.Errorf("failed to read Datadog Config: %w", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	ctx := context.Background()

	client := github.NewClient(tc)

	prs, _, err := client.PullRequests.List(ctx, "quipper", "kubernetes-clusters", nil)
	fmt.Printf("%v\n", *prs[0].Labels[0].Name)
	fmt.Printf("%v\n", *prs[0].User.Login)
	fmt.Printf("%v\n", *prs[0].RequestedReviewers[0].Login)

	var prinfos []PR

	for _, pr := range prs {
		prinfos = append(prinfos, PR{
			Number: pr.Number,
			Labels: pr.Labels,
			User:  pr.User.Login,
			RequestedReviewers: pr.RequestedReviewers,
		})
	}

	fmt.Println("%v",prinfos)

	ddClient := datadog.NewClient(apikey, appkey)

	const timeDifferencesToJapan = +9 * 60 * 60
	tz := time.FixedZone("JST", timeDifferencesToJapan)

	var customMetrics []datadog.Metric
	now := time.Now().In(tz)
	nowF := float64(now.Unix())

	count := "1"
	countf, err := strconv.ParseFloat(count, 64)
	if err != nil {
		return fmt.Errorf("failed to parse. count %v, error %w", count, err)
	}

	for _, prinfo := range prinfos {
		customMetrics = append(customMetrics, datadog.Metric{
			Metric: datadog.String("datadog.custom.github.pr.count"),
			Points: []datadog.DataPoint{
				{datadog.Float64(nowF), datadog.Float64(countf)},
			},
			Type: datadog.String("gauge"),
			Tags: []string{"number:" + strconv.Itoa(*prinfo.Number),"author:" + *prinfo.User, "repo:" + "quipper/kubernetes-clusters"},
		})
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
		return  "", fmt.Errorf("missing environment variable: GITHUB_TOKEN")
	}

	return githubToken, nil
}
