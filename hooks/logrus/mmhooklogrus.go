package logrus

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	//
	// _hook contain a singleton instance of `mmHookLogrus`. The
	// subsequence call of `NewHook` will only replace the field value.
	//
	_hook *mmHookLogrus

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
	}
)

//
// mmHookLogrus contains configuration for Mattermost (server address,
// channel, username) and reusable http transport and client.
//
type mmHookLogrus struct {
	Endpoint string
	Channel  string
	Username string
	DefAttc  *Attachment
	hostname string
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
func NewHook(endpoint, channel, username string, attc *Attachment) logrus.Hook {
	var err error

	if _hook == nil {
		_hook = &mmHookLogrus{
			Endpoint: endpoint,
			Channel:  channel,
			Username: username,
			DefAttc:  attc,
		}

		_hook.hostname, err = os.Hostname()
		if err != nil {
			_hook.hostname = os.Getenv("HOSTNAME")
		}
	} else {
		_hook.Endpoint = endpoint
		_hook.Channel = channel
		_hook.Username = username
		_hook.DefAttc = attc
	}

	return _hook
}

//
// Levels will return all logrus level that will be send to Mattermost.
//
func (hook *mmHookLogrus) Levels() []logrus.Level {
	return logrus.AllLevels
}

//
// Fire will send logrus `entry` to Mattermost.
//
func (hook *mmHookLogrus) Fire(entry *logrus.Entry) (err error) {
	if entry == nil {
		return
	}

	msg := NewMessage(hook.DefAttc, entry)

	_chanMsg <- msg

	return
}
