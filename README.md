[![Build Status](https://travis-ci.com/src-d/gitbase-web.svg?branch=master)](https://travis-ci.com/src-d/gitbase-web)
[![codecov.io](https://codecov.io/github/src-d/gitbase-web/coverage.svg)](https://codecov.io/github/src-d/gitbase-web)

# Gitbase Web

Application to query Git repositories using SQL. Powered by [gitbase](https://github.com/src-d/gitbase), it allows to perform SQL queries on the Git history and the [Universal AST](https://doc.bblf.sh/) of the code itself.

![Screenshot](.github/screenshot.png?raw=true)

# Usage

## With Docker Compose

The easiest way to run Gitbase Web and its dependencies is using [Docker Compose](https://docs.docker.com/compose/install/).

The first step is to populate a directory with some Git repositories to be served by gitbase before running it. For example:

```bash
$ mkdir $HOME/repos
$ cd $HOME/repos
$ git clone git@github.com:src-d/gitbase.git
$ git clone git@github.com:bblfsh/bblfshd.git
$ git clone git@github.com:src-d/gitbase-web.git
```

Next you will need to download the `docker-compose.yml` file from this repository and run `docker-compose up`. This tool will run three different containers: the gitbase-web frontend itself, gitbase, and bblfshd. To kill the running containers use `Ctrl+C`.

```bash
$ wget https://raw.githubusercontent.com/src-d/gitbase-web/master/docker-compose.yml
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

The application will run the queries against a [gitbase](https://docs.sourced.tech/gitbase) server, and will request UASTs to a [bblfsh](https://doc.bblf.sh/) server. Make sure both are properly configured.

Then you can choose to run the web client either as a docker container, or as a binary.

### As a Docker

```bash
$ docker pull srcd/gitbase-web:latest
$ docker run -d \
    --publish 8080:8080 \
    --env GITBASEPG_DB_CONNECTION="root@tcp(<gitbase-ip>:3306)/none" \
    --env GITBASEPG_BBLFSH_SERVER_URL="<bblfshd-ip>:9432" \
    srcd/gitbase-web:latest
```

### As a Binary

Download the binary from our [releases section](https://github.com/src-d/gitbase-web/releases), and run it:

```bash
$ export GITBASEPG_DB_CONNECTION="root@tcp(<gitbase-ip>:3306)/none"
$ export GITBASEPG_BBLFSH_SERVER_URL="<bblfshd-ip>:9432"
$ ./gitbase-web serve
```

# Configuration

Any of the previous execution methods accept configuration through the following environment variables or CLI arguments.

| Variable | Argument | Default value | Meaning |
| -- | -- | -- | -- |
| `GITBASEPG_HOST` | `--host` | `0.0.0.0` | IP address to bind the HTTP server |
| `GITBASEPG_PORT` | `--port` | `8080` | Port to bind the HTTP server |
| `GITBASEPG_SERVER_URL` | `--server` | | URL used to access the application in the form `HOSTNAME[:PORT]`. Leave it unset to allow connections from any proxy or public address |
| `GITBASEPG_DB_CONNECTION` | `--db` | `root@tcp(localhost:3306)/none?maxAllowedPacket=4194304` | gitbase connection string. Use the DSN (Data Source Name) format described in the [Go MySQL Driver docs](https://github.com/go-sql-driver/mysql#dsn-data-source-name). |
| `GITBASEPG_CONN_MAX_LIFETIME` | `--conn-max-lifetime` | `30` | Maximum amount of time a SQL connection may be reused, in seconds. Make sure this value is lower than the timeout configured in the gitbase server, set with [`GITBASE_CONNECTION_TIMEOUT`](https://docs.sourced.tech/gitbase/using-gitbase/configuration#environment-variables) |
| `GITBASEPG_BBLFSH_SERVER_URL` | `--bblfsh` | `127.0.0.1:9432` | Address where bblfsh server is listening |
| `GITBASEPG_SELECT_LIMIT` | `--select-limit` | `100` | Default `LIMIT` forced on all the SQL queries done from the UI. Set it to 0 to remove any limit |
| `GITBASEPG_FOOTER_HTML` | `--footer` | | Allows to add any custom html to the page footer. It must be a string encoded in base64. Use it, for example, to add your analytics tracking code snippet  |
| `LOG_LEVEL` | `--log-level=`  | `info` | Logging level (`info`, `debug`, `warning` or `error`) |
| `LOG_FORMAT` | `--log-format=`  |  | log format (`text` or `json`), defaults to `text` on a terminal and `json` otherwise |
| `LOG_FIELDS` | `--log-fields=`  |  | default fields for the logger, specified in json |
| `LOG_FORCE_FORMAT` | `--log-force-format` | | ignore if it is running on a terminal or not |


# Contribute

[Contributions](https://github.com/src-d/gitbase-web/issues) are more than welcome, if you are interested please take a look to our [Contributing Guidelines](docs/CONTRIBUTING.md). There you will also find information on how to build and run the project from sources.

# Code of Conduct

All activities under source{d} projects are governed by the [source{d} code of conduct](https://github.com/src-d/guide/blob/master/.github/CODE_OF_CONDUCT.md).

# License

GPL v3.0, see [LICENSE](LICENSE)
