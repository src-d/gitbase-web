package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/service"

	"github.com/stretchr/testify/suite"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type SchemaSuite struct {
	suite.Suite
	db      service.SQLDB
	mock    sqlmock.Sqlmock
	handler http.Handler
}

func (suite *SchemaSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to initialize the mock DB. '%s'", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	h := handler.APIHandlerFunc(handler.Schema(suite.db))
	suite.handler = lg.RequestLogger(logger)(h)
}

func (suite *SchemaSuite) TearDownTest() {
	suite.db.Close()
}

// Tests
// -----------------------------------------------------------------------------

func TestSchemaSuite(t *testing.T) {
	suite.Run(t, new(SchemaSuite))
}

func (suite *SchemaSuite) TestGet() {
	suite.mock.ExpectQuery("SHOW TABLES").WillReturnRows(
		sqlmock.NewRows([]string{"table"}).
			AddRow("foo"))

	suite.mock.ExpectQuery("DESCRIBE TABLE.*").WillReturnRows(
		sqlmock.NewRows([]string{"name", "type"}).
			AddRow("foo", "TEXT").
			AddRow("bar", "TEXT"))

	req, _ := http.NewRequest("GET", "/schema", nil)

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)
}
