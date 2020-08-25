package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tsconn23/linotwit/internal/config"
	"github.com/tsconn23/linotwit/internal/bootstrap/flags"
	"github.com/tsconn23/linotwit/internal/bootstrap/interfaces"
)

func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	commonFlags flags.Common,
	handlers []interfaces.BootstrapHandler) {

	wg, _ := initWaitGroup(ctx, cancel, commonFlags, handlers)

	wg.Wait()
}

func initWaitGroup(
	ctx context.Context,
	cancel context.CancelFunc,
	commonFlags flags.Common,
	handlers []interfaces.BootstrapHandler) (*sync.WaitGroup, bool) {

	var cfg config.ConfigInfo
	startedSuccessfully := true
	loader := config.NewLoader(commonFlags)

	err := loader.Process(&cfg)
	if err != nil {
		startedSuccessfully = false
		fmt.Println(err)
	} else {
		fmt.Printf("Access Token = %s\r\n", cfg.Credentials.AccessToken)
	}

	var wg sync.WaitGroup
	// call individual bootstrap handlers.
	if startedSuccessfully {
		translateInterruptToCancel(ctx, &wg, cancel)
		for i := range handlers {
			if handlers[i](ctx, &wg) == false {
				cancel()
				startedSuccessfully = false
				break
			}
		}
	}

	return &wg, startedSuccessfully
}

// translateInterruptToCancel spawns a go routine to translate the receipt of a SIGTERM signal to a call to cancel
// the context used by the bootstrap implementation.
func translateInterruptToCancel(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		signalStream := make(chan os.Signal)
		defer func() {
			signal.Stop(signalStream)
			close(signalStream)
		}()
		signal.Notify(signalStream, os.Interrupt, syscall.SIGTERM)
		select {
		case <-signalStream:
			cancel()
			return
		case <-ctx.Done():
			return
		}
	}()
}