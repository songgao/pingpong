# pingpong

`pingpong` is a utility commandline tool that makes testing server/client programs easier by interleaving output from different commands and labeling them with different colors.

## Installation
```shell
$ go get -u github.com/songgao/pingpong
```

## Usage
```
Usage: pingpong [OPTIONS] -- cmd1, cmd2, ...
  -h=false: Print help message; an empty string means no logging
  -help=false: Print help message; an empty string means no logging
  -log="": Directory for logging
  -time=true: Prepend time before each line
```
Each `cmd` is fed into `bash -c` to be executed.

## Example

Here's a simple piece of code composed with server and clients:
```go
package main

import (
  "net"
  "time"
  "fmt"
  "os"
  "io"
  "strconv"
)

const STR = `Go is expressive, concise, clean, and efficient. Its concurrency mechanisms make it easy to write programs that get the most out of multicore and networked machines, while its novel type system enables flexible and modular program construction. Go compiles quickly to machine code yet has the convenience of garbage collection and the power of run-time reflection. It's a fast, statically typed, compiled language that feels like a dynamically typed, interpreted language.`

func main() {
  switch os.Args[1] {
  case "server": server()
  case "client": client()
  default:
  }
}

func server() {
  ln, err := net.Listen("tcp", ":9090")
  if err != nil{ return }
  go func() {
    flag := true
    for {
      conn, _ := ln.Accept()
      if err != nil { continue }
      fmt.Printf("-- Client connected: %s\n", conn.RemoteAddr())
      fmt.Fprintln(conn, "Hello client, greetings from server!")
      if flag {
        fmt.Fprintf(conn, "%s\n", STR)
        flag = false
      }
      go func() {
        io.Copy(os.Stdout, conn)
        fmt.Printf("-- Client disconnected: %s\n", conn.RemoteAddr())
      }()
    }
  }()
  time.Sleep(8 * time.Second)
}

func client() {
  id, _ := strconv.Atoi(os.Args[2])
  time.Sleep(2 * time.Duration(id) * time.Second)
  conn, err := net.Dial("tcp", "localhost:9090")
  if err != nil { return }
  go io.Copy(os.Stdout, conn)
  fmt.Fprintf(conn, "Hello server, greetings from client %d!\n", id)
  time.Sleep(time.Second)
  conn.Close()
}
```

Running 4 commands with `pingpong`:
```shell
$ pingpong -log="./logs" "go run main.go server" "go run main.go client 1" "go run main.go client 2" "exit 1"
```

Screenshot:

![](http://songgao.github.io/pingpong/images/screenshot.png)
