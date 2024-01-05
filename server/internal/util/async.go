package util

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func EnsureRunGoroutin(f func(), tryCount ...int) {
	try := 0
	if len(tryCount) > 0 {
		try = tryCount[0]
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("paniced: %v", r)
				fmt.Printf("%s", string(debug.Stack()))
				time.Sleep(1 * time.Second)
				if try > 10 {
					os.Exit(1)
				}
				EnsureRunGoroutin(f, try+1)
			}
		}()

		f()
	}()
}
