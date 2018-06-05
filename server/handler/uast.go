package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	bblfsh "gopkg.in/bblfsh/client-go.v2"
	"gopkg.in/bblfsh/client-go.v2/tools"
	"gopkg.in/bblfsh/sdk.v1/protocol"
	"gopkg.in/bblfsh/sdk.v1/uast"

	"github.com/src-d/gitbase-playground/server/serializer"
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
			return nil, errors.New("bblfsh returend not OK response")
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

		return serializer.NewParseResponse(resp), nil
	}
}

// Filter : placeholder method
func Filter() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}
}

// unmarshallUAST tries to cast data as [][]byte and unmarshall uast nodes
func unmarshallUAST(data interface{}) ([]*node, error) {
	var protobufs [][]byte
	if err := json.Unmarshal(*data.(*[]byte), &protobufs); err != nil {
		return nil, err
	}

	nodes := make([]*node, len(protobufs))

	for i, v := range protobufs {
		n := uast.NewNode()
		if err := n.Unmarshal(v); err != nil {
			return nil, err
		}
		nodes[i] = (*node)(n)
	}

	return nodes, nil
}

type node uast.Node

// MarshalJSON returns the JSON representation of the Node
func (n *node) MarshalJSON() ([]byte, error) {
	var nodes = make([]*node, len(n.Children))
	for i, n := range n.Children {
		nodes[i] = (*node)(n)
	}

	var roles = make([]string, len(n.Roles))
	for i, r := range n.Roles {
		roles[i] = r.String()
	}

	node := struct {
		*uast.Node
		Roles    []string
		Children []*node
	}{
		(*uast.Node)(n),
		roles,
		nodes,
	}

	return json.Marshal(node)
}
