//
// Package logrus contains logrus's hook for Mattermost.
//
// Features:
// - Asynchronous
// - No level filter: all levels from logrus will be send to Mattermost.
// - Sending log as message attachment (see NewHook)
//
// Example
//
//```
// import (
//	"github.com/Sirupsen/logrus"
//	mmlogrus "github.com/shuLhan/mattermost-integration/hooks/logrus"
// )
//
// int main() {
//	_endpoint := https://my.mattermost.org/hooks/xxx"
//	_channel := "log_alpha"
//	_username := "app-name"
//
//	logrus.AddHook(mmlogrus.NewHook(_endpoint, _channel, _username, nil))
//
//	logrus.WithFields(logrus.Fields{
//		"k1": "v1",
//		"k2": "v2",
//	}).Info("Test info")
//
//	mmlogrus.Stop()
// }
//```
//
package logrus

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	_httpTr   *http.Transport
	_httpCl   *http.Client
	_chanMsg  chan *Message
	_chanSent chan string
	_chanSig  chan os.Signal
	_wg       sync.WaitGroup
)

func send(msg *Message) {
	var (
		err              error
		reqBody, resBody []byte
		body             *bytes.Reader
		req              *http.Request
		res              *http.Response
	)

	reqBody, err = json.Marshal(msg)
	if err != nil {
		println(">>> json.Marshal:", err)
		goto out
	}

	body = bytes.NewReader(reqBody)

	req, err = http.NewRequest("POST", _hook.Endpoint, body)
	if err != nil {
		println(">>> http.NewRequest:", err)
		goto out
	}

	req.Header.Set("Content-Type", "application/json")

	res, err = _httpCl.Do(req)
	if err != nil {
		println(">>> client.Do:", err)
		goto out
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		println(">>> ioutil.ReadAll:", err)
		goto out
	}

	err = res.Body.Close()
	if err != nil {
		println(">>> res body:", string(resBody))
		println(">>> res.Body.Close:", err)
	}
out:
	go func() {
		if err != nil {
			_chanSent <- err.Error()
		} else {
			_chanSent <- string(resBody)
		}
	}()
}

//
// consumer will consume message from channel `_chanMsg` to be send to
// Mattermost.
//
func consumer() {
	for {
		select {
		case msg, ok := <-_chanMsg:
			if ok {
				send(msg)
			} else {
				goto out
			}
		}
	}
out:
	_wg.Done()
}

//
// Stop will wait for all message to be send and close all channels.
//
func Stop() {
	close(_chanMsg)
	_wg.Wait()
	close(_chanSent)
}

func init() {
	_httpTr = &http.Transport{
		MaxIdleConns:       3,
		IdleConnTimeout:    time.Minute,
		DisableCompression: false,
	}

	_httpCl = &http.Client{
		Transport: _httpTr,
	}

	_chanMsg = make(chan *Message, 1)
	_chanSent = make(chan string, 1)

	_wg.Add(1)
	go consumer()

	_chanSig = make(chan os.Signal, 1)
	signal.Notify(_chanSig, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGSEGV)
	go func() {
		<-_chanSig
		Stop()
		close(_chanSig)
	}()
}
