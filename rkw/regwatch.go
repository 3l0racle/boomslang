package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// see: https://docs.microsoft.com/en-us/windows/win32/api/winreg/nf-winreg-regnotifychangekeyvalue

type dword = uint32

const (
	regKeyPath = `SOFTWARE\Wow6432Node\SAAZOD\ManagedPosix`

	REG_NOTIFY_CHANGE_NAME       = dword(0x00000001)
	REG_NOTIFY_CHANGE_ATTRIBUTES = dword(0x00000002)
	REG_NOTIFY_CHANGE_LAST_SET   = dword(0x00000004)
	REG_NOTIFY_CHANGE_SECURITY   = dword(0x00000008)
	REG_NOTIFY_THREAD_AGNOSTIC   = dword(0x10000000)
)

var (
	advapi32                    = syscall.NewLazyDLL("Advapi32.dll")
	procRegNotifyChangeKeyValue = advapi32.NewProc("RegNotifyChangeKeyValue")
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func try(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
	}
}

func main() {
	//k, err := registry.OpenKey(registry.LOCAL_MACHINE, regKeyPath, registry.NOTIFY|registry.QUERY_VALUE|registry.WOW64_64KEY)
	//must(err)

	regEv, err := windows.CreateEvent(nil, 1, 0, nil)
	must(err)
	defer try(windows.Close(regEv))
	//dwFilter := REG_NOTIFY_CHANGE_NAME | REG_NOTIFY_CHANGE_ATTRIBUTES | REG_NOTIFY_CHANGE_LAST_SET | REG_NOTIFY_CHANGE_SECURITY

	queue := make(chan error)
	ctx, _ := cancelHandler()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go notifyRegChange(ctx, registry.LOCAL_MACHINE, regKeyPath, queue)
	go func() {
		fmt.Println("consumer start")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("consumer dead")
				return
			case err, ok := <-queue:
				if !ok {
					fmt.Printf("consumer chan closed\n")
				}
				fmt.Printf("consumer got data %s\n", err)
			}
		}
	}()

	//go regListen(ctx, wg, regEv)
	wg.Wait()
}

func notifyRegChange(ctx context.Context, key registry.Key, path string, notifyCh chan error) (err error) {
	k, err := registry.OpenKey(key, path, syscall.KEY_NOTIFY)
	if err != nil {
		return
	}
	for {
		select {
		case <-ctx.Done():
			close(notifyCh)
			return
		default:
		}
		procRegNotifyChangeKeyValue.Call(uintptr(k), 1, 0x00000001|0x00000004, 0, 0)
		notifyCh <- nil
	}
}

// func regListen(ctx context.Context, wg *sync.WaitGroup, h windows.Handle) {
// 	defer wg.Done()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("regListen: ctx.Done()")
// 			return
// 		default:
// 		}

// 		ev, err := windows.WaitForSingleObject(h, 5000)
// 		if err != nil {
// 			fmt.Printf("regListen: error - %s\n", err)
// 			return
// 		}

// 		handleEvent(ev)
// 	}
// }

// func handleEvent(ev uint32) {
// 	switch ev {
// 	case uint32(windows.WAIT_TIMEOUT):
// 		fmt.Println("handleEvent: timeout")
// 		return
// 	case windows.WAIT_OBJECT_0:
// 		windows.ResetEvent(windows.Handle(ev))
// 		fmt.Println("handleEvent: change detected")
// 	default:
// 		fmt.Printf("handleEvent: unknown event %#010x\n", ev)
// 	}
// }

// func testFlags() {
// 	dwFilter := REG_NOTIFY_CHANGE_NAME | REG_NOTIFY_CHANGE_ATTRIBUTES | REG_NOTIFY_CHANGE_LAST_SET | REG_NOTIFY_CHANGE_SECURITY
// 	keys := map[string]dword{
// 		"REG_NOTIFY_CHANGE_NAME":       REG_NOTIFY_CHANGE_NAME,
// 		"REG_NOTIFY_CHANGE_ATTRIBUTES": REG_NOTIFY_CHANGE_ATTRIBUTES,
// 		"REG_NOTIFY_CHANGE_LAST_SET":   REG_NOTIFY_CHANGE_LAST_SET,
// 		"REG_NOTIFY_CHANGE_SECURITY":   REG_NOTIFY_CHANGE_SECURITY,
// 	}

// 	for k, v := range keys {
// 		fmt.Printf("%s:\t%#010x\n", k, v)
// 	}

// 	fmt.Printf("\nWant Flag:\t%d\t%#010x\n", dwFilter, dwFilter)
// }

// cancelHandler returns cancellation context and function for graceful shutdown
func cancelHandler() (context.Context, context.CancelFunc) {
	ctx, cancelFn := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		cancelFn()
		signal.Stop(signals)
	}()

	return ctx, cancelFn
}
