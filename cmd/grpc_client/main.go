package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pg "calendar.com/proto"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	c := pg.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	token, err := c.Login(ctx, &pg.Credentials{
		Login: "test1123",
		//Password: "testtest",
	})
	if err != nil {
		panic(err.Error())
	}

	log.Printf(`
Token: %v
ExpiresAt: %v`, token.Token, token.ExpiresAt)
}
