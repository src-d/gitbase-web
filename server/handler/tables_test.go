package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/service"

	"github.com/pressly/lg"
	"github.com/stretchr/testify/suite"
)

// Suite setup
// -----------------------------------------------------------------------------

type TablesSuite struct {
	suite.Suite
	db      service.SQLDB
	handler http.Handler
}

func (suite *TablesSuite) SetupSuite() {
	suite.db = setupDB(suite.Require())

	// logger
	logger := logrus.New()

	// handler
	tablesHandler := handler.APIHandlerFunc(handler.Tables(suite.db))
	suite.handler = lg.RequestLogger(logger)(tablesHandler)
}

func (suite *TablesSuite) TearDownSuite() {
	suite.db.Close()
}

// Tests
// -----------------------------------------------------------------------------

func (suite *TablesSuite) TestGet() {
	req, _ := http.NewRequest("GET", "/tables", nil)

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	okResponse(suite.Require(), res)

	firstRow := firstRow(suite.Require(), res)
	suite.IsType("string", firstRow["table"])
}

// Main test to run the suite

func TestTablesSuite(t *testing.T) {
	suite.Run(t, new(TablesSuite))
}
