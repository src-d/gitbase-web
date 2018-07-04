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
const uastProtoMsgBase64List = "WyJDZ1JHYVd4bEd1a0RDZ2RRY205bmNtRnRFaFFLQ25OdmRYSmpaVlI1Y0dVU0JtMXZaSFZzWlJJWENneHBiblJsY201aGJGSnZiR1VTQjNCeWIyZHlZVzBhblFNS0UwVjRjSEpsYzNOcGIyNVRkR0YwWlcxbGJuUVNGQW9NYVc1MFpYSnVZV3hTYjJ4bEVnUmliMlI1R3Q0Q0NnNURZV3hzUlhod2NtVnpjMmx2YmhJYUNneHBiblJsY201aGJGSnZiR1VTQ21WNGNISmxjM05wYjI0YTFBRUtFRTFsYldKbGNrVjRjSEpsYzNOcGIyNFNGZ29NYVc1MFpYSnVZV3hTYjJ4bEVnWmpZV3hzWldVU0VRb0lZMjl0Y0hWMFpXUVNCV1poYkhObEdqOEtDa2xrWlc1MGFXWnBaWElTRmdvTWFXNTBaWEp1WVd4U2IyeGxFZ1p2WW1wbFkzUWlCMk52Ym5OdmJHVXFCQkFCR0FFeUJnZ0hFQUVZQ0RvQ0VnRWFQd29LU1dSbGJuUnBabWxsY2hJWUNneHBiblJsY201aGJGSnZiR1VTQ0hCeWIzQmxjblI1SWdOc2IyY3FCZ2dJRUFFWUNUSUdDQXNRQVJnTU9nSVNBU29FRUFFWUFUSUdDQXNRQVJnTU9nVUNFZ0ZVVlJwSENnMVRkSEpwYm1kTWFYUmxjbUZzRWhrS0RHbHVkR1Z5Ym1Gc1VtOXNaUklKWVhKbmRXMWxiblJ6SWdSMFpYTjBLZ1lJREJBQkdBMHlCZ2dTRUFFWUV6b0ZFbGhpVkRFcUJCQUJHQUV5QmdnVEVBRVlGRG9DRWxRcUJCQUJHQUV5QmdnVEVBRVlGRG9CRXlvRUVBRVlBVElHQ0JRUUFoZ0JPZ0U1S2dRUUFSZ0JNZ1lJRkJBQ0dBRTZBU0k9Il0="
