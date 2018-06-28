package handler

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	bblfsh "gopkg.in/bblfsh/client-go.v2"
	"gopkg.in/bblfsh/client-go.v2/tools"
	"gopkg.in/bblfsh/sdk.v1/protocol"
	"gopkg.in/bblfsh/sdk.v1/uast"

	"github.com/src-d/gitbase-playground/server/serializer"
	"github.com/src-d/gitbase-playground/server/service"
)

type parseRequest struct {
	ServerURL string `json:"serverUrl"`
	Language  string `json:"language"`
	Filename  string `json:"filename"`
	Content   string `json:"content"`
	Filter    string `json:"filter"`
}

// Parse returns a function that parse content using bblfsh and returns UAST
func Parse(bbblfshServerURL string) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		var req parseRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if req.ServerURL != "" {
			bbblfshServerURL = req.ServerURL
		}

		cli, err := bblfsh.NewClient(bbblfshServerURL)
		if err != nil {
			return nil, err
		}

		resp, err := cli.NewParseRequest().
			Language(req.Language).
			Filename(req.Filename).
			Content(req.Content).
			Do()
		if err != nil {
			return nil, err
		}

		if resp.Status == protocol.Error {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, "incorrect request")
		}

		if resp.Status != protocol.Ok {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, strings.Join(resp.Errors, "\n"))
		}

		if resp.UAST != nil && req.Filter != "" {
			filtered, err := tools.Filter(resp.UAST, req.Filter)
			if err != nil {
				return nil, err
			}

			resp.UAST = &uast.Node{
				InternalType: "Search results",
				Children:     filtered,
			}
		}

		return serializer.NewParseResponse((*service.ParseResponse)(resp)), nil
	}
}

type filterRequest struct {
	Protobufs string `json:"protobufs"`
	Filter    string `json:"filter"`
}

// Filter returns a function that filters UAST protobuf and returns UAST JSON
func Filter() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		var req filterRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		data, err := base64.StdEncoding.DecodeString(req.Protobufs)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		nodes, err := service.UnmarshallUAST(&data)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp := &uast.Node{InternalType: "Search results"}

		if req.Filter != "" {
			for _, n := range nodes {
				filtered, err := tools.Filter((*uast.Node)(n), req.Filter)
				if err != nil {
					return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
				}

				resp.Children = append(resp.Children, filtered...)
			}
		} else {
			for _, n := range nodes {
				resp.Children = append(resp.Children, (*uast.Node)(n))
			}
		}

		return serializer.UASTFilterResponse((*service.Node)(resp)), nil
	}
}

// GetLanguages returns a list of supported languages by bblfsh
func GetLanguages(bbblfshServerURL string) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		cli, err := bblfsh.NewClient(bbblfshServerURL)
		if err != nil {
			return nil, err
		}

		resp, err := cli.NewSupportedLanguagesRequest().Do()
		if err != nil {
			return nil, err
		}

		return serializer.NewLanguagesResponse(service.DriverManifestsToLangs(resp.Languages)), nil
	}
}
