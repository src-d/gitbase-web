[![Build Status](https://travis-ci.org/src-d/gitbase-playground.svg)](https://travis-ci.org/src-d/gitbase-playground)
[![codecov.io](https://codecov.io/github/src-d/gitbase-playground/coverage.svg)](https://codecov.io/github/src-d/gitbase-playground)

# Gitbase Playground

Web application to query Git repositories using SQL. Powered by [gitbase](https://github.com/src-d/gitbase), it allows to perform SQL queries on the Git history and the [Universal AST](https://doc.bblf.sh/) of the code itself.

![Screenshot](.github/screenshot.png?raw=true)

# Usage

## With Docker Compose

The easiest way to run Gitbase Playground and its dependencies is using [Docker Compose](https://docs.docker.com/compose/install/).

The first step is to populate a directory with some Git repositories to be served by gitbase before running it. For example:

```bash
$ mkdir $HOME/repos
$ cd $HOME/repos
$ git clone git@github.com:src-d/gitbase.git
$ git clone git@github.com:bblfsh/bblfshd.git
$ git clone git@github.com:src-d/gitbase-playground.git
```

Next you will need to download the `docker-compose.yml` file from this repository and run `docker-compose`. This tool will run three different containers: the playground frontend itself, gitbase, and bblfshd. To kill the running containers use `Ctrl+C`.

```bash
$ wget https://raw.githubusercontent.com/src-d/gitbase-playground/master/docker-compose.yml
$ docker-compose pull
$ GITBASEPG_REPOS_FOLDER=$HOME/repos docker-compose up --force-recreate
```

The server should be now available at [http://localhost:8080](http://localhost:8080).

In case there are any containers left, you can use
```bash
docker-compose down
```
for cleanup.

## Without Docker Compose

The playground will run the queries against a [gitbase](https://docs.sourced.tech/gitbase) server, and will request UASTs to a [bblfsh](https://doc.bblf.sh/) server. Make sure both are properly configured.

Then you can choose to run the Gitbase Playground either as a docker container, or as a binary.

### As a Docker

```bash
$ docker pull srcd/gitbase-playground:latest
$ docker run -d \
    --publish 8080:8080 \
    --env GITBASEPG_DB_CONNECTION="root@tcp(<gitbase-ip>:3306)/none" \
    --env GITBASEPG_BBLFSH_SERVER_URL="<bblfshd-ip>:9432" \
    srcd/gitbase-playground:latest
```

### As a Binary

Download the binary from our [releases section](https://github.com/src-d/gitbase-playground/releases), and run it:

```bash
$ export GITBASEPG_DB_CONNECTION="root@tcp(<gitbase-ip>:3306)/none"
$ export GITBASEPG_BBLFSH_SERVER_URL="<bblfshd-ip>:9432"
$ ./gitbase-playground
```

# Configuration

Any of the previous execution methods accept configuration through the following environment variables.

| Variable | Default value | Meaning |
| -- | -- | -- |
| `GITBASEPG_HOST` | `0.0.0.0` | IP address to bind the HTTP server |
| `GITBASEPG_PORT` | `8080` | Port to bind the HTTP server |
| `GITBASEPG_SERVER_URL` | | URL used to access the application in the form `HOSTNAME[:PORT]`. Leave it unset to allow connections from any proxy or public address |
| `GITBASEPG_DB_CONNECTION` | `root@tcp(localhost:3306)/none?maxAllowedPacket=4194304` | gitbase connection string |
| `GITBASEPG_BBLFSH_SERVER_URL` | `127.0.0.1:9432` | Address where bblfsh server is listening |
| `GITBASEPG_ENV` | `production` | Sets the log level. Use `dev` to enable debug log messages |
| `GITBASEPG_SELECT_LIMIT` | `100` | Default `LIMIT` forced on all the SQL queries done from the UI. Set it to 0 to remove any limit |
| `GITBASEPG_FOOTER_HTML` | | Allows to add any custom html to the page footer. It must be a string encoded in base64. Use it, for example, to add your analytics tracking code snippet  |

# Contribute

[Contributions](https://github.com/src-d/gitbase-playground/issues) are more than welcome, if you are interested please take a look to our [Contributing Guidelines](docs/CONTRIBUTING.md). There you will also find information on how to build and run the project from sources.

# Code of Conduct

All activities under source{d} projects are governed by the [source{d} code of conduct](https://github.com/src-d/guide/blob/master/.github/CODE_OF_CONDUCT.md).

# License

GPL v3.0, see [LICENSE](LICENSE)
