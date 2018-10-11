package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/src-d/gitbase-web/server"
	"github.com/src-d/gitbase-web/server/handler"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/src-d/go-cli.v0"
	"gopkg.in/src-d/go-log.v1"
)

// version will be replaced automatically by the CI build.
// See https://github.com/src-d/ci/blob/v1/Makefile.main#L56
var (
	name    = "gitbase-web"
	version = "undefined"
	build   = "undefined"
)

var app = cli.New(name, version, build, "gitbase web client")

// Note: maxAllowedPacket must be explicitly set for go-sql-driver/mysql v1.3.
// Otherwise gitbase will be asked for the max_allowed_packet column and the
// query will fail.
// The next release should make this parameter optional for us:
// https://github.com/go-sql-driver/mysql/pull/680
type ServeCommand struct {
	cli.PlainCommand `name:"serve" short-description:"serve the app" long-description:"starts serving the application"`
	cli.LogOptions   `group:"Log Options"`
	Host             string `long:"host" env:"GITBASEPG_HOST" default:"0.0.0.0" description:"IP address to bind the HTTP server"`
	Port             int    `long:"port" env:"GITBASEPG_PORT" default:"8080" description:"Port to bind the HTTP server"`
	ServerURL        string `long:"server" env:"GITBASEPG_SERVER_URL" description:"URL used to access the application in the form 'HOSTNAME[:PORT]'. Leave it unset to allow connections from any proxy or public address"`
	DBConn           string `long:"db" env:"GITBASEPG_DB_CONNECTION" default:"root@tcp(localhost:3306)/none?maxAllowedPacket=4194304" description:"gitbase connection string. Use the DSN (Data Source Name) format described in the Go MySQL Driver docs: https://github.com/go-sql-driver/mysql#dsn-data-source-name"`
	SelectLimit      int    `long:"select-limit" env:"GITBASEPG_SELECT_LIMIT" default:"100" description:"Default 'LIMIT' forced on all the SQL queries done from the UI. Set it to 0 to remove any limit"`
	BblfshServerURL  string `long:"bblfsh" env:"GITBASEPG_BBLFSH_SERVER_URL" default:"127.0.0.1:9432" description:"Address where bblfsh server is listening"`
	FooterHTML       string `long:"footer" env:"GITBASEPG_FOOTER_HTML" description:"Allows to add any custom html to the page footer. It must be a string encoded in base64. Use it, for example, to add your analytics tracking code snippet"`
}

func (c *ServeCommand) Execute(args []string) error {
	c.initLog()

	// database
	db, err := sql.Open("mysql", c.DBConn)
	if err != nil {
		return fmt.Errorf("error opening the database: %s", err.Error())
	}
	defer db.Close()

	db.SetMaxIdleConns(0)

	static := handler.NewStatic("build/public", c.ServerURL, c.SelectLimit, c.FooterHTML)

	// start the router
	router := server.Router(logrus.StandardLogger(), static, version, db, c.BblfshServerURL)

	log.With(log.Fields{"version": version, "build": build}).
		Infof("listening on %s:%d", c.Host, c.Port)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", c.Host, c.Port), router)
	log.Errorf(err, "")
	return err
}

func (c *ServeCommand) initLog() {
	if c.LogFields == "" {
		bytes, err := json.Marshal(log.Fields{"app": name})
		if err != nil {
			panic(err)
		}
		c.LogFields = string(bytes)
	}

	log.DefaultFactory = &log.LoggerFactory{
		Level:       c.LogLevel,
		Format:      c.LogFormat,
		Fields:      c.LogFields,
		ForceFormat: c.LogForceFormat,
	}
	log.DefaultFactory.ApplyToLogrus()

	log.DefaultLogger = log.New(nil)
}
func main() {
	app.AddCommand(&ServeCommand{})

	app.RunMain()
}
