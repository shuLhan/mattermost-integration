// Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	//
	// _hook contain a singleton instance of `mmHookLogrus`. The
	// subsequence call of `NewHook` will only replace the field value.
	//
	_hook       *mmHookLogrus
	_hookLocker sync.Mutex

	//
	// _iconsLevel contains list of icon to be displayed before log
	// message based on log level.
	//
	_iconsLevel = []string{
		":x:",            // Panic
		":bangbang:",     // Fatal
		":exclamation:",  // Error
		":interrobang:",  // Warning
		":white_circle:", // Info
		":black_circle:", // Debug
		":mag_right:",    // Trace
	}
)

//
// mmHookLogrus contains configuration for Mattermost (server address,
// channel, username) and reusable http transport and client.
//
type mmHookLogrus struct {
	endpoint string
	channel  string
	username string
	defAttc  *Attachment
	hostname string
	levels   []logrus.Level
}

//
// NewHook will create a log hook for mattermost. The log will be send to
// server at `endpoint` inside `channel` name, using `username`.
//
// `channel` and `username` is optional.
// If channel is empty then it will use the default channel defined in
// incoming webhook setting.
// If username is empty then it will use the hostname.
//
// If attachment parameter is not nil, each log will be send as attachment
// [1].  The parameter will act as default attachment value, and it will
// replace the `Text` with `Entry.Message` and `Fields` with `Entry.Data`.
//
// [1] https://docs.mattermost.com/developer/message-attachments.html
//
func NewHook(endpoint, channel, username string, attc *Attachment, minLevel logrus.Level) logrus.Hook {
	levels := make([]logrus.Level, 0, len(logrus.AllLevels))
	for _ lvl := range logrus.AllLevels {
		if lvl >= minLevel {
			levels = append(levels, lvl)	
		}
	}
	
	var err error

	if _hook == nil {
		_hook = &mmHookLogrus{
			endpoint: endpoint,
			channel:  channel,
			username: username,
			defAttc:  attc,
			levels: levels,
		}

		_hook.hostname, err = os.Hostname()
		if err != nil {
			_hook.hostname = os.Getenv("HOSTNAME")
		}
	} else {
		_hookLocker.Lock()

		_hook.endpoint = endpoint
		_hook.channel = channel
		_hook.username = username
		_hook.defAttc = attc

		_hookLocker.Unlock()
	}

	if !_running {
		Start()
	}

	return _hook
}

//
// Levels will return all logrus level that will be send to Mattermost.
//
func (hook *mmHookLogrus) Levels() []logrus.Level {
	return hook.levels
}

//
// Fire will send logrus `entry` to Mattermost.
//
func (hook *mmHookLogrus) Fire(entry *logrus.Entry) (err error) {
	if entry == nil {
		return
	}

	msg := NewMessage(hook.Channel(), hook.Username(), hook.Hostname(),
		hook.Attachment(), entry)

	_chanMsg <- msg

	return
}

//
// Endpoint will return Mattermost endpoint defined in hook.
//
func (hook *mmHookLogrus) Endpoint() (endpoint string) {
	_hookLocker.Lock()
	endpoint = hook.endpoint
	_hookLocker.Unlock()
	return
}

//
// Channel will return Mattermost channel defined in hook.
//
func (hook *mmHookLogrus) Channel() (channel string) {
	_hookLocker.Lock()
	channel = hook.channel
	_hookLocker.Unlock()
	return
}

//
// Username will return Mattermost username defined in hook.
//
func (hook *mmHookLogrus) Username() (username string) {
	_hookLocker.Lock()
	username = hook.username
	_hookLocker.Unlock()
	return
}

//
// Attachment will return default Mattermost attachment defined in hook.
//
func (hook *mmHookLogrus) Attachment() (attc *Attachment) {
	_hookLocker.Lock()
	attc = hook.defAttc
	_hookLocker.Unlock()
	return
}

//
// Hostname will return hostname of current hook.
//
func (hook *mmHookLogrus) Hostname() (hostname string) {
	_hookLocker.Lock()
	hostname = hook.hostname
	_hookLocker.Unlock()
	return
}
