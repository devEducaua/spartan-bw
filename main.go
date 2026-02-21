package main

import (
    "fmt"
    "net"
    "net/url"
    "os"
    "strings"
)

func main() {
    argv := os.Args[1:];

    if len(argv) < 1 {
        fmt.Println("ERROR: pass the url");
        os.Exit(1);
    }

    spartanUrl := argv[0];
    host, port, path := parseUrl(spartanUrl);

    conn, err := net.Dial("tcp", net.JoinHostPort(host, port));
    if err != nil {
        panic(err);
    }

    defer conn.Close();

    res := request(conn, host, path);

    parseResult(conn, host, res);
}

func request(conn net.Conn, host string, path string) strings.Builder {

    request := fmt.Sprintf("%v %v %v\n", host, path, 0);
    conn.Write([]byte(request));

    buf := make([]byte, 4096);
    var res strings.Builder;

    for {
        n, err := conn.Read(buf);
        if n > 0 {
            res.Write(buf[:n]);
        }
        if err != nil {
            break;
        }
    }

    return res;
}

func parseUrl(spartanUrl string) (string, string, string) {
    if !strings.HasPrefix(spartanUrl, "spartan://") {
        spartanUrl = "spartan://" + spartanUrl;
    }

    u, err := url.Parse(spartanUrl);
    if err != nil {
        panic(err);
    }

    host := u.Hostname();
    port := u.Port();
    if port == "" {
        port = "300";
    }

    path := u.EscapedPath();
    if path == "" {
        path = "/";
    }

    if strings.HasSuffix(path, ".gmi/") {
		path = path[:len(path)-1]
    }

    return host, port, path;
}

func parseResult(conn net.Conn, host string, result strings.Builder) {
    content := result.String();

    lines := strings.Split(content, "\n");

    if content == "" {
        fmt.Println("TODO: for whatever reason has no response");
        return
    }

    head := lines[0];
    statusCode := string(head[0]);

    switch statusCode {
    case "2": 
        fmt.Println(content);
        fmt.Println("2 success request");
    case "3":
        redirectPath := head[2:];
		fmt.Println(redirectPath);
		fmt.Println(host);

		result := request(conn, host, redirectPath);

		fmt.Println(result)

		parseResult(conn, host, result)
    case "4": 
        error := head[2:];
        fmt.Println("server returned status 4, client-error: ", error);
    case "5": 
        error := head[2:];
        fmt.Println("server returned status 5, server-error: ", error);
    }
}
