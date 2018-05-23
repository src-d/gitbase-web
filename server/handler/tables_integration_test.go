package handler_test

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/serializer"
	"github.com/stretchr/testify/suite"
)

type TablesSuite struct {
	HandlerSuite
}

// Tests
// -----------------------------------------------------------------------------

func TestTablesSuite(t *testing.T) {
	flag.Parse()
	if !*gitbase {
		return
	}
	q := new(TablesSuite)
	q.requestProcessFunc = handler.Tables
	suite.Run(t, q)
}

func (suite *TablesSuite) TestGet() {
	req, _ := http.NewRequest("GET", "/tables", nil)

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
