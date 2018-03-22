package main

import (

	"log"
	"google.golang.org/grpc"
	"rloop/Go-Ground-Station/proto"
	"context"
	"fmt"
)

func main() {
	count := 1
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9800", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := proto.NewGroundStationServiceClient(conn)

	stream, err := client.StreamPackets(context.Background(), &proto.StreamRequest{})

	if err == nil{
		for {
			_, err := stream.Recv()
			if err == nil {
				fmt.Printf("\n count: %d \n", count)
				count++
			}
		}
	}else{
		fmt.Println(err)
	}


}
