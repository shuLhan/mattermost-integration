// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"testing"
)

func TestFieldsUnmarshalJSON(t *testing.T) {
	tests := []struct {
		desc string
		in   Fields
		exp  string
	}{
		{
			desc: "With empty fields",
			in:   Fields{},
			exp:  "[]",
		},
		{
			desc: "With one field",
			in: Fields{
				{
					Title: "t1",
					Value: "v1",
				},
			},
			exp: `[{"short":false,"title":"t1","value":"v1"}]`,
		},
		{
			desc: "With empty field in middle",
			in: Fields{
				{
					Title: "t1",
					Value: "v1",
				},
				{
					Short: true,
				},
				{
					Title: "t3",
					Value: "v3",
				},
			},
			exp: `[{"short":false,"title":"t1","value":"v1"},{"short":false,"title":"t3","value":"v3"}]`,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		got, err := test.in.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		assert(t, test.exp, string(got), true)
	}
}
