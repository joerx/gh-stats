package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

// Formatter defines an interface to output data in various formats
type Formatter interface {
	Write(data interface{}) error
	Flush()
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "json",
			Usage: "Output format",
		},
		cli.StringFlag{
			Name:   "org",
			Usage:  "Github org to list repos for",
			EnvVar: "GITHUB_ORG_NAME",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Github access token",
			EnvVar: "GITHUB_TOKEN",
		},
	}

	app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func listOrgRepos(ctx context.Context, client *github.Client, orgName string) ([]*github.Repository, error) {
	allRepos := []*github.Repository{}

	opts := &github.RepositoryListByOrgOptions{Type: "private"}
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opts)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func action(c *cli.Context) {

	format := c.String("output")
	token := c.String("token")
	orgName := c.String("org")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repos, err := listOrgRepos(ctx, client, orgName)

	if err != nil {
		log.Fatal(err)
	}

	var f Formatter

	switch format {
	case "json":
		f = NewJSONFormatter(os.Stdout)
	case "csv":
		f = NewCSVFormatter(os.Stdout, []string{"ID", "Name", "CreatedAt", "UpdatedAt"}, nil)
	}

	for _, r := range repos {
		f.Write(r)
	}
	f.Flush()
}
