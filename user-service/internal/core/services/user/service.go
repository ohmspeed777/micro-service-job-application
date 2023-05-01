package user

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	pb "app/internal/core/proto/order"
	"app/repository/user"
	"context"
	"net/http"

	"github.com/ohmspeed777/go-pkg/corex"
	"github.com/ohmspeed777/go-pkg/errorx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type IUserService interface {
	GetMe(ctx context.Context) (*domain.User, error)
	GetOne(ctx context.Context, req domain.GetOneReq) (*domain.User, error)
	GetMyOrderHistory(ctx context.Context) (*domain.ResponseInfo, error)
}

type Service struct {
	user      ports.IUserRepo
	orderGRPC pb.OrderServiceClient
}

func NewService(db *mongo.Database, orderGRPC *grpc.ClientConn) IUserService {
	return &Service{
		user:      user.NewRepo(db),
		orderGRPC: pb.NewOrderServiceClient(orderGRPC),
	}
}

func (s *Service) GetMe(ctx context.Context) (*domain.User, error) {
	c, err := corex.NewFromOutingContext(ctx)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not bind my context", err)
	}

	user, err := s.user.FindOneByID(c.User.ID)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not find my profile", err)
	}

	return user, nil
}

func (s *Service) GetOne(ctx context.Context, req domain.GetOneReq) (*domain.User, error) {
	user, err := s.user.FindOneByID(req.ID)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not find my profile", err)
	}

	return user, nil
}

func (s *Service) GetMyOrderHistory(ctx context.Context) (*domain.ResponseInfo, error) {
	c, err := corex.NewFromOutingContext(ctx)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not bind my context", err)
	}

	res, err := s.orderGRPC.GetMyOrder(ctx, &pb.GetMyOrderReq{Id: c.User.ID})
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, err.Error(), err)
	}

	entities := []*domain.Order{}
	for _, v := range res.Entities {
		_id, _ := primitive.ObjectIDFromHex(v.GetId())
		_userID, _ := primitive.ObjectIDFromHex(c.User.ID)

		tmp := &domain.Order{
			Status: domain.OrderStatus(v.Status),
			Total:  v.Total,
		}
		tmp.ID = _id
		tmp.CreatedAt = v.GetCreatedAt().AsTime()
		tmp.UpdatedAt = v.GetUpdatedAt().AsTime()
		tmp.User = _userID

		items := []*domain.OrderItem{}
		for _, item := range v.Items {
			p := &domain.Product{
				Price: item.GetProduct().Price,
				Stock: uint(item.GetProduct().Stock),
				Name:  item.GetProduct().Name,
			}

			p.CreatedAt = item.GetProduct().GetCreatedAt().AsTime()
			p.UpdatedAt = item.GetProduct().GetUpdatedAt().AsTime()
			_pid, _ := primitive.ObjectIDFromHex(item.GetProduct().GetId())
			p.ID = _pid
			items = append(items, &domain.OrderItem{
				Quantity: uint(item.Quantity),
				Product:  p,
			})
		}
		tmp.Items = items
		entities = append(entities, tmp)

	}

	return domain.NewPagination(entities, res.PageInfo.Page, res.PageInfo.Size, res.PageInfo.AllOfEntities), nil
}
