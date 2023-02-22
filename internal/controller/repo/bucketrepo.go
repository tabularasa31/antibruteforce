package repo

import (
	"fmt"
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

func (b *BucketRepo) Allow(request models.Request) bool {
	ip := request.Ip
	login := request.Login
	password := request.Pass

	// Check if the rate limits have been exceeded for each key
	if request.Ip != "" {
		ipC, ipD, ipOK := b.ipLimiter.AllowMinute(ip, int64(b.cfg.IpLimit))
		fmt.Printf("ip count - %v | ip delay - %v | ipOK - %v\n", ipC, ipD, ipOK)
		if !ipOK {
			return false
		}
	}

	if request.Login != "" {
		loginC, loginD, loginOK := b.loginLimiter.AllowMinute(login, int64(b.cfg.LoginLimit))
		fmt.Printf("login count - %v | login delay - %v | loginOK - %v\n", loginC, loginD, loginOK)
		if !loginOK {
			return false
		}
	}
	if request.Pass != "" {
		passC, passD, passwordOK := b.passwordLimiter.AllowMinute(password, int64(b.cfg.PassLimit))
		fmt.Printf("password count - %v | password delay - %v | passwordOK - %v\n", passC, passD, passwordOK)
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
