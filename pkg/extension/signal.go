package extension

import (
	"context"
	"log"
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
		if callback != nil {
			callback()
		}
		for _, server := range servers {
			log.Println("Stopping the server...")
			server.Stop()
		}
	}(ctx)
}
