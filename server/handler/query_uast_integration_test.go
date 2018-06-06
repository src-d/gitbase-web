package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/stretchr/testify/suite"
)

type QueryUast struct {
	HandlerSuite
}

// Tests
// -----------------------------------------------------------------------------

func TestUastFunctions(t *testing.T) {
	q := new(QueryUast)
	q.requestProcessFunc = handler.Query

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, q)
}

// This test requires that gitbase can reach bblfshd and that it's serving the
// repository https://github.com/src-d/gitbase-playground
func (suite *QueryUast) TestUastFunctions() {
	req, _ := http.NewRequest("POST", "/query", strings.NewReader(
		`{ "query": "SELECT blob_hash, blob_content, uast(blob_content, 'go') as uast FROM blobs WHERE blob_hash='fd30cea52792da5ece9156eea4022bdd87565633'" }`))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	okResponse(suite.Require(), res)

	firstRow := firstRow(suite.Require(), res)
	suite.IsType("string", firstRow["blob_hash"])
	suite.IsType("string", firstRow["blob_content"])

	var arr []interface{}
	suite.IsType(arr, firstRow["uast"])

	uasts := firstRow["uast"].([]interface{})

	var jsonObj map[string]interface{}
	suite.IsType(jsonObj, uasts[0])

	// check that roles were converted correctly
	uast := uasts[0].(map[string]interface{})
	roles := uast["Roles"].([]interface{})
	suite.Equal("File", roles[1].(string))
}
