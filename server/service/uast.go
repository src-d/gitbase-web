package service

import (
	"bytes"
	"fmt"

	bblfsh "github.com/bblfsh/go-client"
	"gopkg.in/bblfsh/sdk.v2/uast/nodes"
	"gopkg.in/bblfsh/sdk.v2/uast/nodes/nodesproto"
	errors "gopkg.in/src-d/go-errors.v1"
)

// need to move to service to avoid circular imports

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

// UnmarshalNodes returns UAST nodes from data marshaled by gitbase
func UnmarshalNodes(data []byte) (nodes.Array, error) {
	if len(data) == 0 {
		return nil, nil
	}

	buf := bytes.NewReader(data)
	n, err := nodesproto.ReadTree(buf)
	if err != nil {
		return nil, err
	}
	if n.Kind() != nodes.KindArray {
		return nil, fmt.Errorf("unmarshal: wrong kind of node found %q, expected %q",
			n.Kind(), nodes.KindArray.String())
	}

	return n.(nodes.Array), nil
}

type ParseResponse struct {
	UAST nodes.Node `json:"uast"`
	Lang string     `json:"language"`
}

type Language struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func DriverManifestsToLangs(drivers []bblfsh.DriverManifest) []Language {
	result := make([]Language, len(drivers))

	for i, driver := range drivers {
		result[i] = Language{
			ID:   driver.Language,
			Name: driver.Name,
		}
	}

	return result
}
