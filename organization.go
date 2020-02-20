package main

import (
	"context"
	"fmt"
	"os"

	graphql "github.com/shurcooL/githubv4"
)

type organization struct {
	Login graphql.String
}

var (
	organizationsQuery struct {
		Organizations struct {
			Nodes    []organization
			PageInfo struct {
				EndCursor   graphql.String
				HasNextPage bool
			}
		} `graphql:"organizations(first: 100, after: $organizationsPage)"`
	}

	organizations []organization
)

func getOrganizations() ([]organization, int) {
	variables := map[string]interface{}{
		"organizationsPage": (*graphql.String)(nil),
	}

	for {
		err := graphqlClient.Query(context.Background(), &organizationsQuery, variables)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", red(err))
			os.Exit(1)
		}

		organizations = append(organizations, organizationsQuery.Organizations.Nodes...)

		// break on last page
		if !organizationsQuery.Organizations.PageInfo.HasNextPage {
			break
		}

		variables["organizationsPage"] = graphql.NewString(organizationsQuery.Organizations.PageInfo.EndCursor)
	}

	return organizations, len(organizations)
}
