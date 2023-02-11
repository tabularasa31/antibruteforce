package usecase

import (
	"github.com/tabularasa31/antibruteforce/internal/controller/repo"
	"github.com/tabularasa31/antibruteforce/internal/models"
)

// BucketUseCase -.
type BucketUseCase struct {
	repo repo.BucketRepo
}

// New -.
func New(r repo.BucketRepo) *BucketUseCase {
	return &BucketUseCase{
		repo: r,
	}
}

func (u *BucketUseCase) Check(q *models.Request) bool {
	return false
}

//
//rpc Query(Request) returns (Response);
//rpc ClearBucket(Request) returns(google.protobuf.Empty);
//rpc AddToBlackList(Subnet) returns(Response);
//rpc AddToWhiteList(Subnet) returns(Response);
//rpc RemoveFromBlackList(Subnet) returns(Response);
//rpc RemoveFromWhiteList(Subnet) returns(Response);
