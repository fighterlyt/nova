package conversation

import (
	"context"
	"strings"

	"github.com/fighterlyt/nova/databases/pools"
	"gopkg.in/telebot.v4"
)

// Adapter 适配器，用来和外部IO
type Adapter interface {
	GetChan() <-chan Item
}

// Item 外部消息
type Item struct {
	Context   context.Context // 上下文
	Prefix    string          // 前缀
	Argument  string          // 参数
	Responser Responser       // 应答
}

// Response 应答
type Response struct {
	Msg string
	Err error
}

type Responser interface {
	Response(msg string, err error)
}

type channelBaseResponser struct {
	ch   chan Response
	pool *pools.Pool[chan Response]
}

func newChannelBaseResponser(pool *pools.Pool[chan Response], ch chan Response) *channelBaseResponser {
	return &channelBaseResponser{pool: pool, ch: ch}
}

func (r *channelBaseResponser) Response(msg string, err error) {
	r.ch <- Response{Msg: msg, Err: err}
	r.pool.Put(r.ch)
}

type TelegramAdapter struct {
	bot  *telebot.Bot
	ch   chan Item
	pool *pools.Pool[chan Response]
}

func NewTelegramAdapter(bot *telebot.Bot) *TelegramAdapter {
	var (
		ch   = make(chan Item, 1)
		pool = pools.NewPool[chan Response](10, func() chan Response {
			return make(chan Response, 1)
		})
	)

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.

		var (
			text = c.Text()
		)

		fields := strings.Split(text, " ")

		responseCh := pool.Get()

		ch <- Item{
			Context:   context.Background(),
			Prefix:    fields[0],
			Argument:  strings.Join(fields[1:], " "),
			Responser: newChannelBaseResponser(pool, responseCh),
		}

		result := <-responseCh

		if result.Err != nil {
			c.Send(result.Err.Error())
		} else {
			c.Send(result.Msg)
		}

		return nil
	})

	go bot.Start()

	return &TelegramAdapter{
		bot:  bot,
		ch:   ch,
		pool: pool,
	}
}

func (t *TelegramAdapter) GetChan() <-chan Item {
	return t.ch
}
