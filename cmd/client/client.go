package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gustavodaltoe/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // conecta com o server
	if err != nil {
		log.Fatalf("could not connect to gRPC Server: %v", err)
	}
	defer connection.Close() // defer observa quando ninguem está usando a connection e ela é fechada

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Gustavo",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Gustavo",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not receive the message: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "g1",
			Name:  "gustavo",
			Email: "gus@gus.com",
		},
		&pb.User{
			Id:    "g2",
			Name:  "gustavo2",
			Email: "gus2@gus.com",
		},
		&pb.User{
			Id:    "g3",
			Name:  "gustavo3",
			Email: "gus3@gus.com",
		},
		&pb.User{
			Id:    "g4",
			Name:  "gustavo4",
			Email: "gus4@gus.com",
		},
		&pb.User{
			Id:    "g5",
			Name:  "gustavo5",
			Email: "gus5@gus.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "g1",
			Name:  "gustavo",
			Email: "gus@gus.com",
		},
		&pb.User{
			Id:    "g2",
			Name:  "gustavo2",
			Email: "gus2@gus.com",
		},
		&pb.User{
			Id:    "g3",
			Name:  "gustavo3",
			Email: "gus3@gus.com",
		},
		&pb.User{
			Id:    "g4",
			Name:  "gustavo4",
			Email: "gus4@gus.com",
		},
		&pb.User{
			Id:    "g5",
			Name:  "gustavo5",
			Email: "gus5@gus.com",
		},
	}

	// channel é um local onde é mandada uma comunicação entre uma go-rotina e outra
	wait := make(chan int)

	go func() { // fica rodando em background enviando a informação
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() { // em paralelo, fica recebendo a resposta
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait) // finaliza o channel para finalizar o programa
	}()

	<-wait // enquanto esse channel não morrer, a aplicação não para
}
