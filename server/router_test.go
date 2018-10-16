package server_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/gitbase-web/server"
	"github.com/src-d/gitbase-web/server/handler"
	"github.com/src-d/gitbase-web/server/service"
	testingTools "github.com/src-d/gitbase-web/server/testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	log "gopkg.in/src-d/go-log.v1"
)

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

type RouterTestSuite struct {
	ClientTestSuite
	router http.Handler
	server *httptest.Server
	db     service.SQLDB
}

const version = "test-version"

func (s *RouterTestSuite) SetupSuite() {
	(&log.LoggerFactory{}).ApplyToLogrus()

	staticHandler := &handler.Static{}
	s.db = &testingTools.MockDB{}
	s.router = server.Router(
		logrus.StandardLogger(),
		staticHandler,
		version,
		s.db,
		"",
	)
}

func (s *RouterTestSuite) SetupTest() {
	s.server = httptest.NewServer(s.router)
}

func (s *RouterTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *RouterTestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *RouterTestSuite) GetResponse(method string, path string, body io.Reader) *http.Response {
	url := s.server.URL + path
	response, err := GetResponse(method, url, body)
	if err != nil {
		s.Fail(err.Error())
	}

	return response
}

func (s *RouterTestSuite) TestVersion() {
	expectedVersion := fmt.Sprintf(
		`{"status":200,"data":{"version":"%s","bblfsh":"undefined","gitbase":"undefined"}}`, version)
	response := s.GetResponse("GET", "/version", nil)
	s.AssertResponseBodyStatus(response, 200, expectedVersion, "version should be served")
}
