package main

import (
	"fmt"
	"greet/greet/greetpb"
	"io"
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

	// rpc calls
	//doUnaryRequest(c)
	doServerStreamingRequest(c)
}

func doUnaryRequest(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
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

func doServerStreamingRequest(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mohamed",
			LastName:  "Fadel",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("W: Failed calling GreetManyTimes: %v", err)
	}
	for {
		resp, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("W: Failed Reading from stream: %v", err)
		}
		fmt.Printf("Response from GreetManyTimes %v \n", resp.GetResult())
	}
}
