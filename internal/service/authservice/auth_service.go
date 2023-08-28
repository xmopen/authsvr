// Package authservice 鉴权service.
package authservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/xmopen/authsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/database/xmuser"
)

const (
	XMUserWithTokenKey   = "xm_user_token_%s"
	XMUserWithAccountKey = "xm_user_account_%s"
)

var (
	authServiceInstance *AuthService
	authServiceOnce     sync.Once
)

// AuthService 鉴权服务.
type AuthService struct {
	xredis *redis.Client
}

// Service 获取AuthService instance.
func Service() *AuthService {
	if authServiceInstance == nil {
		authServiceOnce.Do(func() {
			authServiceInstance = &AuthService{
				xredis: config.GlobalRedis(),
			}
		})
	}
	return authServiceInstance
}

// XMUserWithToken 获取XMUser通过Token.
func (a *AuthService) XMUserWithToken(token string) (*xmuser.XMUser, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	return a.getXMUserFromRedisWithKey(fmt.Sprintf(XMUserWithTokenKey, token))
}

// XMUserWithAccount 获取XMUser通过Account.
func (a *AuthService) XMUserWithAccount(account string) (*xmuser.XMUser, error) {
	if account == "" {
		return nil, fmt.Errorf("account is empty")
	}
	return a.getXMUserFromRedisWithKey(fmt.Sprintf(XMUserWithAccountKey, account))
}

func (a *AuthService) getXMUserFromRedisWithKey(key string) (*xmuser.XMUser, error) {
	result, err := a.xredis.Get(context.TODO(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	xmUser := xmuser.New()
	if err := json.Unmarshal([]byte(result), xmUser); err != nil {
		return nil, err
	}
	return xmUser, nil
}
