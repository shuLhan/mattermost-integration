// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"encoding/json"
	"testing"
)

func TestFieldUnmarshalJSON(t *testing.T) {
	tests := []struct {
		desc string
		in   Field
		exp  string
	}{
		{
			desc: "With empty field",
			in:   Field{},
			exp:  "{}",
		},
		{
			desc: "With empty title",
			in: Field{
				Short: false,
				Title: "",
				Value: "value",
			},
			exp: "{}",
		},
		{
			desc: "With empty value",
			in: Field{
				Short: false,
				Title: "title",
				Value: "",
			},
			exp: "{}",
		},
		{
			desc: "With empty title and value",
			in: Field{
				Short: false,
				Title: "",
				Value: "",
			},
			exp: "{}",
		},
		{
			desc: "With all field set",
			in: Field{
				Short: true,
				Title: "title",
				Value: "value",
			},
			exp: `{"short":true,"title":"title","value":"value"}`,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		got, err := json.Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}

		assert(t, test.exp, string(got), true)
	}
}
