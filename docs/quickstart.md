# Quickstart

You can locally build and deploy `gitbase-playground` and its dependencies using [`docker-compose`](https://docs.docker.com/compose/install/)

If you preffer to run `gitbase-playground` dependencies manually, you can follow [the alternative playground quickstart](quickstart-manually.md)

## Populate the database

Populate a directory with some git repositories to be served by [gitbase](https://github.com/src-d/gitbase):

```bash
$ git clone git@github.com:src-d/gitbase-playground.git ./repos/gitbase-playground
$ git clone git@github.com:src-d/go-git-fixtures.git ./repos/go-git-fixtures
```

## Run the application

Once bblfsh and gitbase are running and accessible, you can serve the playground:

```bash
$ GITBASEPG_ENV=dev REPOS_FOLDER=./repos GITBASEPG_ENV=dev docker-compose up
```

Once the server is running &ndash;with its default values&ndash;, it will be accessible through: http://localhost:8080

You have more information about the [playground architecture](CONTRIBUTING.md#architecture), [development guides](CONTRIBUTING.md#development) and [configuration options](CONTRIBUTING.md#configuration) in the [CONTRIBUTING.md](CONTRIBUTING.md).


## Run a Query

You will find more info about how to run queries using the playground API on the [rest-api guide](rest-api.md)
