package order

import (
	"app/internal/core/domain"
	pb "app/internal/core/proto/order"
	"app/internal/core/services/order"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPC struct {
	order order.IOrderService
	pb.UnimplementedOrderServiceServer
}

func NewGrpcServer(order order.IOrderService, server *grpc.Server) {
	pb.RegisterOrderServiceServer(server, &GRPC{
		order: order,
	})
}

func (g *GRPC) GetMyOrder(ctx context.Context, req *pb.GetMyOrderReq) (*pb.GetMyOrderRes, error) {
	res, err := g.order.GetAll(ctx, order.GetAllReq{User: req.Id})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := []*pb.Order{}
	for _, v := range res.Entities.([]*domain.OrderLookedUp) {
		tmp := &pb.Order{
			Id:        v.ID.Hex(),
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
			Status:    int64(v.Status),
			Total:     v.Total,
		}

		items := []*pb.OrderItem{}
		for _, item := range v.Items {
			items = append(items, &pb.OrderItem{
				Quantity:  int64(item.Quantity),
				ProductId: item.ProductID.Hex(),
				Product: &pb.Product{
					Id:        item.ProductID.Hex(),
					CreatedAt: timestamppb.New(item.Product.CreatedAt),
					UpdatedAt: timestamppb.New(item.Product.UpdatedAt),
					Price:     item.Product.Price,
					Name:      item.Product.Name,
					Stock:     int64(item.Product.Stock),
				},
			})
		}
		tmp.Items = items
		response = append(response, tmp)
	}

	return &pb.GetMyOrderRes{
		PageInfo: &pb.PageInfo{
			Page:          res.PageInfo.Page,
			Size:          res.PageInfo.Size,
			AllOfEntities: res.PageInfo.AllOFEntities,
			NumOfPages:    res.PageInfo.NumOfPages,
		},
		Entities: response,
	}, nil
}
