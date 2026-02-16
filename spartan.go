package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

func main() {
	argv := os.Args[1:]

	if len(argv) < 1 {
		fmt.Println("pass the url")
		os.Exit(1);
	}

	fullUrl := argv[0]

	if !strings.HasPrefix(fullUrl, "spartan://") {
		fullUrl = fmt.Sprintf("spartan://%v", fullUrl)
	}

	u, err := url.Parse(fullUrl)

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "300"
	}

	path := u.EscapedPath()
	if path == "" {
		path = "/"
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(host, port));
	if err != nil {
		log.Fatal(err);
	}

	defer conn.Close()

	request := fmt.Sprintf("%v %v %v\n", u.Hostname(), path, 0)

	conn.Write([]byte(request))

	buf := make([]byte, 4086)

	var res strings.Builder

	for {
		n, err := conn.Read(buf);
		if n > 0 {
			res.Write(buf[:n])
		}

		if err != nil {
			break
		}
	}

	content := res.String()

	if content == "" {
		fmt.Println("no response")
		os.Exit(1);
	}

	lines := strings.Split(content, "\n")
	head := lines[0]

	status := string(head[0])

	if status == "2" {
		fmt.Println(strings.Join(lines[1:], "\n"))
	}
}
