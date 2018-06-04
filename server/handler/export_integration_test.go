package handler_test

import (
	"encoding/csv"
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/gitbase-playground/server/handler"
	"github.com/src-d/gitbase-playground/server/service"

	"github.com/stretchr/testify/suite"
)

type ExportSuite struct {
	suite.Suite
	db      service.SQLDB
	handler http.HandlerFunc
}

func TestExportSuite(t *testing.T) {
	flag.Parse()
	if !*gitbase {
		return
	}

	db, err := getDB(true)
	if err != nil {
		t.Fatal(db)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(db)
	}

	s := new(ExportSuite)
	s.db = db
	s.handler = handler.Export(db)

	suite.Run(t, s)
}

func (suite *ExportSuite) TestSuccess() {
	req, _ := http.NewRequest("GET", "/export/?query=select+*+from+repositories", nil)
	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)

	r := csv.NewReader(res.Body)

	record, err := r.Read()
	suite.Nil(err)
	suite.Equal(record, []string{"id"})

	record, err = r.Read()
	suite.Nil(err)
	suite.Equal(len(record), 1)
	suite.True(len(record[0]) > 0)
}

func (suite *ExportSuite) TestError() {
	req, _ := http.NewRequest("GET", "/export/?query=select+*+from+not_exist", nil)
	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *ExportSuite) TearDownSuite() {
	suite.db.Close()
}
