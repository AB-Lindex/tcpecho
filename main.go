package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//go:embed version.txt
var version string

var (
	port     int = 43210
	workerID int
	hostname string
)

func main() {
	if len(os.Args) > 1 {
		v, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		port = v
	}

	hostname, _ = os.Hostname()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s Listening on port %d...", version, port)
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		workerID++
		go handle(conn, workerID)
	}
}

func handle(c net.Conn, id int) {
	log.Printf("%d: starting connection with %v", id, c.RemoteAddr().String())
	defer func() {
		log.Printf("%d: closing connection with %v", id, c.RemoteAddr().String())
		c.Close()
	}()
	fmt.Fprintf(c, "Welcome to tcpecho on %s:%d\nThis is handler #%d, type '.' to quit\n", hostname, port, id)
	rd := bufio.NewScanner(c)
	for {
		if rd.Scan() {
			txt := rd.Text()
			if txt == "." {
				fmt.Fprintln(c, "Goodbye.")
				return
			}
			do(id, c, txt)
		}
	}
}

type dooer func(id int, c net.Conn, cmd, args string)

var commands = map[string]dooer{
	"date": dateCmd,
	"time": timeCmd,
	"info": infoCmd,
}

func do(id int, c net.Conn, txt string) {
	parts := strings.SplitN(txt, " ", 2)
	var args string
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}
	for cmd, fn := range commands {
		if strings.EqualFold(cmd, parts[0]) {
			log.Printf("%d: doing '%s'", id, cmd)
			fn(id, c, cmd, args)
			return
		}
	}
	fmt.Fprintf(c, "unknown command: '%s'\n", parts[0])
}

func dateCmd(id int, c net.Conn, cmd, args string) {
	fmt.Fprintf(c, "The date is %s\n", time.Now().Format(time.DateOnly))
}

func timeCmd(id int, c net.Conn, cmd, args string) {
	fmt.Fprintf(c, "The time is %s\n", time.Now().Format(time.TimeOnly))
}

func infoCmd(id int, c net.Conn, cmd, args string) {
	fmt.Fprintf(c, "You are connecting from %s\n", c.RemoteAddr().String())
}
