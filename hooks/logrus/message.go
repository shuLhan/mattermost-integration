// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

//
// Message define the message that will be send to Mattermost.
//
type Message struct {
	channel    string
	username   string
	hostname   string
	attc       *Attachment
	entryData  logrus.Fields
	entryLevel logrus.Level
	entryMsg   string
}

//
// NewMessage will create and return new Message.
//
func NewMessage(channel, username, hostname string, attc *Attachment,
	entry *logrus.Entry,
) (msg *Message) {
	msg = &Message{
		attc:       NewAttachment(attc, entry),
		entryData:  entry.Data,
		entryLevel: entry.Level,
		entryMsg:   entry.Message,
	}

	return
}

//
// getText will convert Message into text. The text output format,
//
// `:icon: <field-key=field-value ...> msg=Message`
//
func (msg Message) getText() (str string) {
	var out []byte

	out = append(out, []byte(_iconsLevel[msg.entryLevel])...)

	for k, v := range msg.entryData {
		out = append(out, ' ')
		out = append(out, []byte(k)...)
		out = append(out, '=')

		str = fmt.Sprintf("%+v", v)
		for _, c := range []byte(str) {
			if c == '\\' {
				out = append(out, []byte(`\`)...)
				out = append(out, []byte(`\`)...)
				continue
			}
			if c == '"' {
				out = append(out, []byte(`\`)...)
				out = append(out, []byte(`"`)...)
				continue
			}
			out = append(out, c)
		}
	}

	if len(msg.entryMsg) > 0 {
		out = append(out, ' ')
		out = append(out, []byte("msg=")...)

		for _, c := range []byte(msg.entryMsg) {
			if c == '\\' {
				out = append(out, []byte(`\`)...)
				out = append(out, []byte(`\`)...)
				continue
			}
			if c == '"' {
				out = append(out, []byte(`\`)...)
				out = append(out, []byte(`"`)...)
				continue
			}
			out = append(out, c)
		}
	}

	return string(out)
}

//
// MarshalJSON will convert message to JSON.
//
func (msg *Message) MarshalJSON() (out []byte, err error) {
	str := `{`

	if len(msg.channel) > 0 {
		str += `"channel":"` + msg.channel + `",`
	}

	if len(msg.username) > 0 {
		str += `"username":"` + msg.username + `",`
	} else {
		str += `"username":"` + msg.hostname + `",`
	}

	if msg.attc != nil {
		var attc []byte

		attc, err = msg.attc.MarshalJSON()
		if err != nil {
			return
		}

		str += `"attachments":[` + string(attc) + `]`
	} else {
		str += `"text":"` + msg.getText() + `"`
	}

	str += `}`
	out = []byte(str)

	return
}
