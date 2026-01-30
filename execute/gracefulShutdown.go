package execute

import (
	"context"
	"os/signal"
	"syscall"
)

var (
	debug = false
)

func CatchExit(ctx context.Context, operations ...func(ctx context.Context)) <-chan struct{} {
	var (
		stop       context.CancelFunc
		notifyChan = make(chan struct{}, 0)
	)

	ctx, stop = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer stop()

		<-ctx.Done()

		logger.Infof(`need exit`)

		for _, operation := range operations {
			logger.Infof(`call operation`)

			if operation != nil {
				operation(ctx)
			}
		}

		stop()

		close(notifyChan)
	}()

	return notifyChan
}

func CatchReload(ctx context.Context, operations ...func(ctx context.Context)) {
	var (
		stop context.CancelFunc
	)

	ctx, stop = signal.NotifyContext(ctx, syscall.SIGUSR1)

	defer stop()

	for _, operation := range operations {
		if operation != nil {
			operation(ctx)
		}
	}

	<-ctx.Done()

	stop()
}
