# ghe-get-all-owners

[![build](https://github.com/stoe/ghe-get-all-owners/workflows/build/badge.svg)](https://github.com/stoe/ghe-get-all-owners/actions?query=workflow%3Abuild) [![release](https://github.com/stoe/ghe-get-all-owners/workflows/release/badge.svg)](https://github.com/stoe/ghe-get-all-owners/actions?query=workflow%3Arelease)

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
Create a Personal Access Token (PAT) for GitHub Enterprise Server, `https://HOSTNAME/settings/tokens/new?description=ghe-get-all-owners&scopes=read:org,user:email,read:user,site_admin`

### Usage Example

```sh
$ ./ghe-get-all-owners -h github.example.com -t AA123...
```

## License

MIT © [Stefan Stölzle](https://github.com/stoe)
