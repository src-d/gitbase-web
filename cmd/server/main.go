package main

import (
	"fmt"
	"net/http"

	"github.com/src-d/gitbase-playground/server"
	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/service"

	"github.com/kelseyhightower/envconfig"
)

// version will be replaced automatically by the CI build.
// See https://github.com/src-d/ci/blob/v1/Makefile.main#L56
var version = "dev"

type appConfig struct {
	Env       string `envconfig:"ENV" default:"production"`
	Host      string `envconfig:"HOST" default:"0.0.0.0"`
	Port      int    `envconfig:"PORT" default:"8080"`
	ServerURL string `envconfig:"SERVER_URL"`
}

func main() {
	// main configuration
	var conf appConfig
	envconfig.MustProcess("GITBASEPG", &conf)
	if conf.ServerURL == "" {
		conf.ServerURL = fmt.Sprintf("//localhost:%d", conf.Port)
	}

	// logger
	logger := service.NewLogger(conf.Env)

	static := handler.NewStatic("build", conf.ServerURL)

	// start the router
	router := server.Router(logger, static, version)
	logger.Infof("listening on %s:%d", conf.Host, conf.Port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router)
	logger.Fatal(err)
}
