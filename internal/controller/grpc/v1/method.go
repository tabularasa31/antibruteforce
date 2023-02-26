//nolint:dupl
package grpcv1

import (
	"context"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	proto "github.com/tabularasa31/antibruteforce/api"
	"github.com/tabularasa31/antibruteforce/internal/models"
	"github.com/tabularasa31/antibruteforce/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	whitelist = "white"
	blacklist = "black"
	ok        = "OK"
)

type AntibruteforceService struct {
	useCases *usecase.UseCases
	logg     *zap.Logger
	proto.UnimplementedAntiBruteforceServer
}

func NewAntibruteforceService(u *usecase.UseCases, logg *zap.Logger) *AntibruteforceService {
	return &AntibruteforceService{useCases: u, logg: logg}
}

func (a *AntibruteforceService) AllowRequest(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	request := models.Request{Login: in.GetLogin(), Pass: in.GetPass(), IP: in.GetIp()}

	res := a.useCases.AllowRequest(ctx, request)

	return &proto.Response{
		Ok: &wrapperspb.BoolValue{Value: res},
	}, nil
}

func (a *AntibruteforceService) ClearBucket(_ context.Context, in *proto.Request) (*proto.Response, error) {
	request := models.Request{Login: in.GetLogin(), IP: in.GetIp()}

	if err := a.useCases.ClearBucket(request); err != nil {
		return nil, err
	}
	return &proto.Response{
		Ok: &wrapperspb.BoolValue{Value: true},
	}, nil
}

func (a *AntibruteforceService) AddToBlackList(ctx context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subnet, er := a.getSubnet(in.GetSubnet())
	if er != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false},
			},
			status.Error(codes.InvalidArgument, er.Error())
	}

	message, err := a.useCases.Add(ctx, subnet, blacklist)
	if err != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false},
			},
			status.Error(codes.Internal, err.Error())
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: message,
	}, nil
}

func (a *AntibruteforceService) AddToWhiteList(ctx context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subnet, er := a.getSubnet(in.GetSubnet())
	if er != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false},
			},
			status.Error(codes.InvalidArgument, er.Error())
	}

	message, err := a.useCases.Add(ctx, subnet, whitelist)
	if err != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false},
			},
			status.Error(codes.Internal, err.Error())
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: message,
	}, nil
}

func (a *AntibruteforceService) RemoveFromBlackList(ctx context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subnet, er := a.getSubnet(in.GetSubnet())
	if er != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false},
			},
			status.Error(codes.InvalidArgument, er.Error())
	}

	if err := a.useCases.Remove(ctx, subnet, blacklist); err != nil {
		return nil, err
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: ok,
	}, nil
}

func (a *AntibruteforceService) RemoveFromWhiteList(ctx context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subnet, er := a.getSubnet(in.GetSubnet())
	if er != nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false},
			},
			status.Error(codes.InvalidArgument, er.Error())
	}

	if err := a.useCases.Remove(ctx, subnet, whitelist); err != nil {
		return nil, err
	}

	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: ok,
	}, nil
}

func (a AntibruteforceService) getSubnet(in string) (string, error) {
	_, subnet, err := net.ParseCIDR(in)

	s := subnet.String()
	return s, err
}
