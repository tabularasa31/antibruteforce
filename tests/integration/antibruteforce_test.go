//go:build integration

package integration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	proto "github.com/tabularasa31/antibruteforce/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AbfSuite struct {
	suite.Suite
	ctx    context.Context
	conn   *grpc.ClientConn
	client proto.AntiBruteforceClient
}

func (s *AbfSuite) SetupSuite() {
	grpcHost := os.Getenv("GRPC_HOST")
	if grpcHost == "" {
		grpcHost = "localhost:50051"
	}
	var err error
	s.conn, err = grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.client = proto.NewAntiBruteforceClient(s.conn)
}

type testCase struct {
	description string
	request     *proto.Request
	expectedOk  bool
	expectedErr error
}

func (s *AbfSuite) Test_AllowRequest() {
	testCases := []testCase{
		{
			description: "success case",
			request: &proto.Request{
				Login: "test",
				Pass:  "secret",
				Ip:    "192.168.0.9",
			},
			expectedOk:  true,
			expectedErr: nil,
		},
		{
			description: "empty login",
			request: &proto.Request{
				Login: "",
				Pass:  "secret",
				Ip:    "192.168.0.9",
			},
			expectedOk:  false,
			expectedErr: nil,
		},
		{
			description: "empty password",
			request: &proto.Request{
				Login: "test",
				Pass:  "",
				Ip:    "192.168.0.9",
			},
			expectedOk:  false,
			expectedErr: nil,
		},
		{
			description: "empty IP",
			request: &proto.Request{
				Login: "test",
				Pass:  "secret",
				Ip:    "",
			},
			expectedOk:  false,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n Test Case: %s \n", tc.description)
		resp, err := s.client.AllowRequest(s.ctx, tc.request)
		s.Require().Equal(tc.expectedOk, resp.GetOk().Value)
		s.Require().Equal(tc.expectedErr, err)
	}

	tc := testCase{
		description: "invalid IP",
		request: &proto.Request{
			Login: "test",
			Pass:  "secret",
			Ip:    "200.0.1.2.",
		},
		expectedErr: status.Error(codes.InvalidArgument, "invalid IP"),
	}
	fmt.Printf("\n Test Case: %s \n", tc.description)
	_, err := s.client.AllowRequest(s.ctx, tc.request)
	s.Require().Equal(tc.expectedErr, err)
}

func (s *AbfSuite) Test_ClearBucket() {
	tc := testCase{
		description: "success clear bucket",
		request: &proto.Request{
			Login: "test",
			Pass:  "secret",
			Ip:    "192.168.1.2",
		},
		expectedOk: true,
	}
	fmt.Printf("\n Test Case: %s \n", tc.description)
	_, err := s.client.AllowRequest(s.ctx, tc.request)
	s.Require().NoError(err)
	resp, err := s.client.ClearBucket(s.ctx, tc.request)
	s.Require().NoError(err)
	s.Require().Equal(resp.GetOk().Value, tc.expectedOk)
}

type testSubnet struct {
	description  string
	request      *proto.Subnet
	expectedOk   bool
	expectedMess string
	expectedErr  error
}

func (s *AbfSuite) Test_WhiteList() {
	ts := testSubnet{
		description: "success add to whitelist",
		request: &proto.Subnet{
			Subnet: "100.0.0.0/25",
		},
		expectedOk:   true,
		expectedMess: "",
	}

	resp, err := s.client.RemoveFromWhiteList(s.ctx, ts.request)
	s.Require().NoError(err)

	fmt.Printf("\n Test Case: %s \n", ts.description)
	resp, err = s.client.AddToWhiteList(s.ctx, ts.request)
	s.Require().NoError(err)
	s.Require().Equal(ts.expectedOk, resp.GetOk().Value)
	s.Require().Equal(ts.expectedMess, resp.GetMessage())

	fmt.Printf("\n Test Case: %s \n", ts.description)
	resp, err = s.client.RemoveFromWhiteList(s.ctx, ts.request)
	s.Require().NoError(err)
	s.Require().Equal(ts.expectedOk, resp.GetOk().Value)
	s.Require().Equal(ts.expectedMess, resp.GetMessage())
}

func (s *AbfSuite) Test_BlackList() {
	ts := testSubnet{
		description: "success add to whitelist",
		request: &proto.Subnet{
			Subnet: "100.0.1.0/25",
		},
		expectedOk:   true,
		expectedMess: "",
	}

	resp, err := s.client.RemoveFromBlackList(s.ctx, ts.request)
	s.Require().NoError(err)

	fmt.Printf("\n Test Case: %s \n", ts.description)
	resp, err = s.client.AddToBlackList(s.ctx, ts.request)
	s.Require().NoError(err)
	s.Require().Equal(ts.expectedOk, resp.GetOk().Value)
	s.Require().Equal(ts.expectedMess, resp.GetMessage())

	fmt.Printf("\n Test Case: %s \n", ts.description)
	resp, err = s.client.RemoveFromBlackList(s.ctx, ts.request)
	s.Require().NoError(err)
	s.Require().Equal(ts.expectedOk, resp.GetOk().Value)
	s.Require().Equal(ts.expectedMess, resp.GetMessage())
}

func TestAbfSuite(t *testing.T) {
	suite.Run(t, new(AbfSuite))
}
