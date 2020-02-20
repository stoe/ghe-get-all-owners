package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/gookit/color"
	"golang.org/x/oauth2"

	rest "github.com/google/go-github/v29/github"
	graphql "github.com/shurcooL/githubv4"
	flag "github.com/spf13/pflag"
)

type records [][]string

// Meta https://developer.github.com/enterprise/v3/meta/
type Meta struct {
	Auth    bool   `json:"verifiable_password_authentication"`
	Version string `json:"installed_version"`
}

var (
	// options
	help     bool
	hostname string
	token    string

	httpClient    *http.Client
	restClient    *rest.Client
	graphqlClient *graphql.Client

	blue   = color.FgBlue.Render
	dimmed = color.OpFuzzy.Render
	green  = color.FgGreen.Render
	red    = color.FgRed.Render

	ctx = context.Background()
)

const minVersion = "2.19.0"

func main() {
	err := setup()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", red(err))
		os.Exit(1)
	}

	start := time.Now()

	orgs, _ := getOrganizations()

	// csv headers
	r := records{
		[]string{"organization", "login", "name", "email"},
	}

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
				r = append(r, []string{
					fmt.Sprintf("%s", login),
					fmt.Sprintf("%s", u.User.Login),
					fmt.Sprintf("%s", u.User.Name),
					fmt.Sprintf("%s", u.User.Email),
				})
			}
		} else {
			r = append(r, []string{fmt.Sprintf("%s", login), "", "", ""})
		}
	}

	fp := r.saveCSV("ghes-find-owners.csv")

	fmt.Printf("\nFile saved to %s\n", blue(fp))

	fmt.Printf(
		"\nDone after %s\n\n",
		time.Now().Sub(start).Round(time.Millisecond),
	)
}

func (r records) saveCSV(filename string) string {
	path := "./dist/"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	filepath := path + filename

	// delete previousely generated file
	os.Remove(filepath)

	file, err := os.Create(filepath)
	defer file.Close()

	if err != nil {
		log.Fatalln("error creating csv:", err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(r)

	if err := writer.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	return filepath
}

// helpers -------------------------------------------------------------------------------------------------------------

func setup() error {
	flags := flag.NewFlagSet("ghe-get-all-owners", flag.ContinueOnError)

	flags.StringVarP(&hostname, "hostname", "h", "", "hostname")
	flags.StringVarP(&token, "token", "t", "", "personal access token")
	flags.BoolVar(&help, "help", false, "print this help")

	flags.SortFlags = false

	err := flags.Parse(os.Args[1:])

	if err != nil {
		printHelp(flags)

		return errors.New(red(err))
	}

	args := flags.Args()
	if len(args) != 0 {
		printHelp(flags)

		return errors.New(red("excess arguments"))
	}

	if help {
		printHelp(flags)

		os.Exit(0)
	}

	if hostname == "" {
		return printHelpOnError("hostname", flags)
	}

	if hostname == "github.com" {
		return printHelpOnError("github.com is not supported", flags)
	}

	if token == "" {
		return printHelpOnError("token", flags)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient = oauth2.NewClient(ctx, src)

	graphqlURL := fmt.Sprintf("https://%s/api/graphql", hostname)
	graphqlClient = graphql.NewEnterpriseClient(graphqlURL, httpClient)

	restURL := fmt.Sprintf("https://%s/api/v3", hostname)
	restClient, _ = rest.NewEnterpriseClient(restURL, restURL, httpClient)

	checkVersion()

	return nil
}

func checkVersion() {
	minV := *semver.New(minVersion)

	res, err := httpClient.Get(
		fmt.Sprintf("https://%s/api/v3/meta", hostname),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", red(err))
		os.Exit(1)
	}

	defer res.Body.Close()

	meta := &Meta{}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", red(err))
		os.Exit(1)
	}

	json.Unmarshal(body, meta)

	v := semver.New(meta.Version)

	if v.LessThan(minV) {
		msg := fmt.Sprintf(
			"Need GHES version >= %s, but got %s.",
			green(minV.String()),
			red(v.String()),
		)

		fmt.Fprintf(
			os.Stderr,
			"error: %s\n",
			msg,
		)
		os.Exit(1)
	}
}

func printHelp(f *flag.FlagSet) {
	fmt.Println(`USAGE:
  ghe-get-all-owners [OPTIONS]

OPTIONS:`)
	f.PrintDefaults()
	fmt.Println(`
EXAMPLES:
  $ ghe-get-all-owners-darwin-amd64 -h github.example.com -t AA123...
  $ ghe-get-all-owners-windows-amd64.exe -h github.example.com -t AA123...`)
	fmt.Println("")
}

func printHelpOnError(t string, f *flag.FlagSet) error {
	printHelp(f)

	msg := fmt.Sprintf(red("%s missing"), t)

	return errors.New(msg)
}
