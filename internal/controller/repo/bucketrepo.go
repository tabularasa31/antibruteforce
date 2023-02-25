package repo

import (
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"github.com/tabularasa31/antibruteforce/config"
	"github.com/tabularasa31/antibruteforce/internal/models"
	"time"
)

// BucketRepo -.
type BucketRepo struct {
	Redis           *redis.Client
	cfg             *config.AppConfig
	ipLimiter       *redis_rate.Limiter
	loginLimiter    *redis_rate.Limiter
	passwordLimiter *redis_rate.Limiter
}

// NewBucketRepo -.
func NewBucketRepo(r *redis.Client, cfg *config.AppConfig) *BucketRepo {
	return &BucketRepo{
		Redis:           r,
		cfg:             cfg,
		ipLimiter:       redis_rate.NewLimiter(r),
		loginLimiter:    redis_rate.NewLimiter(r),
		passwordLimiter: redis_rate.NewLimiter(r),
	}
}

func (b *BucketRepo) CheckLimit(request models.Request) bool {
	ip := request.Ip
	login := request.Login
	password := request.Pass

	// Check if the rate limits have been exceeded for each key
	if request.Ip != "" {
		_, _, ipOK := b.ipLimiter.AllowMinute(ip, int64(b.cfg.IpLimit))
		if !ipOK {
			return false
		}
	}

	if request.Login != "" {
		_, _, loginOK := b.loginLimiter.AllowMinute(login, int64(b.cfg.LoginLimit))
		if !loginOK {
			return false
		}
	}
	if request.Pass != "" {
		_, _, passwordOK := b.passwordLimiter.AllowMinute(password, int64(b.cfg.PassLimit))
		if !passwordOK {
			return false
		}
	}
	return true
}

func (b *BucketRepo) ClearBucket(request models.Request) error {
	if errIp := b.ipLimiter.Reset(request.Ip, 1*time.Minute); errIp != nil {
		return errIp
	}
	if errLogin := b.loginLimiter.Reset(request.Login, 1*time.Minute); errLogin != nil {
		return errLogin
	}
	return nil
}
