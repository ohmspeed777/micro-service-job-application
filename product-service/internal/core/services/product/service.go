package product

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/repository/product"
	"context"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/ohmspeed777/go-pkg/errorx"
	"go.mongodb.org/mongo-driver/mongo"
)

type IProductService interface {
	GetAll(ctx context.Context, req GetAllReq) (*domain.ResponseInfo, error)
	Create(ctx context.Context, req CreateReq) (*domain.Product, error)
	GetOne(ctx context.Context, req domain.GetOneReq) (*domain.Product, error)
}

type Service struct {
	product ports.IProductRepo
}

func NewService(db *mongo.Database) IProductService {
	return &Service{
		product: product.NewRepo(db),
	}
}

func (s *Service) GetAll(ctx context.Context, req GetAllReq) (*domain.ResponseInfo, error) {
	products, counter, err := s.product.FindAll(req.Query)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not find all products", err)
	}

	return domain.NewPagination(products, req.GetPage(), req.GetLimit(), counter), nil
}

func (s *Service) Create(ctx context.Context, req CreateReq) (*domain.Product, error) {
	entity := &domain.Product{}
	copier.Copy(entity, &req)

	err := s.product.Create(entity)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not create product", err)
	}

	return entity, nil
}

func (s *Service) GetOne(ctx context.Context, req domain.GetOneReq) (*domain.Product, error) {
	fmt.Println()
	fmt.Println(req.ID)
	fmt.Println()
	
	product, err := s.product.FindOneByID(req.ID)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not find a product", err)
	}

	return product, nil
}