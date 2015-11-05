package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")

	if err != nil {
		panic("Error loading X509 key pair")
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	conn, err := tls.Dial("tcp", service, &config)
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	go func() {
		for {
			buf := make([]byte, 1500)
			_, err := conn.Read(buf)
			if err != nil {
				conn.Close()
				os.Exit(1)
			}
			stringCleaned := bytes.Trim(buf, "\x00")
			var str string = fmt.Sprintf("%s", stringCleaned)
			var message Message
			err = json.Unmarshal([]byte(stringCleaned), &message)
			if err != nil {
				fmt.Println("Error parsing JSON", err.Error())
				fmt.Println("Raw Message: ", str)
			} else {
				fmt.Println(str)
				fmt.Printf("%d: %s", message.ClientId, message.Message)
				fmt.Print("> ")

			}

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
