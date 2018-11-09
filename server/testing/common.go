package testing

import (
	"context"
	"database/sql"
)

// MockDB is a mock of *sql.DB
type MockDB struct{}

// Close closes the DB connection
func (db *MockDB) Close() error {
	return nil
}

// Ping pings the DB
func (db *MockDB) Ping() error {
	return nil
}

// Query sends a query to the DB
func (db *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

// QueryContext executes a query
func (db *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

// QueryRow executes a query that is expected to return at most one row
func (db *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

// Exec executes a query without returning any rows
func (db *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

// As returned by gitbase v0.18.0-beta.1, Bblfsh v2.9.2-drivers
// SELECT UAST('console.log("test")', 'JavaScript') AS uast
const (
	UASTMarshaled     = "\x00bgr\x01\x00\x00\x00\x04\bW\x10\x01\x03B\x01\x02\x0e:\x05\x03\x04\x05\x06\aB\x05\b\x16\x17\x18\x19\x06\x12\x04@pos\a\x12\x05@role\a\x12\x05@type\n\x12\bcomments\t\x12\aprogram\n:\x03\x05\t\nB\x03\v\f\x14\x05\x12\x03end\a\x12\x05start\x10\x12\x0euast:Positions\f:\x04\x05\r\x0e\x0fB\x04\x10\x11\x12\x13\x05\x12\x03col\x06\x12\x04line\b\x12\x06offset\x0f\x12\ruast:Position\x02 \x14\x02 \x01\x02 \x13\bB\x04\x10\x12\x12\x15P\f\x02 \x00\x03B\x01\x17\x06\x12\x04File\x00\x10:\x06\x03\x04\x05\x1a\x1b\x1cB\x06\b\x1d\x1f \x18V\x06\x12\x04body\f\x12\ndirectives\f\x12\nsourceType\x03B\x01\x1e\b\x12\x06Module\t\x12\aProgram\x03B\x01!\f:\x04\x03\x04\x05\"B\x04\b#%&\f\x12\nexpression\x03B\x01$\v\x12\tStatement\x15\x12\x13ExpressionStatement\x0e:\x05\x03\x04\x05'(B\x05\b),-<\v\x12\targuments\b\x12\x06callee\x04B\x02*+\f\x12\nExpression\x06\x12\x04Call\x10\x12\x0eCallExpression\x03B\x01.\x0e:\x05\x03\x04\x05/0B\x05179:;\b\x12\x06Format\a\x12\x05Value\aB\x03\v24P\b\bB\x04\x10\x13\x123P\f\x02 \x12\bB\x04\x105\x126P\f\x02 \r\x02 \f\x04B\x02+8\n\x12\bArgument\r\x12\vuast:String\x02\x12\x00\x06\x12\x04test\x10:\x06\x03\x04\x05=>?B\x06@CGHIQ\n\x12\bcomputed\b\x12\x06object\n\x12\bproperty\aB\x03\vA\x14P\b\bB\x04\x106\x12BP\f\x02 \v\aB\x05D*E+F\v\x12\tQualified\f\x12\nIdentifier\b\x12\x06Callee\x12\x12\x10MemberExpression\x020\x00\n:\x03\x03\x05JB\x03KOP\x06\x12\x04Name\aB\x03\vL\x14P\b\bB\x04\x10M\x12NP\f\x02 \b\x02 \a\x11\x12\x0fuast:Identifier\t\x12\aconsole\aB\x03ROUPI\aB\x03\vASP\b\bB\x04\x10T\x12MP\f\x02 \t\x05\x12\x03log\b\x12\x06module"
	UASTMarshaledJSON = `[
  {
    "@pos": {
      "@type": "uast:Positions",
      "end": {
        "@type": "uast:Position",
        "col": 20,
        "line": 1,
        "offset": 19
      },
      "start": {
        "@type": "uast:Position",
        "col": 1,
        "line": 1,
        "offset": 0
      }
    },
    "@role": [
      "File"
    ],
    "@type": "File",
    "comments": [],
    "program": {
      "@pos": {
        "@type": "uast:Positions",
        "end": {
          "@type": "uast:Position",
          "col": 20,
          "line": 1,
          "offset": 19
        },
        "start": {
          "@type": "uast:Position",
          "col": 1,
          "line": 1,
          "offset": 0
        }
      },
      "@role": [
        "Module"
      ],
      "@type": "Program",
      "body": [
        {
          "@pos": {
            "@type": "uast:Positions",
            "end": {
              "@type": "uast:Position",
              "col": 20,
              "line": 1,
              "offset": 19
            },
            "start": {
              "@type": "uast:Position",
              "col": 1,
              "line": 1,
              "offset": 0
            }
          },
          "@role": [
            "Statement"
          ],
          "@type": "ExpressionStatement",
          "expression": {
            "@pos": {
              "@type": "uast:Positions",
              "end": {
                "@type": "uast:Position",
                "col": 20,
                "line": 1,
                "offset": 19
              },
              "start": {
                "@type": "uast:Position",
                "col": 1,
                "line": 1,
                "offset": 0
              }
            },
            "@role": [
              "Expression",
              "Call"
            ],
            "@type": "CallExpression",
            "arguments": [
              {
                "@pos": {
                  "@type": "uast:Positions",
                  "end": {
                    "@type": "uast:Position",
                    "col": 19,
                    "line": 1,
                    "offset": 18
                  },
                  "start": {
                    "@type": "uast:Position",
                    "col": 13,
                    "line": 1,
                    "offset": 12
                  }
                },
                "@role": [
                  "Call",
                  "Argument"
                ],
                "@type": "uast:String",
                "Format": "",
                "Value": "test"
              }
            ],
            "callee": {
              "@pos": {
                "@type": "uast:Positions",
                "end": {
                  "@type": "uast:Position",
                  "col": 12,
                  "line": 1,
                  "offset": 11
                },
                "start": {
                  "@type": "uast:Position",
                  "col": 1,
                  "line": 1,
                  "offset": 0
                }
              },
              "@role": [
                "Qualified",
                "Expression",
                "Identifier",
                "Call",
                "Callee"
              ],
              "@type": "MemberExpression",
              "computed": false,
              "object": {
                "@pos": {
                  "@type": "uast:Positions",
                  "end": {
                    "@type": "uast:Position",
                    "col": 8,
                    "line": 1,
                    "offset": 7
                  },
                  "start": {
                    "@type": "uast:Position",
                    "col": 1,
                    "line": 1,
                    "offset": 0
                  }
                },
                "@type": "uast:Identifier",
                "Name": "console"
              },
              "property": {
                "@pos": {
                  "@type": "uast:Positions",
                  "end": {
                    "@type": "uast:Position",
                    "col": 12,
                    "line": 1,
                    "offset": 11
                  },
                  "start": {
                    "@type": "uast:Position",
                    "col": 9,
                    "line": 1,
                    "offset": 8
                  }
                },
                "@type": "uast:Identifier",
                "Name": "log"
              }
            }
          }
        }
      ],
      "directives": [],
      "sourceType": "module"
    }
  }
]`
)
