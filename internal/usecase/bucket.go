package usecase

import (
	"context"
	"github.com/tabularasa31/antibruteforce/internal/models"
	"net"
)

func (u *UseCases) AllowRequest(ctx context.Context, request models.Request) bool {
	switch u.lists.SearchIPInList(ctx, net.IP(request.Ip)) {
	case "white":
		return true
	case "black":
		return false
	default:
		return u.buckets.Allow(request)
	}
}

func (u *UseCases) ClearBucket(request models.Request) error {
	return u.buckets.ClearBucket(request)
}
