package execute

import (
	"context"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCatchExit(t *testing.T) {
	var (
		ctx  = context.Background()
		exit = int64(0)
	)

	debug = true

	t.Log(`enter`, time.Now().String())

	ctx, _ = context.WithTimeout(ctx, time.Second*5)

	exitFunc := func(ctx context.Context) {
		atomic.AddInt64(&exit, 1)
	}

	wait := CatchExit(ctx, exitFunc)

	go func() {
		time.Sleep(time.Second)
		process, err := os.FindProcess(os.Getpid())

		require.NoError(t, err, `findProcess`)

		require.NoError(t, process.Signal(syscall.SIGINT), `send signal`)
	}()

	<-wait

	t.Log(`sleep`, time.Now().String())

	t.Log(`exitValue`, atomic.LoadInt64(&exit))
}
