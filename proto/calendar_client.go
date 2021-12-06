package proto

import "google.golang.org/grpc"

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	c := pg.New
}
