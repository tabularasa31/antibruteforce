package models

import "time"

type (
	TokenBucket struct {
		Cap   int
		Token int
		Reset time.Time
	}
	Request struct {
		Login string
		Pass  string
		Ip    string
	}
)
