package main

import (
	"fmt"
	"greet/greet/greetpb"
	"log"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("W: Failed creating Client Connection %v", err)
	}

	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	doUnaryRequest(c)

}

func doUnaryRequest(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greet: &greetpb.Greet{
			FirstName: "Mohamed",
			LastName:  "Fadel",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("W: Failed making request %v", err)
	}
	fmt.Println("Result: " + res.GetResult())
}
