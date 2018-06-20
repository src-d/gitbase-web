package service

import (
	"encoding/json"
	"fmt"

	"gopkg.in/bblfsh/sdk.v1/protocol"
	"gopkg.in/bblfsh/sdk.v1/uast"
)

// need to move to service to avoid circular imports

// UnmarshallUAST tries to cast data as [][]byte and unmarshall uast nodes
func UnmarshallUAST(data interface{}) ([]*Node, error) {
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

// ParseResponse amends default MarshalJSON to be compatible with frontend
type ParseResponse protocol.ParseResponse

// MarshalJSON returns the JSON representation of the protocol.ParseResponse
func (r *ParseResponse) MarshalJSON() ([]byte, error) {
	fmt.Println("MarshalJSON")
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
