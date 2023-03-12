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

	if request.Login == "" || request.Pass == "" || request.IP == "" {
		return &proto.Response{
			Ok:      &wrapperspb.BoolValue{Value: false},
			Message: "login/password/IP should not be empty",
		}, nil
	}

	if res := net.ParseIP(request.IP); res == nil {
		return &proto.Response{
				Ok: &wrappers.BoolValue{Value: false},
			},
			status.Error(codes.InvalidArgument, "invalid IP")
	}

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
		return &proto.Response{},
			status.Error(codes.InvalidArgument, er.Error())
	}

	ok, message, err := a.useCases.Add(ctx, subnet, blacklist)
	if err != nil {
		a.logg.Error(err.Error())
		return &proto.Response{},
			status.Error(codes.Internal, "Internal problems")
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: ok},
		Message: message,
	}, nil
}

func (a *AntibruteforceService) AddToWhiteList(ctx context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subnet, er := a.getSubnet(in.GetSubnet())
	if er != nil {
		return &proto.Response{},
			status.Error(codes.InvalidArgument, er.Error())
	}

	ok, message, err := a.useCases.Add(ctx, subnet, whitelist)
	if err != nil {
		a.logg.Error(err.Error())
		return &proto.Response{},
			status.Error(codes.Internal, "Internal problems")
	}
	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: ok},
		Message: message,
	}, nil
}

func (a *AntibruteforceService) RemoveFromBlackList(ctx context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subnet, er := a.getSubnet(in.GetSubnet())
	if er != nil {
		return &proto.Response{},
			status.Error(codes.InvalidArgument, "invalid subnet")
	}

	if subnet == "" {
		return &proto.Response{}, status.Error(codes.InvalidArgument, "subnet should not be empty")
	}

	mess, err := a.useCases.Remove(ctx, subnet, blacklist)
	if err != nil {
		a.logg.Error(err.Error())
		return &proto.Response{},
			status.Error(codes.Internal, "Internal problems")
	}
	if mess != "" {
		return &proto.Response{
			Ok:      &wrappers.BoolValue{Value: false},
			Message: mess,
		}, nil
	}

	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: "",
	}, nil
}

func (a *AntibruteforceService) RemoveFromWhiteList(ctx context.Context, in *proto.Subnet) (*proto.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subnet, er := a.getSubnet(in.GetSubnet())
	if er != nil {
		return &proto.Response{}, status.Error(codes.InvalidArgument, "invalid subnet")
	}

	if subnet == "" {
		return &proto.Response{}, status.Error(codes.InvalidArgument, "subnet should not be empty")
	}

	mess, err := a.useCases.Remove(ctx, subnet, whitelist)
	if err != nil {
		a.logg.Error(err.Error())
		return &proto.Response{},
			status.Error(codes.Internal, "Internal problems")
	}
	if mess != "" {
		return &proto.Response{
			Ok:      &wrappers.BoolValue{Value: false},
			Message: mess,
		}, nil
	}

	return &proto.Response{
		Ok:      &wrappers.BoolValue{Value: true},
		Message: "",
	}, nil
}

func (a *AntibruteforceService) getSubnet(in string) (string, error) {
	_, subnet, err := net.ParseCIDR(in)
	return subnet.String(), err
}
