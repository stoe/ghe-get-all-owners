package main

import (
	"context"
	"fmt"
	"os"

	graphql "github.com/shurcooL/githubv4"
)

type member struct {
	User struct {
		Login graphql.String
		Email graphql.String
		Name  graphql.String
	} `graphql:"... on User"`
}

var (
	enterpriseQuery struct {
		Enterprise struct {
			Members struct {
				PageInfo struct {
					EndCursor   graphql.String
					HasNextPage bool
				}
				Nodes []member
			} `graphql:"members(organizationLogins: [$login], role: OWNER, first: 100, after: $memberPage)"`
		} `graphql:"enterprise(slug: \"github\")"`
	}

	members []member
)

func getOwners(login graphql.String, c chan []member) {
	variables := map[string]interface{}{
		"login":      login,
		"memberPage": (*graphql.String)(nil),
	}

	for {
		err := graphqlClient.Query(context.Background(), &enterpriseQuery, variables)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", red(err))
			os.Exit(1)
		}
		members = enterpriseQuery.Enterprise.Members.Nodes

		// break on last page
		if !enterpriseQuery.Enterprise.Members.PageInfo.HasNextPage {
			break
		}

		variables["memberPage"] = graphql.NewString(enterpriseQuery.Enterprise.Members.PageInfo.EndCursor)
	}

	c <- members
}
