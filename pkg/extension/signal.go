package extension

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

type StopAble interface {
	Stop()
}

func RegisterStopSignal(ctx context.Context, callback func(), servers ...StopAble) {
	endChan := make(chan os.Signal, 1)
	signal.Notify(endChan, os.Interrupt, os.Kill)
	go func(ctx context.Context) {
		select {
		case <-endChan:
		case <-ctx.Done():
		}
		fmt.Println("Stopping the server...")
		if callback != nil {
			callback()
		}
		for _, server := range servers {
			server.Stop()
		}
	}(ctx)
}
