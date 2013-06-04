package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	printTime bool
	logDir    string
	help      bool
)

func init() {
	flag.BoolVar(&printTime, "time", true, "Prepend time before each line")
	flag.StringVar(&logDir, "log", "", "Directory for logging")
	flag.BoolVar(&help, "h", false, "Print help message; an empty string means no logging")
	flag.BoolVar(&help, "help", false, "Print help message; an empty string means no logging")
}

func main() {
	flag.Parse()
	if help {
		printHelp()
	} else {
		if logDir != "" {
			os.MkdirAll(logDir, os.ModeDir|0755)
		}
		run(flag.Args())
	}
}

func run(cmds []string) {
	r := NewRender(len(cmds), printTime)
	processes := make([]*os.Process, len(cmds))
	for i := 0; i < len(cmds); i++ {
		i_ := i
		cmd := cmds[i]
		processes[i] = startProcess(cmd, r.In[i], func(err error) {
			msg := ""
			if err == nil {
				msg = fmt.Sprintf("Command [%s] exited with code 0", cmd)
			} else {
				msg = fmt.Sprintf("Command [%s] exited with error %v", cmd, err)
			}
			r.PrintStatus([]byte(msg), i_)
		})
	}

	// catch Interrupt and Kill signals, making sure when pingpong exits all sub-processes exit.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)
	go func() {
		sig := <-ch
		fmt.Println() // prevent "^C" from ruining output
		r.PrintStatus([]byte(fmt.Sprintf("Got signal %v; killing processes...", sig)), -1)
		for _, proc := range processes {
			if proc != nil {
				proc.Kill()
			}
		}
	}()

	fmt.Println("")
	r.PrintLegend(cmds)
	fmt.Println("")
	r.Run()
	fmt.Println("")
	r.PrintLegend(cmds)
	fmt.Println("")
}

func printHelp() {
	fmt.Printf("Usage: %s [OPTIONS] -- cmd1, cmd2, ...\n", os.Args[0])
	flag.PrintDefaults()
}
