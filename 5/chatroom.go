package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type chatMessage struct {
	User    string    `json="username"`
	Message string    `json="msg"`
	Time    time.Time `json="submitted_at"`
}

const MESSAGE_FILE = "messages.json"

var _messages_lock = sync.RWMutex{}
var _messages []chatMessage = nil

func getMessages() []chatMessage {
	_messages_lock.RLock()
	defer _messages_lock.RUnlock()

	res := make([]chatMessage, len(_messages))
	for i := 0; i < len(res); i++ {
		mi := len(_messages) - 1 - i
		res[i] = _messages[mi]
	}
	return res
}

func saveMessage(username string, message string) error {
	_messages_lock.Lock()
	defer _messages_lock.Unlock()

	addMessage := chatMessage{
		User:    username,
		Message: message,
		Time:    time.Now(),
	}

	_messages = append(_messages, addMessage)
	asBytes, err := json.Marshal(_messages)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(MESSAGE_FILE, asBytes, 0666)
}

func init() {
	// No need to acquire _messages_lock -- this is called before
	// main can run.
	txt, err := ioutil.ReadFile(MESSAGE_FILE)
	if err != nil {
		_messages = []chatMessage{}
		return
	}

	err = json.Unmarshal(txt, &_messages)
	if err != nil {
		log.Fatal("Could not parse JSON messages file:", err)
	}
}
