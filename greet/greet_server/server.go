package main

import (
	"context"
	"fmt"
	"greet/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("Greet function call")
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GreetManyTimes function call")
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " for the \"" + strconv.Itoa(i) + "\"th time"
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function call")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// reached the end of the stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("W: Failed reading from client stream: %v", err)
		}
		result += "Hello " + req.GetGreeting().GetFirstName() + "!\t"
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone rpc Call")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// client closed stream should stop streaming data
			return nil
		}
		if err != nil {
			log.Fatalf("W: Error reading data from client stream: %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!\n"
		sendErr := stream.Send(&greetpb.GreetEveryoneRespone{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("W: Error sending data to client stream: %v", sendErr)
			return sendErr
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("W:Failed creating listner: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("W:Failed creating server: %v", err)
	}
}
