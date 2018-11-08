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
	jsonRequest := `{ "content": "function(} ][", "language": "javascript" }`
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

type UASTModeSuite struct {
	suite.Suite
	handler http.Handler
}

func TestUASTModeSuite(t *testing.T) {
	q := new(UASTModeSuite)
	q.handler = lg.RequestLogger(logrus.New())(handler.APIHandlerFunc(handler.Parse(bblfshServerURL())))

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, q)
}

func (suite *UASTModeSuite) TestSuccess() {
	testCases := []string{
		`{ "content": "console.log('test')", "language": "javascript", "mode": "" }`,
		`{ "content": "console.log('test')", "language": "javascript", "mode": "native" }`,
		`{ "content": "console.log('test')", "language": "javascript", "mode": "annotated" }`,
		`{ "content": "console.log('test')", "language": "javascript", "mode": "semantic" }`,
	}

	for _, tc := range testCases {
		suite.T().Run(tc, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/parse", strings.NewReader(tc))

			res := httptest.NewRecorder()
			suite.handler.ServeHTTP(res, req)

			suite.Require().Equal(http.StatusOK, res.Code, res.Body.String())

			var resBody serializer.Response
			err := json.Unmarshal(res.Body.Bytes(), &resBody)
			suite.Nil(err)

			suite.Equal(res.Code, resBody.Status)
			suite.NotEmpty(resBody.Data)
		})
	}
}

func (suite *UASTModeSuite) TestWrongMode() {
	jsonRequest := `{ "content": "console.log('test')", "language": "javascript", "mode": "foo" }`
	req, _ := http.NewRequest("POST", "/parse", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

// JSON: [<UAST(console.log("test"))>]
// Easy to obtain in the frontend with SELECT UAST('console.log("test")', 'JavaScript') AS uast
// Gitbase v0.18.0-beta.1, Bblfsh v2.9.2-drivers
const uastProtoMsgBase64List = "AGJncgEAAAAECFcQAQNCAQIOOgUDBAUGB0IFCBYXGBkGEgRAcG9zBxIFQHJvbGUHEgVAdHlwZQoSCGNvbW1lbnRzCRIHcHJvZ3JhbQo6AwUJCkIDCwwUBRIDZW5kBxIFc3RhcnQQEg51YXN0OlBvc2l0aW9ucww6BAUNDg9CBBAREhMFEgNjb2wGEgRsaW5lCBIGb2Zmc2V0DxINdWFzdDpQb3NpdGlvbgIgFAIgAQIgEwhCBBASEhVQDAIgAANCARcGEgRGaWxlABA6BgMEBRobHEIGCB0fIBhWBhIEYm9keQwSCmRpcmVjdGl2ZXMMEgpzb3VyY2VUeXBlA0IBHggSBk1vZHVsZQkSB1Byb2dyYW0DQgEhDDoEAwQFIkIECCMlJgwSCmV4cHJlc3Npb24DQgEkCxIJU3RhdGVtZW50FRITRXhwcmVzc2lvblN0YXRlbWVudA46BQMEBScoQgUIKSwtPAsSCWFyZ3VtZW50cwgSBmNhbGxlZQRCAiorDBIKRXhwcmVzc2lvbgYSBENhbGwQEg5DYWxsRXhwcmVzc2lvbgNCAS4OOgUDBAUvMEIFMTc5OjsIEgZGb3JtYXQHEgVWYWx1ZQdCAwsyNFAICEIEEBMSM1AMAiASCEIEEDUSNlAMAiANAiAMBEICKzgKEghBcmd1bWVudA0SC3Vhc3Q6U3RyaW5nAhIABhIEdGVzdBA6BgMEBT0+P0IGQENHSElRChIIY29tcHV0ZWQIEgZvYmplY3QKEghwcm9wZXJ0eQdCAwtBFFAICEIEEDYSQlAMAiALB0IFRCpFK0YLEglRdWFsaWZpZWQMEgpJZGVudGlmaWVyCBIGQ2FsbGVlEhIQTWVtYmVyRXhwcmVzc2lvbgIwAAo6AwMFSkIDS09QBhIETmFtZQdCAwtMFFAICEIEEE0STlAMAiAIAiAHERIPdWFzdDpJZGVudGlmaWVyCRIHY29uc29sZQdCA1JPVVBJB0IDC0FTUAgIQgQQVBJNUAwCIAkFEgNsb2cIEgZtb2R1bGU="
