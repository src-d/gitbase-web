package handler_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type QuerySuite struct {
	HandlerSuite
}

// Tests
// -----------------------------------------------------------------------------

func TestQuerySuite(t *testing.T) {
	flag.Parse()
	if !*gitbase {
		return
	}
	q := new(QuerySuite)
	q.requestProcessFunc = handler.Query
	suite.Run(t, q)
}

func (suite *QuerySuite) TestSelectAll() {
	testCases := []string{
		"blobs",
		"commits",
		"refs",
		"remotes",
		"repositories",
		"tree_entries",
	}

	for _, tc := range testCases {
		suite.T().Run(tc, func(t *testing.T) {
			jsonRequest := fmt.Sprintf(`{ "query": "select * from %s", "limit": 100 }`, tc)
			req, _ := http.NewRequest("POST", "/query", strings.NewReader(jsonRequest))

			res := httptest.NewRecorder()
			suite.handler.ServeHTTP(res, req)

			okResponse(require.New(t), res)
		})
	}
}

func (suite *QuerySuite) TestLimit() {
	testCases := []string{
		`{ "query": "select * from refs", "limit": 100 }`,
		`{ "query": "select * from refs", "limit": 0 }`,
		`{ "query": "select * from refs" }`,
	}

	for _, tc := range testCases {
		suite.T().Run(tc, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/query", strings.NewReader(tc))

			res := httptest.NewRecorder()
			suite.handler.ServeHTTP(res, req)

			okResponse(require.New(t), res)
		})
	}
}

func (suite *QuerySuite) TestBoolFunctions() {
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(
		`{ "query": "select name, is_remote(name) as remote, is_tag(name) as tag from refs" }`))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	okResponse(suite.Require(), res)

	firstRow := firstRow(suite.Require(), res)
	suite.IsType("string", firstRow["name"])
	suite.IsType(true, firstRow["remote"])
	suite.IsType(true, firstRow["tag"])
}

func (suite *QuerySuite) TestWrongSQLSyntax() {
	jsonRequest := `{ "query": "selectSELECT * from commits", "limit": 100 }`
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Require().Equal(http.StatusBadRequest, res.Code)

	resBody, err := errorResponse(res)
	suite.Require().Nil(err)
	suite.EqualValues(res.Code, resBody["status"])

	firstErr := firstErr(suite.Require(), resBody)
	suite.EqualValues(1105, firstErr["mysqlCode"])
	suite.EqualValues(res.Code, firstErr["status"])
	suite.Contains(firstErr["title"], "syntax error")
}

func (suite *QuerySuite) TestWrongLimit() {
	testCases := []string{
		`[1, 2]`,
		`"10"`,
		`{ "a" : 5 }`,
	}

	for _, tc := range testCases {
		suite.T().Run(tc, func(t *testing.T) {
			jsonRequest := fmt.Sprintf(`{ "query": "select * from commits", "limit": %s }`, tc)
			req, _ := http.NewRequest("POST", "/query", strings.NewReader(jsonRequest))

			res := httptest.NewRecorder()
			suite.handler.ServeHTTP(res, req)

			require := require.New(t)

			require.Equal(http.StatusBadRequest, res.Code)

			resBody, err := errorResponse(res)
			require.Nil(err)
			require.EqualValues(res.Code, resBody["status"])

			firstErr := firstErr(require, resBody)
			require.EqualValues(res.Code, firstErr["status"])
			require.Contains(firstErr["title"], "Bad Request")
		})
	}
}

// This test requires that gitbase can reach bblfshd and that it's serving the
// repository https://github.com/src-d/gitbase-playground
func (suite *QuerySuite) TestUastFunctions() {
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(
		`{ "query": "SELECT hash, content, uast(content, 'go') as uast FROM blobs WHERE hash='fd30cea52792da5ece9156eea4022bdd87565633'" }`))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	if false {
		okResponse(suite.Require(), res)

		firstRow := firstRow(suite.Require(), res)
		suite.IsType("string", firstRow["hash"])
		suite.IsType("string", firstRow["content"])

		var arr []interface{}
		suite.IsType(arr, firstRow["uast"])

		var jsonObj map[string]interface{}
		suite.IsType(jsonObj, firstRow["uast"].([]interface{})[0])
	}
}
