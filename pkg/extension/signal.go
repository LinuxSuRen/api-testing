package extension

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/linuxsuren/api-testing/pkg/logging"
)

var (
	signalLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("signal")
)

type StopAble interface {
	Stop()
}

func RegisterStopSignal(ctx context.Context, callback func(), servers ...StopAble) {
	endChan := make(chan os.Signal, 1)
	signal.Notify(endChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func(ctx context.Context) {
		select {
		case <-endChan:
		case <-ctx.Done():
		}
		if callback != nil {
			callback()
		}
		for _, server := range servers {
			signalLogger.Info("Stopping the server...")
			server.Stop()
		}
	}(ctx)
}
