package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type ExportSuite struct {
	suite.Suite
	db      service.SQLDB
	mock    sqlmock.Sqlmock
	handler http.HandlerFunc
}

func (suite *ExportSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to initialize the mock DB. '%s'", err)
	}

	suite.handler = handler.Export(suite.db)
}

func (suite *ExportSuite) TearDownTest() {
	suite.db.Close()
}

// Tests
// -----------------------------------------------------------------------------

func TestExportSuite(t *testing.T) {
	suite.Run(t, new(ExportSuite))
}

func (suite *ExportSuite) TestSuccess() {
	rows := sqlmock.NewRows([]string{"a", "b"}).
		AddRow(1, "one")

	suite.mock.ExpectQuery(".*").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/export/?query=select+*+from+repositories", nil)
	res := httptest.NewRecorder()

	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)
}

func (suite *ExportSuite) TestDBError() {
	suite.mock.ExpectQuery(".*").WillReturnError(fmt.Errorf("forced err"))

	req, _ := http.NewRequest("GET", "/export/?query=select+*+from+not_exist", nil)
	res := httptest.NewRecorder()

	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *ExportSuite) TestBadRequest() {
	testCases := []string{
		"/export/?query",
		"/export/?query=",
		"/export",
		"/export/?foo=bar",
	}

	suite.mock.ExpectQuery(".*").WillReturnError(fmt.Errorf("forced err"))

	for _, tc := range testCases {
		suite.T().Run(tc, func(t *testing.T) {
			a := assert.New(t)

			req, _ := http.NewRequest("GET", tc, nil)
			res := httptest.NewRecorder()

			suite.handler.ServeHTTP(res, req)

			a.Equal(http.StatusBadRequest, res.Code)
			a.Contains(res.Body.String(), "Bad Request")
		})
	}
}
