package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"github.com/src-d/gitbase-playground/server/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type QuerySuite struct {
	suite.Suite
	db      service.SQLDB
	mock    sqlmock.Sqlmock
	handler http.Handler
}

func (suite *QuerySuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to initialize the mock DB. '%s'", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	h := APIHandlerFunc(Query(suite.db))
	suite.handler = lg.RequestLogger(logger)(h)
}

func (suite *QuerySuite) TearDownTest() {
	suite.db.Close()
}

// Tests
// -----------------------------------------------------------------------------

func TestQuerySuite(t *testing.T) {
	suite.Run(t, new(QuerySuite))
}

func (suite *QuerySuite) TestAddLimit() {
	testCases := [][]string{
		{"SHOW TABLES", "SHOW TABLES"},
		{"select * from repositories", "select * from repositories LIMIT 100"},
		{"SELECT * FROM repositories", "SELECT * FROM repositories LIMIT 100"},
		{`
			SELECT * FROM repositories
			`, "SELECT * FROM repositories LIMIT 100"},
		{"  SELECT * FROM repositories  ", "SELECT * FROM repositories LIMIT 100"},
		{"  SELECT * FROM repositories  ; ", "SELECT * FROM repositories   LIMIT 100"},
		{"/* comment */ SELECT * FROM repositories", "SELECT * FROM repositories LIMIT 100"},
		{"SELECT * FROM repositories /* comment */", "SELECT * FROM repositories LIMIT 100"},
		{"SELECT * FROM repositories; /* comment */", "SELECT * FROM repositories LIMIT 100"},
		{`/* comment
			multiline */ SELECT * FROM repositories; /* comment
			multiline */`, "SELECT * FROM repositories LIMIT 100"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc[0], func(t *testing.T) {
			a := assert.New(t)
			a.Equal(tc[1], addLimit(tc[0], 100))
		})
	}
}

func (suite *QuerySuite) TestBadRequest() {
	testCases := []string{
		`{"wrongname": "select * from repositories"}`,
		`name": "select * from repositories"}`,
		`{"query": 1234}`,
		`{"query": "select * from repositories", "limit": "string"}`,
	}

	suite.mock.ExpectQuery(".*").WillReturnError(fmt.Errorf("forced err"))

	for _, tc := range testCases {
		suite.T().Run(tc, func(t *testing.T) {
			a := assert.New(t)

			req, _ := http.NewRequest("POST", "/query", strings.NewReader(tc))
			res := httptest.NewRecorder()
			suite.handler.ServeHTTP(res, req)

			a.Equal(http.StatusBadRequest, res.Code)
			a.Contains(res.Body.String(), "Bad Request")
		})
	}
}

func (suite *QuerySuite) TestQueryErr() {
	suite.mock.ExpectQuery(".*").WillReturnError(fmt.Errorf("forced err"))

	json := `{"query": "select * from repositories"}`
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(json))
	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *QuerySuite) TestQuery() {
	rows := sqlmock.NewRows([]string{"a", "b", "c", "d"}).
		AddRow(1, "one", 1.5, 100).
		AddRow(nil, nil, nil, nil)

	suite.mock.ExpectQuery(".*").WillReturnRows(rows)

	json := `{"query": "select * from repositories"}`
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(json))
	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)
}
