## Development

### Dependencies

Launch [bblfshd](https://github.com/bblfsh/bblfshd) and install the drivers. More info in the [bblfshd documentation](https://doc.bblf.sh/user/getting-started.html).

```bash
docker run -d --name bblfshd --privileged -p 9432:9432 -v /var/lib/bblfshd:/var/lib/bblfshd bblfsh/bblfshd
docker exec -it bblfshd bblfshctl driver install --all
```

Install [gitbase](https://github.com/src-d/gitbase), populate a repository directory, and start it.

```bash
go get github.com/src-d/gitbase/...
cd $GOPATH/src/github.com/src-d/gitbase
make dependencies
mkdir repos
git clone https://github.com/src-d/gitbase-playground.git repos/gitbase-playground
go run cli/gitbase/main.go server -v --git=repos
```

### Build

```bash
go build -o gitbase-playground cmd/server/main.go
```

### Run

Use `GITBASEPG_ENV=dev` for extra logs information.

Development:

```bash
GITBASEPG_ENV=dev go run cmd/server/main.go
```

Built binary:

```bash
GITBASEPG_ENV=dev ./gitbase-playground
```

### Run the Tests

```bash
go test -v server/handler/*
```
