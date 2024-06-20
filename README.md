# Legitrx

This is [Runxi Yu](https://runxiyu.org)'s fork of [icyphox](https://icyphox.sh/)'s [legit](https://git.icyphox.sh/legit/) project. It aims to replace the use of [cgit](https://git.zx2c4.com/cgit/about/) Runxi's servers.

At its current state, you should not use this fork. Use upstream instead.

## Features

- Customizable templates and stylesheets.
- Cloning over HTTPS.
- Not CGI.

## Building

```sh
git clone https://git.sr.ht/~runxiyu/legitrx
go build
```

## Configuration

If no filename is specified via the `--config` option, it looks for `legitrx.yaml` in the current working directory.

```yaml
repo:
  scanPath: /var/www/git
  readme:
    - README
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
  title: Legitrx
  description: Testing
server:
  name: git.runxiyu.org
  host: 127.0.0.1
  port: 5555
  ```

`server.name` is used for `go import` meta tags and clone URLs.

## Notes

- Run legitrx behind a TLS terminating proxy like [relayd(8)](https://man.openbsd.org/relayd.8) or nginx.
- Cloning only works in bare repos.
- Pushing over https, while supported, is disabled because auth is a pain. Use ssh or [Gitolite](https://gitolite.com/gitolite/).
- Paths are [unveil(2)](https://man.openbsd.org/unveil.2)'d on OpenBSD.

## Ideas

- Code highlighting support.

## License

Legitrx is licensed under the MIT license.
