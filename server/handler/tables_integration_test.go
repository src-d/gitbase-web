package handler_test

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/stretchr/testify/suite"
)

type TablesSuite struct {
	HandlerSuite
}

// Tests
// -----------------------------------------------------------------------------

func TestTablesSuite(t *testing.T) {
	flag.Parse()
	if !*gitbase {
		return
	}
	q := new(TablesSuite)
	q.requestProcessFunc = handler.Tables
	suite.Run(t, q)
}

func (suite *TablesSuite) TestGet() {
	req, _ := http.NewRequest("GET", "/tables", nil)

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	okResponse(suite.Require(), res)

	firstRow := firstRow(suite.Require(), res)
	suite.IsType("string", firstRow["table"])
}
