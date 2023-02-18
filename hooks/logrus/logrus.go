// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logrus contains logrus's hook for Mattermost.
//
// Features:
// - Asynchronous
// - No level filter: all levels from logrus will be send to Mattermost.
// - Sending log as message attachment (see NewHook)
//
// # Example
//
//	import (
//		"github.com/sirupsen/logrus"
//		mmlogrus "github.com/shuLhan/mattermost-integration/hooks/logrus"
//	)
//
//	int main() {
//		_endpoint := https://my.mattermost.org/hooks/xxx"
//		_channel := "log_alpha"
//		_username := "app-name"
//
//		logrus.AddHook(mmlogrus.NewHook(_endpoint, _channel, _username, nil))
//
//		logrus.WithFields(logrus.Fields{
//			"k1": "v1",
//			"k2": "v2",
//		}).Info("Test info")
//
//		mmlogrus.Stop()
//	}
package logrus

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	_httpTr   *http.Transport
	_httpCl   *http.Client
	_chanMsg  chan *Message
	_chanSent chan string
	_running  bool
)

// send will send message `msg` to Mattermost.
//
// On success it will return the HTTP response body with nil error.
// On fail it will return empty response with error message.
func send(msg *Message) (sResBody string, err error) {
	var (
		reqBody, resBody []byte
		body             *bytes.Reader
		req              *http.Request
		res              *http.Response
	)

	reqBody, err = msg.MarshalJSON()
	if err != nil {
		return
	}

	body = bytes.NewReader(reqBody)

	req, err = http.NewRequest("POST", _hook.Endpoint(), body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err = _httpCl.Do(req)
	if err != nil {
		return
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err == nil {
		sResBody = string(resBody)
	}

	err = res.Body.Close()

	return
}

// consumer will consume message from channel `_chanMsg` to be send to
// Mattermost.
func consumer() {
	_running = true
	for _running {
		select {
		case msg, ok := <-_chanMsg:
			if !ok {
				goto out
			}

			go func() {
				res, err := send(msg)
				if err != nil {
					_chanSent <- err.Error()
				} else {
					_chanSent <- res
				}
			}()
		case <-_chanSent:
			continue
		}
	}
out:
	_running = false
}

// Stop will wait for all message to be send and close all channels.
func Stop() {
	_running = false
	if _chanMsg != nil {
		close(_chanMsg)
	}
}

// Start will start the message consumer routine.
func Start() {
	_httpTr = &http.Transport{
		MaxIdleConns:       3,
		IdleConnTimeout:    time.Minute,
		DisableCompression: false,
	}

	_httpCl = &http.Client{
		Transport: _httpTr,
	}

	_chanMsg = make(chan *Message, 30)
	_chanSent = make(chan string, 30)

	go consumer()
}
