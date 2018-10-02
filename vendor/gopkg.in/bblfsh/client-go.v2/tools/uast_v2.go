package tools

import (
	"gopkg.in/bblfsh/sdk.v2/uast/nodes"
	"gopkg.in/bblfsh/sdk.v2/uast/query"
	"gopkg.in/bblfsh/sdk.v2/uast/query/xpath"
)

// FilterXPath takes a `Node` as returned by UAST() call and an xpath query and filters the tree,
// returning the iterator of nodes that satisfy the given query.
func (c *Context) FilterXPath(node nodes.External, query string) (query.Iterator, error) {
	return FilterXPath(node, query)
}

// FilterXPath is a shorthand for Context.FilterXPath.
func FilterXPath(node nodes.External, query string) (query.Iterator, error) {
	return xpath.New().Execute(node, query)
}
