package handler

import (
	"net/http"

	"github.com/src-d/gitbase-playground/server/service"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type HandlerUnitSuite struct {
	suite.Suite
	db                 service.SQLDB
	mock               sqlmock.Sqlmock
	handler            http.Handler
	logger             *logrus.Logger
	requestProcessFunc func(db service.SQLDB) RequestProcessFunc
	IsIntegration      bool
}

func (suite *HandlerUnitSuite) SetupSuite() {
	// logger
	suite.logger = logrus.New()
	suite.logger.SetLevel(logrus.FatalLevel)
}

func (suite *HandlerUnitSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to initialize the mock DB. '%s'", err)
	}

	h := APIHandlerFunc(suite.requestProcessFunc(suite.db))
	suite.handler = lg.RequestLogger(suite.logger)(h)
}

func (suite *HandlerUnitSuite) TearDownTest() {
	suite.db.Close()
}
