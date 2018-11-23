package service_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/src-d/gitbase-web/server/service"
	"github.com/stretchr/testify/suite"
)

type UastSuite struct {
	suite.Suite
}

func TestUastSuite(t *testing.T) {
	s := new(UastSuite)
	suite.Run(t, s)
}

func (suite *UastSuite) TestNegativeNodeLen() {
	var nodeLen int32 = -20

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, nodeLen)
	suite.Require().NoError(err)

	nodes, err := service.UnmarshalUAST(buf.Bytes())
	suite.Require().Error(err)
	suite.Require().Nil(nodes)
}
