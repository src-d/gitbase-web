package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/serializer"
)

type UASTSuite struct {
	suite.Suite
	handler http.Handler
}

func TestUASTSuite(t *testing.T) {
	q := new(UASTSuite)
	q.handler = lg.RequestLogger(logrus.New())(handler.APIHandlerFunc(handler.Parse(bblfshServerURL())))

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, q)
}

func (suite *UASTSuite) TestSuccess() {
	jsonRequest := `{ "content": "console.log('test')", "language": "javascript" }`
	req, _ := http.NewRequest("POST", "/parse", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Require().Equal(http.StatusOK, res.Code, res.Body.String())

	var resBody serializer.Response
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	suite.Nil(err)

	suite.Equal(res.Code, resBody.Status)
	suite.NotEmpty(resBody.Data)
}

func (suite *UASTSuite) TestError() {
	jsonRequest := `{ "content": "function() { not_python = 1 }", "language": "python" }`
	req, _ := http.NewRequest("POST", "/parse", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusInternalServerError, res.Code)
}
