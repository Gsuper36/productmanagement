package db

import (
	"context"
	"database/sql"
	"log"
	pb "ordermanagement/service/ecommerce"
)

type Repository interface {
	FindById(string) (*pb.Product, error)
	All() ([]*pb.Product, error)
	Flush() error
	Add(*pb.Product)
}

type ProductRepository struct {
	connection *sql.DB
	products []*pb.Product
}

func (rep *ProductRepository) FindById(id string) (*pb.Product, error)  {
	row := rep.connection.QueryRow("select * from product where id=$1", id)
	product := new(pb.Product);

	err := row.Scan(&product.Id, &product.Title, &product.Description, &product.Price)

	if err != nil {
		return &pb.Product{}, err
	}

	return product, nil
}

func (rep *ProductRepository) All() ([]*pb.Product, error) {
	products := make([]*pb.Product, 0)

	rows, err := rep.connection.Query("select * from product")

	if err != nil {
		return products, err
	}

	defer rows.Close()

	for rows.Next() {
		product := new(pb.Product)

		err := rows.Scan(&product.Id, &product.Title, &product.Description, &product.Price)

		if err != nil {
			return make([]*pb.Product, 0), err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return make([]*pb.Product, 0), err
	}

	return products, nil
}

func (rep *ProductRepository) Flush() error {
	if (len(rep.products) == 0) {
		return nil
	}

	tx, err := rep.connection.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelDefault})

	if err != nil {
		log.Println(err)
		return err
	}

	for _, product := range rep.products {
		log.Println(product)
		_, err := tx.Exec("insert into product (id, title, description, price) values ($1, $2, $3, $4)", &product.Id, &product.Title, &product.Description, &product.Price)
		
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (rep *ProductRepository) Add(product *pb.Product) {
	rep.products = append(rep.products, product)
}

func NewProductRepository(conn *sql.DB) *ProductRepository {
	return &ProductRepository{
		connection: conn,
		products: make([]*pb.Product, 0),
	}
}