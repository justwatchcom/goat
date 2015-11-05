package xsd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testCase struct {
	comment, name  string
	path           []string
	params         map[string]interface{}
	isFaulty       bool
	expectedResult string
}

var (
	testCases = []testCase{
		{
			comment: "basic string encoding",
			name:    "string",
			path:    []string{"simple_path"},
			params: map[string]interface{}{
				"simple_path": "test",
			},
			expectedResult: "test",
		},
		{
			comment:  "basic string encoding with missing value",
			name:     "string",
			path:     []string{"simple_path"},
			params:   map[string]interface{}{},
			isFaulty: true,
		},
		{
			comment: "basic string encoding with complex path",
			name:    "string",
			path:    []string{"path", "to", "glory"},
			params: map[string]interface{}{
				"path/to/glory": "test",
			},
			expectedResult: "test",
		},
		{
			comment: "basic string encoding with int parameter",
			name:    "string",
			path:    []string{"simple_path"},
			params: map[string]interface{}{
				"simple_path": 1,
			},
			isFaulty: true,
		},
		{
			comment: "basic int encoding",
			name:    "int",
			path:    []string{"simple_path"},
			params: map[string]interface{}{
				"simple_path": 1,
			},
			expectedResult: "1",
		},
	}
)

func TestBaseSchema_EncodeType(t *testing.T) {
	Convey("given a baseSchema instance and a SchemaRepository", t, func() {
		bs := baseSchema{}
		sr := SchemaMap{}
		for _, c := range testCases {
			Convey(fmt.Sprintf("test case '%s'", c.comment), func() {
				buf := new(bytes.Buffer)
				enc := xml.NewEncoder(buf)
				err := bs.EncodeType(c.name, enc, sr, c.params, c.path...)
				if c.isFaulty {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					err = enc.Flush()
					So(err, ShouldBeNil)
					So(buf.String(), ShouldEqual, c.expectedResult)
				}
			})
		}
	})
}
