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

	"github.com/src-d/gitbase-web/server/handler"
	"github.com/src-d/gitbase-web/server/serializer"
)

type UASTParseSuite struct {
	suite.Suite
	handler http.Handler
}

func TestUASTParseSuite(t *testing.T) {
	q := new(UASTParseSuite)
	q.handler = lg.RequestLogger(logrus.New())(handler.APIHandlerFunc(handler.Parse(bblfshServerURL())))

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, q)
}

func (suite *UASTParseSuite) TestSuccess() {
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

func (suite *UASTParseSuite) TestError() {
	jsonRequest := `{ "content": "function() { not_python = 1 }", "language": "python" }`
	req, _ := http.NewRequest("POST", "/parse", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

type UASTFilterSuite struct {
	suite.Suite
	handler http.Handler
}

func TestUASTFilterSuite(t *testing.T) {
	q := new(UASTFilterSuite)
	q.handler = lg.RequestLogger(logrus.New())(handler.APIHandlerFunc(handler.Filter()))

	suite.Run(t, q)
}

func (suite *UASTFilterSuite) TestSuccess() {
	jsonRequest := `{ "protobufs": "` + uastProtoMsgBase64List + `", "filter": "//*" }`
	req, _ := http.NewRequest("POST", "/filter", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Require().Equal(http.StatusOK, res.Code, res.Body.String())

	var resBody serializer.Response
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	suite.Nil(err)

	suite.Equal(res.Code, resBody.Status)
	suite.NotEmpty(resBody.Data)
}

func (suite *UASTFilterSuite) TestProtobufError() {
	jsonRequest := `{ "protobufs": "not-proto", "filter": "[" }`
	req, _ := http.NewRequest("POST", "/filter", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *UASTFilterSuite) TestFilterError() {
	jsonRequest := `{ "protobufs": "` + uastProtoMsgBase64List + `", "filter": "[" }`
	req, _ := http.NewRequest("POST", "/filter", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

// JSON: [<UAST(console.log("test"))>]
// Easy to obtain in the frontend with SELECT UAST('console.log("test")', 'JavaScript') AS uast
const uastProtoMsgBase64List = "AAACFgoERmlsZRr8AwoHUHJvZ3JhbRIXCgxpbnRlcm5hbFJvbGUSB3Byb2dyYW0SFAoKc291cmNlVHlwZRIGbW9kdWxlGrADChNFeHByZXNzaW9uU3RhdGVtZW50EhQKDGludGVybmFsUm9sZRIEYm9keRrxAgoOQ2FsbEV4cHJlc3Npb24SGgoMaW50ZXJuYWxSb2xlEgpleHByZXNzaW9uGtwBChBNZW1iZXJFeHByZXNzaW9uEhYKDGludGVybmFsUm9sZRIGY2FsbGVlEhEKCGNvbXB1dGVkEgVmYWxzZRpDCgpJZGVudGlmaWVyEg8KBE5hbWUSB2NvbnNvbGUSFgoMaW50ZXJuYWxSb2xlEgZvYmplY3QqBBABGAEyBggHEAEYCBpDCgpJZGVudGlmaWVyEhgKDGludGVybmFsUm9sZRIIcHJvcGVydHkSCwoETmFtZRIDbG9nKgYICBABGAkyBggLEAEYDCoEEAEYATIGCAsQARgMOgUCEgFUVRpSCgZTdHJpbmcSGQoMaW50ZXJuYWxSb2xlEglhcmd1bWVudHMSCgoGRm9ybWF0EgASDQoFVmFsdWUSBHRlc3QqBggMEAEYDTIGCBIQARgTOgJUMSoEEAEYATIGCBMQARgUOgISVCoEEAEYATIGCBMQARgUOgETKgQQARgBMgYIExABGBQ6ATkqBBABGAEyBggTEAEYFDoBIg=="
