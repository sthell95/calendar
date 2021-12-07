package main

import (
	pg "calendar.com/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	c := pg.NewCalendarClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	token, err := c.Login(ctx, &pg.Credentials{
		Login:    "test",
		Password: "testtest",
	})
	if err != nil {
		panic(err.Error())
	}

	log.Printf(`
Token: %v
ExpiresAt: %v`, token.Token, token.ExpiresAt)
}
