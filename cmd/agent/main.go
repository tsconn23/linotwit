package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/tsconn23/linotwit/internal/bootstrap"
	"github.com/tsconn23/linotwit/internal/bootstrap/flags"
	"github.com/tsconn23/linotwit/internal/bootstrap/interfaces"
)

func main() {
	fmt.Println("Hello world!")

	f := flags.New()
	f.Parse(os.Args[1:])

	ctx, cancel := context.WithCancel(context.Background())
	bootstrap.Run(
		ctx,
		cancel,
		f,
		[]interfaces.BootstrapHandler{
			TestHandler,
		})
}

func TestHandler(
	ctx context.Context,
	wg *sync.WaitGroup) (success bool) {
	wg.Add(1)
	interval := time.Second * time.Duration(3)
	ok := true
	go func(ok *bool) {
		defer wg.Done()
		<-ctx.Done()
		*ok = false
	}(&ok)

	for ok {
		fmt.Println("Program is running...")
		time.Sleep(interval)
	}
	return true
}
