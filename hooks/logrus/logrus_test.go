package logrus

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogrusAddHook(t *testing.T) {
	logrus.AddHook(NewHook(_endpoint, _channel, _username, nil))

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 1,
		"size":   10,
	}).Debug("A walrus debug")

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 2,
		"size":   10,
	}).Info("A walrus info")

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 3,
		"size":   10,
	}).Warn("A walrus warn")

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 3,
		"size":   10,
	}).Error("A walrus error")
}

//
// TestMain will only run if user set the MM_HOOK_LOGRUS_ENDPOINT value in
// environment.
//
func TestMain(m *testing.M) {
	_endpoint = os.Getenv(envEndpointName)
	_channel = os.Getenv(envChannelName)
	_username = os.Getenv(envUsernameName)

	println(">>> Mattermost endpoint: ", _endpoint)
	println(">>> Mattermost channel : ", _channel)
	println(">>> Mattermost username: ", _username)
	println("")

	if len(_endpoint) == 0 {
		println(">>> Environment variable " + envEndpointName + " is empty")
		println(">>> No test will running.")
		os.Exit(0)
	}

	s := m.Run()

	Stop()

	os.Exit(s)
}
