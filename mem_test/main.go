package main

import (
	"context"
	"fmt"
	"time"

	poContext "github.com/disiqueira/PoContext/context"
)

// This is a small test to see how the log pool works when a lot of new instances are created.
// The following code creates 200000 logger instances, but analysing the memory address we can see
// that only 73318 different address are used, so ~63% of the times the logger was reused,
// not allocating more memory.
//
// References:
// https://golang.org/pkg/sync/#Pool
// https://github.com/sirupsen/logrus/blob/master/logger.go#L78-L84
func main() {
	baseContext := context.Background()

	for i := 0; i < 100000; i++ {
		i := i
		go func() {
			ctx := poContext.WithTraceID(baseContext, fmt.Sprintf("TraceID_%d", i))

			fmt.Printf("%p\n", poContext.Logger(ctx))
			fmt.Printf("%p\n", poContext.Logger(ctx))

		}()
	}

	time.Sleep(5 * time.Second)
}
