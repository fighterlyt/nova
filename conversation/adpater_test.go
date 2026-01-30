package conversation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v4"
)

var (
	testTGAdapter Adapter
	err           error
	bot           *telebot.Bot
)

func TestNewTelegramAdapter(t *testing.T) {
	bot, err = telebot.NewBot(telebot.Settings{
		Token:  `7527978229:AAH-g7wf2fr0zE4aZrDXGWrJfysj8NvQ0mc`,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	},
	)

	require.NoError(t, err)

	testTGAdapter = NewTelegramAdapter(bot)
}

func TestTelegramAdapter_Start(t *testing.T) {
	TestNewTelegramAdapter(t)

	for item := range testTGAdapter.GetChan() {
		t.Log(item.Prefix, item.Argument)

		item.Responser.Response(item.Prefix, nil)
	}
}
