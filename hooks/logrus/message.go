package logrus

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
)

//
// Message define the message that will be send to Mattermost.
//
type Message struct {
	attc       *Attachment
	entryData  logrus.Fields
	entryLevel logrus.Level
	entryMsg   string
}

//
// NewMessage will create and return new Message.
//
func NewMessage(attc *Attachment, entry *logrus.Entry) (msg *Message) {
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
// `<username|hostname>: :icon: <field-key=field-value ...> msg=Message`
//
func (msg Message) getText() (str string) {
	str = _iconsLevel[msg.entryLevel]

	for k, v := range msg.entryData {
		str += fmt.Sprintf(" %s=%v", k, v)
	}

	if len(msg.entryMsg) > 0 {
		str += " msg=" + msg.entryMsg
	}

	return
}

//
// MarshalJSON will convert message to JSON.
//
func (msg Message) MarshalJSON() (out []byte, err error) {
	str := `{`

	if len(_hook.Channel) > 0 {
		str += `"channel":"` + _hook.Channel + `",`
	}
	if len(_hook.Username) > 0 {
		str += `"username":"` + _hook.Username + `",`
	} else {
		str += `"username":"` + _hook.hostname + `",`
	}

	if msg.attc != nil {
		var attc []byte

		attc, err = json.Marshal(msg.attc)
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
