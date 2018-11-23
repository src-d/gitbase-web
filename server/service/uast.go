package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/bblfsh/sdk.v1/protocol"
	"gopkg.in/bblfsh/sdk.v1/uast"
	errors "gopkg.in/src-d/go-errors.v1"
)

// need to move to service to avoid circular imports

// UnmarshalUASTOld tries to cast data as [][]byte and unmarshall uast node.
// This is the format returned by gitbase <= v0.16.0
func UnmarshalUASTOld(data interface{}) ([]*Node, error) {
	var protobufs [][]byte
	if err := json.Unmarshal(*data.(*[]byte), &protobufs); err != nil {
		return nil, err
	}

	nodes := make([]*Node, len(protobufs))

	for i, v := range protobufs {
		n := uast.NewNode()
		if err := n.Unmarshal(v); err != nil {
			return nil, err
		}
		nodes[i] = (*Node)(n)
	}

	return nodes, nil
}

// TODO (carlosms): Duplicated code from gitbase,
// (internal/function/uast_utils.go) we should reuse that instead

var (
	// ErrParseBlob is returned when the blob can't be parsed with bblfsh.
	ErrParseBlob = errors.NewKind("unable to parse the given blob using bblfsh: %s")

	// ErrUnmarshalUAST is returned when an error arises unmarshaling UASTs.
	ErrUnmarshalUAST = errors.NewKind("error unmarshaling UAST: %s")

	// ErrMarshalUAST is returned when an error arises marshaling UASTs.
	ErrMarshalUAST = errors.NewKind("error marshaling uast node: %s")
)

// UnmarshalUAST returns UAST nodes from data marshaled by gitbase
func UnmarshalUAST(data []byte) ([]*Node, error) {
	if len(data) == 0 {
		return nil, nil
	}

	nodes := []*Node{}
	buf := bytes.NewBuffer(data)
	for {
		var nodeLen int32
		if err := binary.Read(
			buf, binary.BigEndian, &nodeLen,
		); err != nil {
			if err == io.EOF {
				break
			}

			return nil, ErrUnmarshalUAST.New(err)
		}

		if nodeLen < 1 {
			return nil, ErrUnmarshalUAST.New(fmt.Errorf("malformed data"))
		}

		node := uast.NewNode()
		nodeBytes := buf.Next(int(nodeLen))
		if int32(len(nodeBytes)) != nodeLen {
			return nil, ErrUnmarshalUAST.New(fmt.Errorf("malformed data"))
		}

		if err := node.Unmarshal(nodeBytes); err != nil {
			return nil, ErrUnmarshalUAST.New(err)
		}

		nodes = append(nodes, (*Node)(node))
	}

	return nodes, nil
}

// ParseResponse amends default MarshalJSON to be compatible with frontend
type ParseResponse protocol.ParseResponse

// MarshalJSON returns the JSON representation of the protocol.ParseResponse
func (r *ParseResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		*protocol.ParseResponse
		UAST *Node `json:"uast"`
	}{
		(*protocol.ParseResponse)(r),
		(*Node)(r.UAST),
	}

	return json.Marshal(resp)
}

type Node uast.Node

// MarshalJSON returns the JSON representation of the Node
func (n *Node) MarshalJSON() ([]byte, error) {
	var nodes = make([]*Node, len(n.Children))
	for i, n := range n.Children {
		nodes[i] = (*Node)(n)
	}

	var roles = make([]string, len(n.Roles))
	for i, r := range n.Roles {
		roles[i] = r.String()
	}

	node := struct {
		*uast.Node
		Roles    []string
		Children []*Node
	}{
		(*uast.Node)(n),
		roles,
		nodes,
	}

	return json.Marshal(node)
}

type Language struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func DriverManifestsToLangs(drivers []protocol.DriverManifest) []Language {
	result := make([]Language, len(drivers))

	for i, driver := range drivers {
		result[i] = Language{
			ID:   driver.Language,
			Name: driver.Name,
		}
	}

	return result
}
