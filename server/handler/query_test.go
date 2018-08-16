package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pressly/lg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type QuerySuite struct {
	HandlerUnitSuite
}

// Tests
// -----------------------------------------------------------------------------

func TestQuerySuite(t *testing.T) {
	s := new(QuerySuite)
	s.requestProcessFunc = Query

	suite.Run(t, s)
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
		{"select * from repositories limit 1", "select * from repositories limit 1"},
		{"select * from repositories limit 900", "select * from repositories LIMIT 100"},
		{"select * from repositories limit qwe", "select * from repositories limit qwe LIMIT 100"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc[0], func(t *testing.T) {
			a := assert.New(t)
			result, _ := addLimit(tc[0], 100)
			a.Equal(tc[1], result)
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

func (suite *QuerySuite) TestTypes() {
	columnNames := []string{"a", "b", "c", "d"}
	columnTypes := []string{"BIT", "INT", "DOUBLE", "TEXT"}

	columnValsPtr := genericVals(columnTypes)

	mockRows := sqlmock.NewRows([]string{"a", "b", "c", "d"}).
		AddRow(1, 1234, 1.56, "value").
		AddRow(nil, nil, nil, nil)

	suite.mock.ExpectQuery(".*").WillReturnRows(mockRows)

	rows, err := suite.db.Query("select * from table")
	suite.NoError(err)

	rows.Next()
	err = rows.Scan(columnValsPtr...)
	suite.NoError(err)

	colData, err := columnsData(columnNames, columnTypes, columnValsPtr)
	suite.NoError(err)

	suite.EqualValues(true, colData["a"])
	suite.EqualValues(1234, colData["b"])
	suite.EqualValues(1.56, colData["c"])
	suite.EqualValues("value", colData["d"])

	rows.Next()
	err = rows.Scan(columnValsPtr...)
	suite.NoError(err)

	colData, err = columnsData(columnNames, columnTypes, columnValsPtr)
	suite.NoError(err)

	suite.Nil(colData["a"])
	suite.Nil(colData["b"])
	suite.Nil(colData["c"])
	suite.Nil(colData["d"])
}

func (suite *QuerySuite) TestQueryAbort() {
	// Ideally we would test that the sql query context is canceled, but
	// go-sqlmock does not have something like ExpectContextCancellation

	mockRows := sqlmock.NewRows([]string{"a", "b", "c", "d"})
	suite.mock.ExpectQuery(".*").WillDelayFor(2 * time.Second).WillReturnRows(mockRows)

	json := `{"query": "select * from repositories"}`
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(json))
	res := httptest.NewRecorder()

	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)

	var mockAPIHandlerFunc http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		_, err := suite.requestProcessFunc(suite.db)(r)
		suite.Error(err)
		suite.Equal(context.Canceled, err)
	}

	go func() {
		handler := lg.RequestLogger(suite.logger)(mockAPIHandlerFunc)
		handler.ServeHTTP(res, req)
	}()

	cancel()

	wg.Wait()

	suite.Equal(context.Canceled, ctx.Err())
}
