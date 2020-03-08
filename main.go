package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"

	rest "github.com/google/go-github/v29/github"
	graphql "github.com/shurcooL/githubv4"
)

const minVersion = "2.21.0"

var (
	// options
	help     bool
	hostname string
	token    string

	file     *os.File
	filepath = "./ghes-all-owners.csv"
	header   = []string{"organization", "login", "name", "email"}

	httpClient    *http.Client
	restClient    *rest.Client
	graphqlClient *graphql.Client

	ctx = context.Background()
)

func init() {
	// flags
	pflag.StringVarP(&hostname, "hostname", "h", "", "hostname")
	pflag.StringVarP(&token, "token", "t", "", "personal access token")
	pflag.BoolVar(&help, "help", false, "print this help")
	pflag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient = oauth2.NewClient(ctx, src)

	graphqlURL := fmt.Sprintf("https://%s/api/graphql", hostname)
	graphqlClient = graphql.NewEnterpriseClient(graphqlURL, httpClient)

	restURL := fmt.Sprintf("https://%s/api/v3", hostname)
	restClient, _ = rest.NewEnterpriseClient(restURL, restURL, httpClient)

	validateFlags()
	checkVersion()

	// delete previousely generated file
	os.Remove(filepath)

	var err error
	if file, err = os.Create(filepath); err != nil {
		errorAndExit(err)
	}
}

func main() {
	start := time.Now()

	writer := csv.NewWriter(file)
	writer.Write(header)

	orgs := getOrganizations()
	c := make(chan []member)

	for _, org := range orgs {
		login := org.Login

		// skip the default organization
		if login == "github-enterprise" {
			continue
		}

		fmt.Printf(
			"Looking up owners for https://%s/%s ...",
			hostname,
			login,
		)

		go getOwners(login, c)

		m := <-c
		n := len(m)

		fmt.Printf("found %d\n", n)

		if n > 0 {
			for _, u := range m {
				writer.Write([]string{
					fmt.Sprintf("%s", login),
					fmt.Sprintf("%s", u.User.Login),
					fmt.Sprintf("%s", u.User.Name),
					fmt.Sprintf("%s", u.User.Email),
				})
			}
		} else {
			writer.Write([]string{fmt.Sprintf("%s", login), "", "", ""})
		}

		writer.Flush()
	}

	if err := writer.Error(); err != nil {
		errorAndExit(err)
	}

	if err := file.Close(); err != nil {
		errorAndExit(err)
	}

	fmt.Printf("\nFile saved to %s\n", filepath)

	fmt.Printf(
		"\nDone after %s\n",
		time.Now().Sub(start).Round(time.Millisecond),
	)
}

// helpers -------------------------------------------------------------------------------------------------------------

func validateFlags() {
	if help {
		printHelp()
		os.Exit(0)
	}

	if hostname == "" {
		printHelpOnError("hostname missing")
	}

	if hostname == "github.com" {
		printHelpOnError("github.com is not supported")
	}

	if token == "" {
		printHelpOnError("token missing")
	}
}

func checkVersion() {
	minV := *semver.New(minVersion)

	res, err := httpClient.Get(
		fmt.Sprintf("https://%s/api/v3/meta", hostname),
	)

	if err != nil {
		errorAndExit(err)
	}

	if res.StatusCode != 200 {
		errorAndExit(errors.New(res.Status))
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		errorAndExit(err)
	}

	// https://developer.github.com/enterprise/v3/meta/
	var meta struct {
		Auth    bool   `json:"verifiable_password_authentication"`
		Version string `json:"installed_version"`
	}
	json.Unmarshal(body, &meta)

	v := semver.New(meta.Version)

	if v.LessThan(minV) {
		err := fmt.Errorf(
			"need GHES version >= %s, but got %s",
			minV.String(),
			v.String(),
		)

		errorAndExit(err)
	}
}

func printHelp() {
	fmt.Println(`USAGE:
  ghe-get-all-owners [OPTIONS]

OPTIONS:`)
	pflag.PrintDefaults()
	fmt.Println(`
EXAMPLE:
  $ ghe-get-all-owners -h github.example.com -t AA123...`)
	fmt.Println()
}

func printHelpOnError(s string) {
	printHelp()
	errorAndExit(errors.New(s))
}

func errorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(2)
}
