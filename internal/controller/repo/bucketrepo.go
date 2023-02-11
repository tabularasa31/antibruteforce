package repo

import (
	"github.com/tabularasa31/antibruteforce/internal/models"
	"github.com/tabularasa31/antibruteforce/pkg/redis"
)

// BucketRepo -.
type BucketRepo struct {
	db *redis.Redis
}

// NewBucketRepo -.
func NewBucketRepo(r *redis.Redis) *BucketRepo {
	return &BucketRepo{
		db: r,
	}
}

func (b *BucketRepo) AllowRequest(request models.Request) bool {
	return true
}
func (b *BucketRepo) ClearBucket(request models.Request) {

}

func (b *BucketRepo) FetchToken() {

}
func (b *BucketRepo) UpdateToken() {

}
