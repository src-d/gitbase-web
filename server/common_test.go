package server_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/suite"
)

func GetResponse(method string, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("it should be possible to build a request; %s", err.Error())
	}

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("the server should answer with a response; %s", err.Error())
	}

	return resp, nil
}

type ClientTestSuite struct {
	suite.Suite
}

func (c *ClientTestSuite) AssertResponseBody(resp *http.Response, expectedContent string, msg string) {
	c.Require().NotNil(resp, "the response body should not be nil")
	respBody, err := ioutil.ReadAll(resp.Body)
	c.Require().Nil(err, "the response body should be readable")

	defer resp.Body.Close()
	c.Equal(expectedContent, string(respBody), msg)
}

func (c *ClientTestSuite) AssertResponseStatus(resp *http.Response, expectedStatus int, msg string) {
	c.Require().NotNil(resp, "the response body should not be nil")
	c.Equal(expectedStatus, resp.StatusCode, fmt.Sprintf("status should be %d; %s", expectedStatus, msg))
}

func (c *ClientTestSuite) AssertResponseBodyStatus(resp *http.Response, expectedStatus int, expectedContent string, msg string) {
	c.AssertResponseBody(resp, expectedContent, msg)
	c.AssertResponseStatus(resp, expectedStatus, "")
}
