# ghe-get-all-owners

> Get all organization owners of a GitHub Enterprise Server

[![build](https://github.com/stoe/ghe-get-all-owners/workflows/build/badge.svg)](https://github.com/stoe/ghe-get-all-owners/actions?query=workflow%3Abuild) [![release](https://github.com/stoe/ghe-get-all-owners/workflows/release/badge.svg)](https://github.com/stoe/ghe-get-all-owners/actions?query=workflow%3Arelease)

:information_source: Required [GitHub Enterprise Server](https://github.com/enterprise) version: 2.19 or newer

## Install

```sh
$ go get github.com/stoe/ghe-get-all-owners
```

Or download the the latest release binary for your platform: [github.com/stoe/ghe-get-all-owners](https://github.com/stoe/ghe-get-all-owners/releases)

## Usage

```sh
USAGE:
  ghe-get-all-owners [OPTIONS]

OPTIONS:
  -h, --hostname string       hostname
  -t, --token string          personal access token
      --help                  print this help

EXAMPLES:
  $ ghe-get-all-owners -h github.example.com -t AA123...
```

The scripts requires a personal access token with at least the `read:org,user:email,read:user` scope, better yet to use one with additional `site_admin` scope.

Create a Personal Access Token (PAT) for GitHub Enterprise Server, `https://HOSTNAME/settings/tokens/new?description=ghe-get-all-owners&scopes=read:org,user:email,read:user,site_admin`.

### Usage Example

```sh
$ ghe-get-all-owners -h github.example.com -t AA123...
```

### Usage Result

`ghe-get-all-owners` creates a `ghes-all-owners.csv` file in the current directory. Example:

```csv
organization,login,name,email
foo,hubot,Hubot Robot,hubot@example.com
foo,mona,Mona Octocat,mona@example.com
bar,hubot,Hubot Robot,hubot@example.com
```

## License

MIT © [Stefan Stölzle](https://github.com/stoe)
