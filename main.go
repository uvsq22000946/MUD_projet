package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) error {
	log.Println("I got a connection")

	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 {
			log.Println("Zero bytes, closing connection")
			break
		}
		msg := string(buf[0 : n-2])
		log.Println("Received message:", msg)

		resp := fmt.Sprintf("Vous avez dis, \"%s\"\r\n", msg)
		n, err = conn.Write([]byte(resp))
		if err != nil {
			return err
		}
		if n == 0 {
			log.Println("Zero bytes, closing connection")
			break
		}
	}

	return nil
}

func startserver() error {
	log.Println("Starting server")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		go func() {
			if err := handleConnection(conn); err != nil {
				log.Println("Error handling connection", err)
				return
			}
		}()
	}
}

func main() {
	err := startserver()
	if err != nil {
		log.Fatal(err)
	}
	println("hello world")
}
