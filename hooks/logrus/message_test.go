package logrus

import (
	"testing"

	"github.com/sirupsen/logrus"
)

var expMsgJSON = []byte(`{"channel":"test","username":"testing","text":":white_circle: json={\"msg\":\"field with \\\"JSON\\\"\"} number=1 string=a string struct={n:10 s:string in struct} msg={\"message\":\"this is \\\"JSON\\\"\"}"}`)

func newMessage() *Message {
	return &Message{
		channel:    "test",
		username:   "testing",
		hostname:   "localhost",
		attc:       nil,
		entryLevel: logrus.InfoLevel,
		entryData: logrus.Fields{
			"number": 1,
			"string": "a string",
			"struct": struct {
				n int
				s string
			}{
				n: 10,
				s: "string in struct",
			},
			"json": `{"msg":"field with \"JSON\""}`,
		},
		entryMsg: `{"message":"this is \"JSON\""}`,
	}
}

func BenchmarkMarshalJSONOld(b *testing.B) {
	msg := newMessage()
	got, err := msg._marshalJSON()
	if err != nil {
		b.Fatal(err)
	}

	b.Log(string(got))

	assertb(b, expMsgJSON, got, true)

	for x := 0; x < b.N; x++ {
		_, _ = msg._marshalJSON()
	}
}

func BenchmarkMarshalJSONBuffer(b *testing.B) {
	msg := newMessage()
	got, err := msg.MarshalJSON()
	if err != nil {
		b.Fatal(err)
	}

	b.Log(string(got))

	assertb(b, expMsgJSON, got, true)

	for x := 0; x < b.N; x++ {
		_, _ = msg.MarshalJSON()
	}
}
