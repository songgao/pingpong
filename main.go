package main

import (
	"flag"
	"fmt"
	"os"
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
	for i := 0; i < len(cmds); i++ {
		cmd := cmds[i]
		startProcess(cmd, r.In[i], func(err error) {
			msg := ""
			if err == nil {
				msg = fmt.Sprintf("Command [%s] exited with code 0", cmd)
			} else {
				msg = fmt.Sprintf("Command [%s] exited with error %v", cmd, err)
			}
			r.PrintStatus([]byte(msg), i)
		})
	}
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
