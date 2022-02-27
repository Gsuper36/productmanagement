package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "productmanagement/client/ecommerce"

	"github.com/jaswdr/faker"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var host = flag.String("host", "", "Example localhost:port")

func main() {

	flag.Parse()
	command := make(chan []byte)
	faker := faker.New()
	if *host == "" {
		log.Fatal("Host is empty")
	}

	conn, err := grpc.Dial(*host, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect to: %v", err)
	}

	defer conn.Close()

	client := pb.NewProductManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)

	go readStdin(command)

	for {
		com := string(<-command)

		switch {
			case com == "exit":
				os.Exit(0)
			case com == "add":
				req := &pb.AddProductRequest{
					Title: faker.Lorem().Word(),
					Description: faker.Lorem().Sentence(3),
					Price: faker.Int64Between(10, 100000),
				}

				fmt.Println(req)

				res, err := client.AddProduct(ctx, req)
			if err != nil {
				log.Println(err, res)
			} else {
				log.Println(res.Value)
			}

			case com == "get":
				fmt.Println("Enter id")

				id := string(<-command)

				product, err := client.GetProduct(ctx, &wrapperspb.StringValue{
					Value: id,
				})

				if err != nil {
					log.Println(err)
				} else {
					log.Println(product)
				}

		}
		cancel()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second * 15)
	}
}

func readStdin(out chan<- []byte) {
	for {
		data := make([]byte, 1024)
		l, _ := os.Stdin.Read(data)
		var size int

		for i, v := range data {
			if v == 10 {
				size = i
				break
			}
		}

		if l > 0 {
			out <- data[:size]
		}
	}
}