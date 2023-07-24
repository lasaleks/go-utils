package ie_common_utils_go

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func CreatePidFile(pid_file string) {
	file, err := os.Create(pid_file)
	if err != nil {
		panic(fmt.Errorf("Error create pid file %s", err))
	}
	fmt.Fprint(file, os.Getpid())
	file.Close()
}

func WaitSignalExit(wg *sync.WaitGroup, ctx context.Context, callback func(ctx context.Context)) {
	// notifying the main goroutine that we are done
	defer wg.Done()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Signal(syscall.SIGTERM), os.Interrupt)
	n := <-stop
	fmt.Println("signal:", n)
	callback(ctx)
}
