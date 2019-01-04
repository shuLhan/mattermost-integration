// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
)

var (
	_colorsLevel = []string{
		"#FF0000", // Panic
		"#CC0000", // Fatal
		"#990000", // Error
		"#9F6000", // Warning
		"#FFFFFF", // Info
		"#000000", // Debug
		"#000000", // Trace
	}
)

//
// Attachment define Mattermost message attachment [1].
//
// [1] https://docs.mattermost.com/developer/message-attachments.html
//
type Attachment struct {
	AuthorIcon string
	AuthorLink string
	AuthorName string
	Color      string
	Fallback   string
	Fields     Fields
	ImageURL   string
	Pretext    string
	Text       string
	Title      string
	TitleLink  string
}

//
// NewAttachment will create and return new Attachment with default value set
// from `attc` and Color and Fields based on logrus Entry Level and Data.
//
func NewAttachment(defAttc *Attachment, entry *logrus.Entry) (
	attc *Attachment,
) {
	if defAttc == nil {
		return
	}

	attc = &Attachment{
		AuthorIcon: defAttc.AuthorIcon,
		AuthorLink: defAttc.AuthorLink,
		AuthorName: defAttc.AuthorName,
		ImageURL:   defAttc.ImageURL,
		Pretext:    defAttc.Pretext,
		Title:      defAttc.Title,
		TitleLink:  defAttc.TitleLink,
	}

	if entry != nil {
		attc.Color = _colorsLevel[entry.Level]
		attc.SetFields(entry.Data)
		attc.Text = _iconsLevel[entry.Level] + " " + entry.Message
	}

	return
}

func (attc Attachment) marshalAuthor(buf *bytes.Buffer) {
	if len(attc.AuthorIcon) > 0 {
		_ = bufWriteKV(buf, `"author_icon"`, []byte(attc.AuthorIcon),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}
	if len(attc.AuthorLink) > 0 {
		_ = bufWriteKV(buf, `"author_link"`, []byte(attc.AuthorLink),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}
	if len(attc.AuthorName) > 0 {
		_ = bufWriteKV(buf, `"author_name"`, []byte(attc.AuthorName),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}
}

//
// MarshalJSON will convert Attachment `attc` to JSON.
//
func (attc Attachment) MarshalJSON() (out []byte, err error) {
	var buf bytes.Buffer
	var bFields []byte

	_ = buf.WriteByte('{')

	attc.marshalAuthor(&buf)

	if len(attc.Color) > 0 {
		_ = bufWriteKV(&buf, `"color"`, []byte(attc.Color),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}
	if len(attc.Fallback) > 0 {
		_ = bufWriteKV(&buf, `"fallback"`, []byte(attc.Fallback),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}
	if len(attc.Fields) > 0 {
		bFields, err = attc.Fields.MarshalJSON()
		if err != nil {
			return
		}

		_, _ = buf.WriteString(`"fields":`)
		_, _ = buf.Write(bFields)
		_ = buf.WriteByte(',')
	}
	if len(attc.ImageURL) > 0 {
		_ = bufWriteKV(&buf, `"image_url"`, []byte(attc.ImageURL),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}
	if len(attc.Pretext) > 0 {
		_ = bufWriteKV(&buf, `"pretext"`, []byte(attc.Pretext),
			':', '"', '"')
		_ = buf.WriteByte(',')

	}
	if len(attc.Text) > 0 {
		_ = bufWriteKV(&buf, `"text"`, []byte(attc.Text),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}
	if len(attc.Title) > 0 {
		_ = bufWriteKV(&buf, `"title"`, []byte(attc.Title),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}
	if len(attc.TitleLink) > 0 {
		_ = bufWriteKV(&buf, `"title_link"`, []byte(attc.TitleLink),
			':', '"', '"')
		_ = buf.WriteByte(',')
	}

	out = buf.Bytes()

	out = bytes.TrimSuffix(out, []byte(","))
	out = append(out, []byte("}")...)

	return
}

//
// SetFields will convert logrus Fields data `in` into our Fields.
//
func (attc *Attachment) SetFields(in logrus.Fields) {
	attc.Fields = make(Fields, 0)

	if len(in) == 0 {
		return
	}

	var keys []string
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		attc.Fields = append(attc.Fields, Field{
			Short: true,
			Title: k,
			Value: fmt.Sprintf("%+v", in[k]),
		})
	}
}
