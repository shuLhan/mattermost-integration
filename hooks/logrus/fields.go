// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"bytes"
)

//
// Fields define list of field in message attachments.
//
type Fields []Field

//
// MarshalJSON will convert `field` into a valid JSON. We use manual
// convertion for gaining speed.
//
func (fields Fields) MarshalJSON() (out []byte, err error) {
	var buf bytes.Buffer
	sep := false

	_ = buf.WriteByte('[')

	for _, field := range fields {
		if sep {
			_ = buf.WriteByte(',')
		}

		fout, _ := field.MarshalJSON()

		if len(fout) > 2 {
			_, _ = buf.Write(fout)
			sep = true
		} else {
			sep = false
		}
	}

	_ = buf.WriteByte(']')
	out = buf.Bytes()

	return
}
