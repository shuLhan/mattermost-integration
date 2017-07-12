// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"os"
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/sirupsen/logrus"
)

const (
	envEndpointName = "MM_HOOK_LOGRUS_ENDPOINT"
	envChannelName  = "MM_HOOK_LOGRUS_CHANNEL"
	envUsernameName = "MM_HOOK_LOGRUS_USERNAME"
)

var (
	_endpoint, _channel, _username string
)

func assert(t *testing.T, exp, got interface{}, equal bool) {
	if reflect.DeepEqual(exp, got) != equal {
		debug.PrintStack()
		t.Fatalf("\n"+
			">>> Expecting '%+v'\n"+
			"          got '%+v'\n", exp, got)
		os.Exit(1)
	}
}

func TestFire(t *testing.T) {
	NewHook(_endpoint, _channel, _username, nil)

	tests := []struct {
		desc string
		in   logrus.Entry
		exp  string
	}{
		{
			desc: "With empty fields",
			in: logrus.Entry{
				Level:   logrus.DebugLevel,
				Message: "Test with empty field",
			},
			exp: "ok",
		},
		{
			desc: "With message",
			in: logrus.Entry{
				Message: "Test",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level debug",
			in: logrus.Entry{
				Level:   logrus.DebugLevel,
				Message: "Test debug",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level info",
			in: logrus.Entry{
				Level:   logrus.InfoLevel,
				Message: "Test info",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level warning",
			in: logrus.Entry{
				Level:   logrus.WarnLevel,
				Message: "Test warning",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level error",
			in: logrus.Entry{
				Level:   logrus.ErrorLevel,
				Message: "Test error",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level fatal",
			in: logrus.Entry{
				Level:   logrus.FatalLevel,
				Message: "Test fatal",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level panic",
			in: logrus.Entry{
				Level:   logrus.PanicLevel,
				Message: "Test panic",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		err := _hook.Fire(&test.in)
		if err != nil {
			t.Fatal(err)
		}

		res := <-_chanSent

		assert(t, test.exp, res, true)
	}
}

func TestFireWithAttachment(t *testing.T) {
	attc := Attachment{
		Pretext: "Send from test",
	}

	NewHook(_endpoint, _channel, _username, &attc)

	tests := []struct {
		desc string
		in   logrus.Entry
		exp  string
	}{
		{
			desc: "With empty fields",
			in: logrus.Entry{
				Level:   logrus.DebugLevel,
				Message: "Test attachment with empty field",
			},
			exp: "ok",
		},
		{
			desc: "With message",
			in: logrus.Entry{
				Message: "Test with attachment",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level debug",
			in: logrus.Entry{
				Level:   logrus.DebugLevel,
				Message: "Test attachment debug",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level info",
			in: logrus.Entry{
				Level:   logrus.InfoLevel,
				Message: "Test attachment info",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level warning",
			in: logrus.Entry{
				Level:   logrus.WarnLevel,
				Message: "Test attachment warning",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level error",
			in: logrus.Entry{
				Level:   logrus.ErrorLevel,
				Message: "Test attachment error",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level fatal",
			in: logrus.Entry{
				Level:   logrus.FatalLevel,
				Message: "Test attachment fatal",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
		{
			desc: "With level panic",
			in: logrus.Entry{
				Level:   logrus.PanicLevel,
				Message: "Test attachment panic",
				Data: logrus.Fields{
					"k1": "v1",
					"k2": "v2",
				},
			},
			exp: "ok",
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		err := _hook.Fire(&test.in)
		if err != nil {
			t.Fatal(err)
		}
	}
}
