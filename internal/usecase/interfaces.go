package usecase

import "github.com/tabularasa31/antibruteforce/internal/models"

type (
	BucketRepo interface {
		AllowRequest(request models.Request) bool
		ClearBucket(request models.Request)
	}
	ListRepo interface {
		Add(subnet string, list string) bool
		Remove(subnet string, list string) bool
	}
)

//rpc AllowRequest(Request) returns (Response);
//rpc ClearBucket(Request) returns(google.protobuf.Empty);
//rpc AddToBlackList(Subnet) returns(Response);
//rpc AddToWhiteList(Subnet) returns(Response);
//rpc RemoveFromBlackList(Subnet) returns(Response);
//rpc RemoveFromWhiteList(Subnet) returns(Response);
