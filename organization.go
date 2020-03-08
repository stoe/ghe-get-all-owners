package main

import (
	"context"

	graphql "github.com/shurcooL/githubv4"
)

type organization struct {
	Login graphql.String
}

var (
	oQuery struct {
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

func getOrganizations() (o []organization) {
	variables := map[string]interface{}{
		"organizationsPage": (*graphql.String)(nil),
	}

	for {
		if err := graphqlClient.Query(context.Background(), &oQuery, variables); err != nil {
			errorAndExit(err)
		}

		o = append(o, oQuery.Organizations.Nodes...)

		// break on last page
		if !oQuery.Organizations.PageInfo.HasNextPage {
			break
		}

		variables["organizationsPage"] = graphql.NewString(oQuery.Organizations.PageInfo.EndCursor)
	}

	return
}
