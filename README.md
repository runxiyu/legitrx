# legitrx

This is [Runxi Yu](https://runxiyu.org)'s fork of
[icyphox](https://icyphox.sh/)'s
[legitrx](https://git.icyphox.sh/legitrx/) project

The README hasn't really been updated yet and not much has happened in
the fork yet. Just use upstream for now.

A git web frontend written in Go.

## FEATURES

- Fully customizable templates and stylesheets.
- Cloning over http(s).
- Less archaic HTML.
- Not CGI.

## INSTALLING

Clone it, `go build` it.

## CONFIG

Looks for a `legitrx.yaml` in the current directory by default; pass the
`--config` flag to point it elsewhere.

```yaml
repo:
  scanPath: /var/www/git
  readme:
    - readme
    - README
    - readme.md
    - README.md
  mainBranch:
    - master
    - main
  ignore:
    - foo
    - bar
dirs:
  templates: ./templates
  static: ./static
meta:
  title: git good
  description: i think it's a skill issue
server:
  name: git.icyphox.sh
  host: 127.0.0.1
  port: 5555
  ```

These options are fairly self-explanatory, but of note are:

- repo.scanPath: where all your git repos live. legitrx doesn't
  traverse subdirs yet.
- dirs: use this to override the default templates and static assets.
- repo.readme: readme files to look for.
- repo.mainBranch: main branch names to look for.
- repo.ignore: repos to ignore, relative to scanPath.
- server.name: used for go-import meta tags and clone URLs.

## NOTES

- Run legitrx behind a TLS terminating proxy like relayd(8) or nginx.
- Cloning only works in bare repos -- this is a limitation inherent to git. You
  can still view bare repos just fine in legitrx.
- The default head.html template uses my CDN to fetch fonts -- you may
  or may not want this.
- Pushing over https, while supported, is disabled because auth is a
  pain. Use ssh.
- Paths are unveil(2)'d on OpenBSD.

## IDEAS

- Support for filters for markdown rendering and code highlighting. 

## LICENSE

legitrx is licensed under MIT.
