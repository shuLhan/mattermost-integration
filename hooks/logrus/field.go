// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"bytes"
)

// Field define a single field in message attachment.
type Field struct {
	Short bool
	Title string
	Value string
}

// MarshalJSON will convert `field` into a valid JSON.
//
// (1) The conversion will skip empty field Title or Value.
//
// Returned error always nil.
func (field Field) MarshalJSON() (out []byte, err error) {
	var buf bytes.Buffer

	_ = buf.WriteByte('{')

	// (1)
	if len(field.Title) == 0 || len(field.Value) == 0 {
		goto out
	}

	if field.Short {
		_ = bufWriteKV(&buf, `"short"`, []byte("true"), ':', 0, 0)
	} else {
		_ = bufWriteKV(&buf, `"short"`, []byte("false"), ':', 0, 0)
	}

	_ = buf.WriteByte(',')
	_ = bufWriteKV(&buf, `"title"`, []byte(field.Title), ':', '"', '"')
	_ = buf.WriteByte(',')
	_ = bufWriteKV(&buf, `"value"`, []byte(field.Value), ':', '"', '"')

out:
	_ = buf.WriteByte('}')
	out = buf.Bytes()

	return
}
