package usecase

import (
	"github.com/tabularasa31/antibruteforce/internal/controller/repo"
)

type UseCases struct {
	buckets repo.BucketRepo
	lists   repo.ListRepo
}

// New -.
func New(br *repo.BucketRepo, l *repo.ListRepo) *UseCases {
	return &UseCases{
		buckets: *br,
		lists:   *l,
	}
}
