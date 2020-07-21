package main

import (
	"fmt"
	"greet/greet/greetpb"
	"io"
	"log"
	"time"

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
	//doServerStreamingRequest(c)
	doClientStreaming(c)
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Akari",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Neji",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nekoma",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Shoyo",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("W: Failed calling LongGreet: %v", err)
	}
	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("W: Error geting response from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet response: %v\n", resp.GetResult())
}
