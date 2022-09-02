package flagtrace

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/trace"
)

type traceVar struct {
	traceFile string
	tf        *os.File
}

func (tv *traceVar) String() string {
	if tv == nil {
		return ""
	}
	return tv.traceFile
}

func (tv *traceVar) Set(s string) (err error) {
	tv.traceFile = s
	if len(s) > 0 {
		if tv.tf, err = os.Create(s); err != nil {
			return fmt.Errorf("%s: %s", s, err)
		}
		if err = trace.Start(tv.tf); err != nil {
			return err
		}
		c := make(chan os.Signal, 1)
		go func() {
			for <-c != nil {
				tv.done()
				os.Exit(1)
			}
		}()
		signal.Notify(c, os.Interrupt)
	}
	return nil
}

func (tv *traceVar) Get() interface{} {
	return tv.traceFile
}

func (tv *traceVar) done() {
	if tv.tf == nil {
		return
	}
	defer func() {
		tv.tf.Close()
		tv.tf = nil
	}()
	trace.Stop()
}

// Register registers a [flag] to enable tracing with [runtime/trace].
//
// The returned func must be run defered to stop the tracing and close the file.
//
// When tracing is enabled, a handler for SIGINT (^C) is also registered to properly
// stop tracing and close the file.
func Register(flags *flag.FlagSet, flagName string, usage string) (stopTracing func()) {
	tv := &traceVar{}
	flags.Var(tv, flagName, usage)
	return tv.done
}
