// 接入层:账号处理,用于登录、生成密钥等

package duobb

import (
	"github.com/reechou/duobb_access/models"
)

func (self *DuobbProcess) GetDuobbAccount(userName string) (*models.DuobbAccount, error) {
	account := &models.DuobbAccount{
		UserName: userName,
	}
	err := models.GetDuobbAccount(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
