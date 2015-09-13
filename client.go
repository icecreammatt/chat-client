package main

import (
	"bufio"
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

	rconn, err := net.Dial("tcp4", service)
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	go func() {
		for {
			buf := make([]byte, 512)
			_, err := rconn.Read(buf)
			if err != nil {
				rconn.Close()
				os.Exit(1)
			}
			var str string = fmt.Sprint(string(buf))
			var message Message
			err = json.Unmarshal([]byte(str), &message)
			// TODO: figure out why message is not being deserialized properly
			// checkError(err)

			fmt.Println(str)
			fmt.Printf("%d: %s\n", message.ClientId, message.Message)
			fmt.Print("> ")
		}
	}()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	// Writer
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	checkError(err)
	if err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(1)
	}
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

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
