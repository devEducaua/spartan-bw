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
    res := request(host, path, port);

    parseResult(host, res);
}

func request(host string, path string, port string) string {
    
    conn, err := net.Dial("tcp", net.JoinHostPort(host, port));
    if err != nil {
        panic(err);
    }

    defer conn.Close();

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

    return res.String();
}

func parseUrl(spartanUrl string) (string, string, string) {
    if !strings.HasPrefix(spartanUrl, "spartan://") {
        spartanUrl = "spartan://" + spartanUrl;
    }

    u, err := url.Parse(strings.TrimSpace(spartanUrl));
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

func parseResult(host string, result string) {
    if result == "" {
        fmt.Println("no response");
        return
    }

    lines := strings.Split(result, "\n");
    head := lines[0];

    parts := strings.SplitN(head, " ", 2);
    code := parts[0];

    meta := "";
    if len(parts) > 1 {
        meta = parts[1];
    }

    switch code {
    case "2": 
		body := strings.Join(lines[1:], "\n");
		fmt.Println(body);
    case "3":
        host, port, path := parseUrl(host+meta);
        req := request(host, path, port);
        parseResult(host, req);
    case "4": 
        fmt.Println("server returned status 4, client-error: ", meta);
    case "5": 
        fmt.Println("server returned status 5, server-error: ", meta);
    }
}
