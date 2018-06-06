package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type SchemaSuite struct {
	HandlerUnitSuite
}

// Tests
// -----------------------------------------------------------------------------

func TestSchemaSuite(t *testing.T) {
	s := new(SchemaSuite)
	s.requestProcessFunc = Schema

	suite.Run(t, s)
}

func (suite *SchemaSuite) TestGet() {
	suite.mock.ExpectQuery("SHOW TABLES").WillReturnRows(
		sqlmock.NewRows([]string{"table"}).
			AddRow("foo"))

	suite.mock.ExpectQuery("DESCRIBE TABLE foo").WillReturnRows(
		sqlmock.NewRows([]string{"name", "type"}).
			AddRow("foo", "TEXT").
			AddRow("bar", "TEXT"))

	req, _ := http.NewRequest("GET", "/schema", nil)

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)
}
