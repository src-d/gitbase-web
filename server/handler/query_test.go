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

	common "github.com/src-d/gitbase-web/server/testing"

	"github.com/pressly/lg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gopkg.in/bblfsh/sdk.v2/uast/nodes"
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
		{"  SELECT * FROM repositories  ; ", "SELECT * FROM repositories LIMIT 100"},
		{`  SELECT * FROM repositories  
; `, "SELECT * FROM repositories LIMIT 100"},
		{"/* comment */ SELECT * FROM repositories", "SELECT * FROM repositories LIMIT 100"},
		{"SELECT * FROM repositories /* comment */", "SELECT * FROM repositories LIMIT 100"},
		{"SELECT * FROM repositories; /* comment */", "SELECT * FROM repositories LIMIT 100"},
		{`/* comment
			multiline */ SELECT * FROM repositories; /* comment
			multiline */`, "SELECT * FROM repositories LIMIT 100"},
		{"select * from repositories limit 1", "select * from repositories limit 1"},
		{"select * from repositories limit 1;", "select * from repositories limit 1"},
		{"select * from repositories limit 1 ;", "select * from repositories limit 1"},
		{`select * from repositories limit 1
;`, "select * from repositories limit 1"},
		{`select * from repositories limit 1 
 ; `, "select * from repositories limit 1"},
		{"select * from repositories limit 900", "select * from repositories LIMIT 100"},
		{"select * from repositories limit 900;", "select * from repositories LIMIT 100"},
		{"select * from repositories limit 900 ; ", "select * from repositories LIMIT 100"},
		{`select * from repositories limit 900
 ; `, "select * from repositories LIMIT 100"},
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

	mockRows := sqlmock.NewRows(columnNames).
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

func (suite *QuerySuite) TestUASTType() {
	columnNames := []string{"filename", "uast_a", "uast_b"}
	columnTypes := []string{"TEXT", "TEXT", "TEXT"}

	columnValsPtr := genericVals(columnTypes)

	mockRows := sqlmock.NewRows(columnNames).
		AddRow("hello.js", "", common.UASTMarshaled)

	suite.mock.ExpectQuery(".*").WillReturnRows(mockRows)

	rows, err := suite.db.Query("select * from table")
	suite.NoError(err)

	rows.Next()
	err = rows.Scan(columnValsPtr...)
	suite.NoError(err)

	colData, err := columnsData(columnNames, columnTypes, columnValsPtr)
	suite.NoError(err)

	suite.EqualValues("hello.js", colData["filename"])
	suite.Nil(colData["__filename-protobufs"])

	suite.EqualValues("", colData["uast_a"])
	suite.Nil(colData["__uast_a-protobufs"])

	var nodeArr nodes.Array
	suite.IsType(nodeArr, colData["uast_b"])
	suite.EqualValues(common.UASTMarshaled, colData["__uast_b-protobufs"])
}

func (suite *QuerySuite) TestQueryAbort() {
	// Ideally we would test that the sql query context is canceled, but
	// go-sqlmock does not have something like ExpectContextCancellation

	mockRows := sqlmock.NewRows([]string{"a", "b", "c", "d"}).AddRow(1, "one", 1.5, 100)
	suite.mock.ExpectQuery(`select \* from repositories`).WillDelayFor(2 * time.Second).WillReturnRows(mockRows)

	mockProcessRows := sqlmock.NewRows(
		[]string{"Id", "User", "Host", "db", "Command", "Time", "State", "Info"}).
		AddRow(1234, nil, "localhost:3306", nil, "query", 2, "SquashedTable(refs, commit_files, files)(1/5)", "select * from files").
		AddRow(1288, nil, "localhost:3306", nil, "query", 2, "SquashedTable(refs, commit_files, files)(1/5)", "select * from repositories")
	suite.mock.ExpectQuery("SHOW FULL PROCESSLIST").WillReturnRows(mockProcessRows)

	suite.mock.ExpectExec("KILL 1288")

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

	// Without this wait the Request is cancelled before the handler has time to
	// start the query. Which also works fine, but we want to test a cancellation
	// for a query that is in progress
	time.Sleep(200 * time.Millisecond)

	cancel()

	wg.Wait()

	suite.Equal(context.Canceled, ctx.Err())
}
