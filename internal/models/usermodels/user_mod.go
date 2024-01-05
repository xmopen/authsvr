package usermodels

import (
	"errors"

	"github.com/xmopen/authsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/database/xmuser"
	"gorm.io/gorm"
)

const XMUserTableName = "t_xm_user"

// SaveUser creational user.
func SaveUser(user *xmuser.XMUser) error {
	return config.AuthDataBase().Table(XMUserTableName).Create(user).Error
}

// XMUserWithAccount 根据Account获取XMUser.
func XMUserWithAccount(account string) (*xmuser.XMUser, error) {
	xmUser := xmuser.New()
	err := config.AuthDataBase().Table(XMUserTableName).Where("user_account = ?", account).
		First(xmUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return xmUser, nil
}
