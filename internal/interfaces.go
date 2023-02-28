package internal

import "github.com/tabularasa31/antibruteforce/internal/models"

type UseCases interface {
	AllowRequest(request models.Request) bool
	ClearBucket(request models.Request)
	Add(subnet string, list string)
	Remove(subnet string, list string)
	CheckSubnetInList(subnet string) error
}
