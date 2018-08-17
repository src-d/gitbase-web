package handler

import (
	"net/http"

	bblfsh "gopkg.in/bblfsh/client-go.v2"

	"github.com/src-d/gitbase-web/server/serializer"
	"github.com/src-d/gitbase-web/server/service"
)

// Version returns a function that returns a *serializer.Response
// with a current version of server and dependencies
func Version(version, bbblfshServerURL string, db service.SQLDB) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		// old versions of gitbase don't have VERSION() function
		// so we set it to undefined and ignore error
		gitbaseVersion := "undefined"
		row := db.QueryRow("SELECT VERSION()")
		if row != nil {
			row.Scan(&gitbaseVersion)
		}

		// ignore bblfsh errors and return undefined to be consistent with gitbase
		bblfshVersion := "undefined"
		cli, err := bblfsh.NewClient(bbblfshServerURL)
		if err == nil {
			resp, err := cli.NewVersionRequest().Do()
			if err == nil && len(resp.Errors) == 0 {
				bblfshVersion = resp.Version
			}
		}

		return serializer.NewVersionResponse(version, bblfshVersion, gitbaseVersion), nil
	}
}
