// Package authservice 鉴权service.
package authservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/xmopen/golib/pkg/utils"

	"github.com/redis/go-redis/v9"
	"github.com/xmopen/authsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/database/xmuser"
)

const (
	XMUserWithTokenKey          = "xm_user_token_%s"
	XMUserWithAccountKey        = "xm_user_account_%s"
	XMUserAccountToTokenMapping = "xm_user_account_token_mapping_%s"

	// defaultXMUserAuthTokenExpire  默认Token过期时间
	defaultXMUserAuthTokenExpire = 7 * 24 * time.Hour
)

var (
	authServiceInstance *AuthService
	authServiceOnce     sync.Once
	xmTokenCacheExpire  time.Duration
)

func init() {
	xmTokenCacheExpireTs := config.Config().GetInt64("xmuser.auth.token.expire")
	if xmTokenCacheExpireTs <= 0 {
		xmTokenCacheExpire = defaultXMUserAuthTokenExpire
	} else {
		xmTokenCacheExpire = time.Duration(xmTokenCacheExpireTs)
	}
}

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

// CreateXMUserToken 创建XMUser Token.
func (a *AuthService) CreateXMUserToken(xmUser *xmuser.XMUser) (string, error) {

	xmToken, err := a.xredis.Get(context.TODO(), fmt.Sprintf(XMUserAccountToTokenMapping, xmUser.UserAccount)).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if xmToken != "" {
		return xmToken, nil
	}

	for {
		token := strings.ToUpper(utils.UUID())
		tempXMUser, err := a.getXMUserFromRedisWithKey(fmt.Sprintf(XMUserWithTokenKey, token))
		if err != nil {
			return "", err
		}
		if tempXMUser != nil {
			continue
		}
		xmUserData, err := json.Marshal(xmUser)
		if err != nil {
			return "", err
		}
		if _, err = a.xredis.Set(context.TODO(), fmt.Sprintf(XMUserWithTokenKey, token), xmUserData,
			time.Duration(xmTokenCacheExpire)).Result(); err != nil {
			return "", err
		}
		// 保存account->token的映射.
		if _, err = a.xredis.Set(context.TODO(), fmt.Sprintf(XMUserAccountToTokenMapping, xmUser.UserAccount), token,
			time.Duration(xmTokenCacheExpire)).Result(); err != nil {
			return "", err
		}
		return token, nil
	}
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
