package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	pr "ordermanagement/service/db"
	pb "ordermanagement/service/ecommerce"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type server struct {
	pb.UnimplementedProductManagementServer
	productRepository pr.ProductRepository
}

func (s *server) AddProduct(context context.Context, request *pb.AddProductRequest) (*wrapperspb.StringValue, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return &wrapperspb.StringValue{Value: err.Error()}, err
	}

	s.productRepository.Add(
		&pb.Product{
			Id: id.String(),
			Title: request.GetTitle(),
			Description: request.GetDescription(),
			Price: request.GetPrice(),
		},
	)

	err = s.productRepository.Flush()

	if err != nil {
		log.Println(err)
		fmt.Println(err)

		return &wrapperspb.StringValue{Value: err.Error()}, err
	}

	return &wrapperspb.StringValue{Value: id.String()}, status.New(codes.OK, "").Err()
}

func (s *server) GetProduct(context context.Context, id *wrapperspb.StringValue) (*pb.Product, error) {
	product, err := s.productRepository.FindById(id.Value)

	if err != nil {
		return &pb.Product{}, err
	}

	return product, status.New(codes.OK, "").Err()
}

func (s *server) SearchProducts(*wrapperspb.StringValue, pb.ProductManagement_SearchProductsServer) error {
	return status.Errorf(codes.Unimplemented, "Not implemented")
}

func (s *server) UpdateProducts(pb.ProductManagement_UpdateProductsServer) error {
	return status.Errorf(codes.Unimplemented, "Not implemented")
}

func newServer(conn *sql.DB) *server {
	return &server{
		productRepository: *pr.NewProductRepository(conn),
	}
}