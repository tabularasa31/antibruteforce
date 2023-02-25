package usecase

import (
	"context"
)

func (u *UseCases) Add(ctx context.Context, subnet string, color string) (string, error) {
	return u.lists.SaveToList(ctx, subnet, color)
}

func (u *UseCases) Remove(ctx context.Context, subnet string, color string) error {
	if err := u.lists.DeleteFromList(ctx, subnet, color); err != nil {
		return err
	}
	return nil
}

func (u *UseCases) CheckSubnetColor(ctx context.Context, subnet string) (string, error) {
	color, err := u.lists.CheckColor(ctx, subnet)
	if err != nil {
		return "", err
	}
	return color, nil
}

//func (lu *ListUseCase) SearchIPInLists(ctx context.Context, ip net.IP) bool {
//	if found := lu.repo.SearchIP(ctx, ip); found != "" {
//		return false
//	}
//	return true
//}