package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/src-d/gitbase-web/server/serializer"
	"github.com/src-d/gitbase-web/server/service"

	"gopkg.in/bblfsh/client-go.v3"
	"gopkg.in/bblfsh/client-go.v3/tools"
	"gopkg.in/bblfsh/sdk.v2/uast/nodes"
)

type parseRequest struct {
	ServerURL string `json:"serverUrl"`
	Language  string `json:"language"`
	Filename  string `json:"filename"`
	Content   string `json:"content"`
	Filter    string `json:"filter"`
}

// Parse returns a function that parses text contents using bblfsh and
// returns UAST
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

		resp, lang, err := cli.NewParseRequest().
			Language(req.Language).
			Filename(req.Filename).
			Content(req.Content).
			Mode(bblfsh.Semantic).
			UAST()

		if bblfsh.ErrSyntax.Is(err) {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing UAST: %s", err))
		}
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if req.Filter != "" {
			resp, err = applyXpath(resp, req.Filter)
			if err != nil {
				return nil, err
			}
		}

		return serializer.NewParseResponse(&service.ParseResponse{
			UAST: resp,
			Lang: lang,
		}), nil
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

		reqNodes, err := service.UnmarshalNodes(data)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		var resp nodes.Array

		if req.Filter != "" {
			resp, err = applyXpath(reqNodes, req.Filter)
			if err != nil {
				return nil, err
			}
		} else {
			resp = reqNodes
		}

		return serializer.UASTFilterResponse(resp), nil
	}
}

func applyXpath(n nodes.Node, query string) (nodes.Array, serializer.HTTPError) {
	iter, err := tools.Filter(n, query)
	if err != nil {
		return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	results := nodes.Array{}
	for iter.Next() {
		results = append(results, iter.Node().(nodes.Node))
	}

	return results, nil
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

		langs := service.DriverManifestsToLangs(resp)

		sort.Slice(langs, func(i, j int) bool {
			return langs[i].Name < langs[j].Name
		})

		return serializer.NewLanguagesResponse(langs), nil
	}
}
