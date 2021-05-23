package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

type MessageEvent struct {
	msg string
}

type UserJoinedEvent struct {
}

type ClientInput struct {
	user  *User
	event interface{}
}

type User struct {
	name    string
	session *Session
}

type Session struct {
	conn net.Conn
}

func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
}

type World struct {
	users []*User
}

func generationName() string {
	return fmt.Sprintf("User %d", rand.Intn(100)+1)
}

func handleConnection(conn net.Conn, inputChannel chan ClientInput) error {
	buf := make([]byte, 4096)

	session := &Session{conn}
	user := &User{name: generationName(), session: session}

	inputChannel <- ClientInput{
		user,
		&UserJoinedEvent{},
	}

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

		e := ClientInput{user, &MessageEvent{msg}}
		inputChannel <- e

	}
	return nil
}

func startserver(eventChannel chan ClientInput) error {
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
			if err := handleConnection(conn, eventChannel); err != nil {
				log.Println("Error handling connection", err)
				return
			}
		}()
	}
}

func startGameLoop(clientInputChannel <-chan ClientInput) {
	w := &World{}
	for input := range clientInputChannel {
		switch event := input.event.(type) {
		case *MessageEvent:
			fmt.Println("Message recu", event.msg)
			// TODO: error handing
			input.user.session.WriteLine(fmt.Sprintf("Vous avez dis, \"%s\"\r\n", event.msg))
			for _, user := range w.users {
				if user != input.user {
					user.session.WriteLine(fmt.Sprintf("%s a dis, \"%s\"", input.user.name, event.msg))
				}
			}
		case *UserJoinedEvent:
			fmt.Println("User joined :", input.user.name)
			w.users = append(w.users, input.user)
			input.user.session.WriteLine(fmt.Sprintf("Bienvenue %s", input.user.name))
			for _, user := range w.users {
				if user != input.user {
					user.session.WriteLine(fmt.Sprintf("%s entre dans la salle", input.user.name))
				}
			}
		}
	}
}

func main() {

	ch := make(chan ClientInput)

	go startGameLoop(ch)

	err := startserver(ch)
	if err != nil {
		log.Fatal(err)

	}
}
