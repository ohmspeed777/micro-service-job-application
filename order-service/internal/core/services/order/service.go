package order

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	up "app/internal/core/proto/user"
	"app/repository/order"
	"app/repository/product"
	"context"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/ohmspeed777/go-pkg/corex"
	"github.com/ohmspeed777/go-pkg/errorx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type IOrderService interface {
	Create(ctx context.Context, req CreateReq) (*domain.Order, error)
	Cancel(ctx context.Context, req domain.GetOneReq) (*domain.CommonResponse, error)
	GetOne(ctx context.Context, req domain.GetOneReq) (*domain.OrderLookedUp, error)
	GetAll(ctx context.Context, req GetAllReq) (*domain.ResponseInfo, error)
}

type Service struct {
	product  ports.IProductRepo
	order    ports.IOrderRepo
	userGRPC up.UserClient
}

func NewService(db *mongo.Database, userGRPC *grpc.ClientConn) IOrderService {
	return &Service{
		product:  product.NewRepo(db),
		order:    order.NewRepo(db),
		userGRPC: up.NewUserClient(userGRPC),
	}
}

func (s *Service) Create(ctx context.Context, req CreateReq) (*domain.Order, error) {
	entity := &domain.Order{}
	copier.Copy(entity, &req)

	c, err := corex.NewFromOutingContext(ctx)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not bind my context", err)
	}

	_id, err := primitive.ObjectIDFromHex(c.User.ID)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not convert to primitive object", err)
	}

	entity.User = _id
	entity.Status = domain.Created

	// !to do validate stock
	// !should calculate by own itself

	err = s.order.Create(entity)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not create order", err)
	}

	return entity, nil
}

func (s *Service) Cancel(ctx context.Context, req domain.GetOneReq) (*domain.CommonResponse, error) {
	c, err := corex.NewFromOutingContext(ctx)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not bind my context", err)
	}

	order, err := s.order.FindOneByID(req.ID)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not find a order", err)
	}

	// validate order
	if order.User.Hex() != c.User.ID {
		return nil, errorx.New(http.StatusBadRequest, "you didn't have permission")
	}

	order.Status = domain.Canceled
	err = s.order.Update(order)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not create order", err)
	}

	return domain.NewOkMessage(), nil
}

func (s *Service) GetOne(ctx context.Context, req domain.GetOneReq) (*domain.OrderLookedUp, error) {
	order, err := s.order.AggregateOneByID(req.ID)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not find a order", err)
	}

	res, err := s.userGRPC.GetProfile(ctx, &up.GetUserReq{
		Id: order.User.Hex(),
	})

	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not retrieve a user", err)
	}

	u := &domain.User{
		Email:     res.GetEmail(),
		FirstName: res.GetFirstName(),
		LastName:  res.GetLastName(),
	}

	u.ID = order.User
	u.CreatedAt = res.GetCreatedAt().AsTime()
	u.UpdatedAt = res.GetUpdatedAt().AsTime()

	return order.Format(u), nil
}

func (s *Service) GetAll(ctx context.Context, req GetAllReq) (*domain.ResponseInfo, error) {
	orders, counter, err := s.order.AggregateAllByUser(req.User, req.Query)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not find orders", err)
	}

	// format
	for _, v := range orders {
		v.Format(nil)
	}

	return domain.NewPagination(orders, req.GetPage(), req.GetLimit(), counter), nil
}
