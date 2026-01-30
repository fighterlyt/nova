package conversation

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const (
	initCapacity = 100
)

type Service interface {
	Register(prefix string, op Handler) error
	Serve() error
	Process(ctx context.Context, prefix, argument string) (string, error)
}

// Handler 处理方法,传入上下文，参数，返回结果和错误
type Handler func(ctx context.Context, argument string) (string, error)

type service struct {
	handlers map[string]Handler
	lock     *sync.RWMutex
	notFound Handler
	Adapter  Adapter
}

func newService(adapter Adapter) *service {
	return &service{
		handlers: make(map[string]Handler, initCapacity),
		lock:     &sync.RWMutex{},
		Adapter:  adapter,
	}
}

func (s *service) Register(prefix string, op Handler) error {
	s.lock.Lock()

	defer s.lock.Unlock()

	if prefix == `` {
		return errors.New(`prefix must not be empty`)
	}

	if _, exists := s.handlers[prefix]; exists {
		return fmt.Errorf(`prefix[%s] already exists`, prefix)
	}

	s.handlers[prefix] = op

	return nil
}

func (s *service) Serve() error {
	go func() {
		for msg := range s.Adapter.GetChan() {
			go s.serve(msg.Context, &msg)
		}
	}()

	return nil
}

func (s *service) serve(ctx context.Context, item *Item) {
	msg, err := s.Process(ctx, item.Prefix, item.Argument)

	item.Responser.Response(msg, err)
}

func (s *service) Process(ctx context.Context, prefix, argument string) (string, error) {
	s.lock.RLock()

	handler := s.handlers[prefix]

	s.lock.RUnlock()

	if handler == nil {
		handler = s.notFound
	}

	return handler(ctx, argument)
}
