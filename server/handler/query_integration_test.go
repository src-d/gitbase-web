package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type QueryIntegrationSuite struct {
	HandlerSuite
}

// Tests
// -----------------------------------------------------------------------------

func TestQueryIntegrationSuite(t *testing.T) {
	q := new(QueryIntegrationSuite)
	q.requestProcessFunc = handler.Query

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, q)
}

func (suite *QueryIntegrationSuite) TestSelectAll() {
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

func (suite *QueryIntegrationSuite) TestLimit() {
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

func (suite *QueryIntegrationSuite) TestBoolFunctions() {
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(
		`{ "query": "select ref_name, is_remote(ref_name) as remote, is_tag(ref_name) as tag from refs" }`))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	okResponse(suite.Require(), res)

	firstRow := firstRow(suite.Require(), res)
	suite.IsType("string", firstRow["ref_name"])
	suite.IsType(true, firstRow["remote"])
	suite.IsType(true, firstRow["tag"])
}

func (suite *QueryIntegrationSuite) TestWrongSQLSyntax() {
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

func (suite *QueryIntegrationSuite) TestWrongLimit() {
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
