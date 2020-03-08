package main

import (
	"context"

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
	eQuery struct {
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
		if err := graphqlClient.Query(context.Background(), &eQuery, variables); err != nil {
			errorAndExit(err)
		}
		members = eQuery.Enterprise.Members.Nodes

		// break on last page
		if !eQuery.Enterprise.Members.PageInfo.HasNextPage {
			break
		}

		variables["memberPage"] = graphql.NewString(eQuery.Enterprise.Members.PageInfo.EndCursor)
	}

	c <- members
}
