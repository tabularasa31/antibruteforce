package usecase

import (
	"context"
	"net"

	"github.com/tabularasa31/antibruteforce/internal/models"
)

func (u *UseCases) AllowRequest(ctx context.Context, request models.Request) bool {
	switch u.lists.SearchIPInList(ctx, net.IP(request.IP)) {
	case "white":
		return true
	case "black":
		return false
	default:
		return u.buckets.CheckLimit(request)
	}
}

func (u *UseCases) ClearBucket(request models.Request) error {
	return u.buckets.ClearBucket(request)
}
