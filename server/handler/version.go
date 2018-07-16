package handler

import (
	"net/http"

	bblfsh "gopkg.in/bblfsh/client-go.v2"

	"github.com/src-d/gitbase-playground/server/serializer"
	"github.com/src-d/gitbase-playground/server/service"
)

// Version returns a function that returns a *serializer.Response
// with a current version of server and dependencies
func Version(version, bbblfshServerURL string, db service.SQLDB) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		cli, err := bblfsh.NewClient(bbblfshServerURL)
		if err != nil {
			return nil, err
		}

		resp, err := cli.NewVersionRequest().Do()
		if err != nil {
			return nil, err
		}

		// old versions of gitbase don't have VERSION() function
		// so we set it to undefined and ignore error
		gitbaseVersion := "undefined"
		db.QueryRow("SELECT VERSION()").Scan(&gitbaseVersion)

		return serializer.NewVersionResponse(version, resp.Version, gitbaseVersion), nil
	}
}
