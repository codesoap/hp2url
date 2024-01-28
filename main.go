package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type httpline struct {
	Host string `json:"host"`
	Port int    `json:"port,omitempty"`
	TLS  *bool  `json:"tls,omitempty"`
	Req  string `json:"req"`
}

func init() {
	if len(os.Args) > 1 {
		usage := "%s takes httpipe lines via standard input and\n" +
			"prints the correlating URL for each line.\n"
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rawLine := scanner.Bytes()
		var line httpline
		err := json.Unmarshal(rawLine, &line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not parse line '%s': %v\n", rawLine, err)
			os.Exit(1)
		}
		printURL(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: Could not read standard input:", err)
		os.Exit(1)
	}
}

func printURL(line httpline) {
	var u strings.Builder
	tls := line.TLS == nil || *line.TLS == true
	if tls {
		u.WriteString("https://")
	} else {
		u.WriteString("http://")
	}
	u.WriteString(line.Host)
	if line.Port != 0 &&
		(tls && line.Port != 443 || !tls && line.Port != 80) {
		u.WriteString(fmt.Sprint(":", line.Port))
	}
	u.WriteString(extractPath(line.Req))
	fmt.Println(u.String())
}

func extractPath(req string) string {
	i := strings.Index(req, " ")
	if i == -1 || i == len(req)-1 {
		return ""
	}
	j := strings.Index(req[i+1:], " ")
	if j == -1 {
		return ""
	}
	return req[i+1 : i+1+j]
}
