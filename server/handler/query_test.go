package handler_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/serializer"

	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Suite setup
// -----------------------------------------------------------------------------

type appConfig struct {
	DBConn string `envconfig:"DB_CONNECTION" default:"root@tcp(localhost:3306)/none?maxAllowedPacket=4194304"`
}

type QuerySuite struct {
	suite.Suite
	db      *sql.DB
	handler http.Handler
}

func setupDB(require *require.Assertions) *sql.DB {
	var conf appConfig
	envconfig.MustProcess("GITBASEPG", &conf)

	// db
	var err error
	db, err := sql.Open("mysql", conf.DBConn)
	require.Nil(err)

	err = db.Ping()
	require.Nil(err)

	return db
}

func (suite *QuerySuite) SetupSuite() {
	suite.db = setupDB(suite.Require())

	// logger
	logger := logrus.New()

	// handler
	queryHandler := handler.APIHandlerFunc(handler.Query(suite.db))
	suite.handler = lg.RequestLogger(logger)(queryHandler)
}

func (suite *QuerySuite) TearDownSuite() {
	suite.db.Close()
}

// Helpers
// -----------------------------------------------------------------------------

func errorResponse(res *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)

	return resBody, err
}

func firstErr(require *require.Assertions, resBody map[string]interface{}) map[string]interface{} {
	require.NotEmpty(resBody["errors"].([]interface{}))
	return resBody["errors"].([]interface{})[0].(map[string]interface{})
}

func firstRow(require *require.Assertions, res *httptest.ResponseRecorder) map[string]interface{} {
	var resBody serializer.Response
	json.Unmarshal(res.Body.Bytes(), &resBody)
	require.NotEmpty(resBody.Data.([]interface{}))
	return resBody.Data.([]interface{})[0].(map[string]interface{})
}

func okResponse(require *require.Assertions, res *httptest.ResponseRecorder) {
	require.Equal(http.StatusOK, res.Code)

	var resBody serializer.Response
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	require.Nil(err)

	require.Equal(res.Code, resBody.Status)
	require.NotEmpty(resBody.Data)
	require.NotEmpty(resBody.Meta)
}

// Tests
// -----------------------------------------------------------------------------

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

// This test requires that gitbase can reach bblfshd and that it's serving the
// repository https://github.com/src-d/gitbase-playground
func (suite *QuerySuite) TestUastFunctions() {
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(
		`{ "query": "SELECT hash, content, uast(content, 'go') as uast FROM blobs WHERE hash='fd30cea52792da5ece9156eea4022bdd87565633'" }`))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	okResponse(suite.Require(), res)

	firstRow := firstRow(suite.Require(), res)
	suite.IsType("string", firstRow["hash"])
	suite.IsType("string", firstRow["content"])

	var arr []interface{}
	suite.IsType(arr, firstRow["uast"])

	var jsonObj map[string]interface{}
	suite.IsType(jsonObj, firstRow["uast"].([]interface{})[0])
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

// Main test to run the suite

func TestQuerySuite(t *testing.T) {
	suite.Run(t, new(QuerySuite))
}
