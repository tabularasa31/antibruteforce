package grpcv1

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	proto "github.com/tabularasa31/antibruteforce/api"
	"github.com/tabularasa31/antibruteforce/internal/models"
	"github.com/tabularasa31/antibruteforce/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net"
	"time"
)

type AntibruteforceService struct {
	useCases usecase.UseCases
	logg     zap.Logger
	proto.UnimplementedAntiBruteforceServer
}

func NewAntibruteforceService(u usecase.UseCases, logg zap.Logger) *AntibruteforceService {
	return &AntibruteforceService{useCases: u, logg: logg}
}

func (a *AntibruteforceService) AllowRequest(_ context.Context, in *proto.Request) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := models.Request{Login: in.GetLogin(), Pass: in.GetPass(), Ip: in.GetIp()}

	res := a.useCases.AllowRequest(ctx, request)

	return &proto.Response{
		Ok: &wrapperspb.BoolValue{Value: res},
	}, nil
}

func (a *AntibruteforceService) ClearBucket(_ context.Context, in *proto.Request) (*proto.Response, error) {
	request := models.Request{Login: in.GetLogin(), Ip: in.GetIp()}

	if err := a.useCases.ClearBucket(request); err != nil {
		return nil, err
	}
	return &proto.Response{
		Ok: &wrapperspb.BoolValue{Value: true},
	}, nil
}

func (a *AntibruteforceService) AddToBlackList(_ context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, subnet, err := net.ParseCIDR(in.GetSubnet())
	if err != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false}},
			status.Error(codes.InvalidArgument, err.Error())
	}

	message, err := a.useCases.Add(ctx, subnet.String(), "black")
	if err != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false}},
			status.Error(codes.Internal, err.Error())
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: message}, nil
}

func (a *AntibruteforceService) AddToWhiteList(_ context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, subnet, err := net.ParseCIDR(in.GetSubnet())
	if err != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false}},
			status.Error(codes.InvalidArgument, err.Error())
	}

	message, err := a.useCases.Add(ctx, subnet.String(), "white")
	if err != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false}},
			status.Error(codes.Internal, err.Error())
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: message}, nil
}

func (a *AntibruteforceService) RemoveFromBlackList(_ context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, subnet, err := net.ParseCIDR(in.GetSubnet())
	if err != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false}},
			status.Error(codes.InvalidArgument, err.Error())
	}

	if err := a.useCases.Remove(ctx, subnet.String(), "black"); err != nil {
		return nil, err
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: "OK"}, nil
}

func (a *AntibruteforceService) RemoveFromWhiteList(_ context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, subnet, err := net.ParseCIDR(in.GetSubnet())
	if err != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false}},
			status.Error(codes.InvalidArgument, err.Error())
	}

	if err := a.useCases.Remove(ctx, subnet.String(), "white"); err != nil {
		return nil, err
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: "OK"}, nil
}

// TODO: дописать метод поиска ip в списках
//func (a *AntibruteforceService) SearchInLists(ctx context.Context, in *proto.Subnet) (*proto.Response, error) {
//	ip := net.ParseIP(in.GetSubnet())
//	if ip == nil {
//		return &proto.Response{Ok: &wrappers.BoolValue{Value: false}}, status.Error(codes.InvalidArgument, "invalid IP address")
//	}
//
//	value, err := a.listUseCases.SearchIPInLists(ip)
//	if err != nil {
//		return &proto.Response{Ok: &wrappers.BoolValue{Value: false}}, status.Error(codes.Internal, err.Error())
//	}
//
//	return &proto.Response{Ok: &wrappers.BoolValue{Value: value}}, nil
//}
