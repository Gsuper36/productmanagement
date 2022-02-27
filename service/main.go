package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	pb "ordermanagement/service/ecommerce"

	_ "github.com/jackc/pgx/stdlib"
	"google.golang.org/grpc"
)

func initDB() (*sql.DB, error) {
	driver := os.Getenv("DB_DRIVER")
	url := os.Getenv("DB_URL")

	fmt.Println("is this really works?")

	connection, err := sql.Open(driver, url)

	if err != nil {
		return &sql.DB{}, err
	}

	return connection, nil
}

func main() {
	port := os.Getenv("APP_PORT")

	fmt.Printf("Port: %s\n", port)
	fmt.Println("")
	
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	conn, err := initDB();

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := newServer(conn)
	s := grpc.NewServer()
	pb.RegisterProductManagementServer(s, server)

	log.Printf("Starting gRPC listener on port: %v", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}