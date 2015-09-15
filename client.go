package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	ClientId int    `json:"clientId"`
	Message  string `json:"message"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	go func() {
		for {
			buf := make([]byte, 512)
			_, err := conn.Read(buf)
			if err != nil {
				conn.Close()
				os.Exit(1)
			}
			stringCleaned := bytes.Trim(buf, "\x00")
			var str string = fmt.Sprintf("%s", stringCleaned)
			var message Message
			err = json.Unmarshal([]byte(stringCleaned), &message)
			checkError(err)

			fmt.Println(str)
			fmt.Printf("%d: %s", message.ClientId, message.Message)
			fmt.Print("> ")
		}
	}()

	for {
		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, _ := consoleReader.ReadString('\n')
		input = strings.ToLower(input)
		if strings.HasPrefix(input, "bye") {
			fmt.Println("Good bye!")
			os.Exit(0)
		}
		_, err = conn.Write([]byte(input))
		checkError(err)
	}
	conn.Close()
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
