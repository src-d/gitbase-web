package handler_test

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/service"

	"github.com/stretchr/testify/suite"
)

type ExportIntegrationSuite struct {
	suite.Suite
	db      service.SQLDB
	handler http.HandlerFunc
}

func TestExportIntegrationSuite(t *testing.T) {
	db, err := getDB()

	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	s := new(ExportIntegrationSuite)
	s.db = db
	s.handler = handler.Export(db)

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, s)
}

func (suite *ExportIntegrationSuite) TestSuccess() {
	req, _ := http.NewRequest("GET", "/export/?query=select+*+from+repositories", nil)
	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)

	r := csv.NewReader(res.Body)

	record, err := r.Read()
	suite.Nil(err)
	suite.Equal(record, []string{"repository_id"})

	record, err = r.Read()
	suite.Nil(err)
	suite.Equal(len(record), 1)
	suite.True(len(record[0]) > 0)
}

func (suite *ExportIntegrationSuite) TestError() {
	req, _ := http.NewRequest("GET", "/export/?query=select+*+from+not_exist", nil)
	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *ExportIntegrationSuite) TearDownSuite() {
	suite.db.Close()
}
