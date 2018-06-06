package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/serializer"
	"github.com/stretchr/testify/suite"
)

type SchemaIntegrationSuite struct {
	HandlerSuite
}

// Tests
// -----------------------------------------------------------------------------

func TestSchemaIntegrationSuite(t *testing.T) {
	q := new(SchemaIntegrationSuite)
	q.requestProcessFunc = handler.Schema

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, q)
}

func (suite *SchemaIntegrationSuite) TestGet() {
	req, _ := http.NewRequest("GET", "/schema", nil)

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)

	var resBody serializer.Response
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	suite.Nil(err)

	suite.Equal(res.Code, resBody.Status)
	suite.NotEmpty(resBody.Data)

	firstRow := firstRow(suite.Require(), res)
	suite.IsType("string", firstRow["table"])

	var interfaceSlice []interface{}
	suite.IsType(interfaceSlice, firstRow["columns"])
}
