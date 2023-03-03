package repo

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/models"
)

var (
	cfg = config.AppConfig{LoginLimit: 10, PassLimit: 100, IPLimit: 1000}

	request = models.Request{
		Login: "test_login",
		Pass:  "test_pass",
		IP:    "192.168.0.7",
	}
)

func Test_CheckLimit(t *testing.T) {
	tests := []struct {
		sleep   time.Duration
		name    string
		wantRes bool
	}{
		{
			sleep:   7 * time.Second,
			name:    "Allowed request",
			wantRes: true,
		}, {
			sleep:   1 * time.Millisecond,
			name:    "Not allowed request",
			wantRes: false,
		},
	}
	s := miniredis.RunT(t)
	c := redis.NewClient(&redis.Options{Addr: s.Addr()})

	bucketRepo := NewBucketRepo(c, &cfg)

	for _, tt := range tests {
		_ = bucketRepo.ClearBucket(request)
		for i := 0; i < cfg.LoginLimit; i++ {
			bucketRepo.CheckLimit(request)
			time.Sleep(tt.sleep)
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := bucketRepo.CheckLimit(request); err != tt.wantRes {
				t.Errorf("AllowRequest = %v, wantRes %v", err, tt.wantRes)
			}
		})
	}
}

func Test_ClearBucket(t *testing.T) {
	s := miniredis.RunT(t)
	c := redis.NewClient(&redis.Options{Addr: s.Addr()})

	bucketRepo := NewBucketRepo(c, &cfg)
	_ = bucketRepo.ClearBucket(request)

	for i := 0; i < cfg.LoginLimit+2; i++ {
		bucketRepo.CheckLimit(request)
	}
	t.Run("Bucket cleared successfully", func(t *testing.T) {
		e := bucketRepo.ClearBucket(request)
		require.NoError(t, e)
		if err := bucketRepo.CheckLimit(request); err != true {
			t.Errorf("AllowRequest = %v, wantRes %v", err, true)
		}
	})
}
