package usermodels

import (
	"github.com/xmopen/authsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/database/xmuser"
)

const XMUserTableName = "t_xm_user"

// SaveUser create user.
func SaveUser(user *xmuser.XMUser) error {
	return config.AuthDataBase().Table(XMUserTableName).Create(user).Error
}

// XMUserWithAccount 根据Account获取XMUser.
func XMUserWithAccount(account string) (*xmuser.XMUser, error) {
	xmUser := xmuser.New()
	return xmUser, config.AuthDataBase().Table(XMUserTableName).Where("user_account = ?", account).
		First(xmUser).Error
}
