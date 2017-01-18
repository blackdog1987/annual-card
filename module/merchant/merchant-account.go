package merchant

import (
	"fmt"

	"github.com/molibei/annual-card/lib/data"
	"github.com/molibei/annual-card/module"
)

type MerchantAccountService struct{}

func init() {
	module.MerchantAccount = &MerchantAccountService{}
}

//
func (s *MerchantAccountService) Count(keyword string) (count int64, err error) {
	return data.Db.Where(s.generateParam(keyword)).Count(&data.MerchantAccount{})
}

//
func (s *MerchantAccountService) Search(keyword string, offset, length int) (items []data.MerchantAccount, err error) {
	items = make([]data.MerchantAccount, 0)
	err = data.Db.Where(s.generateParam(keyword)).Limit(length, offset).OrderBy("created ASC").Find(&items)
	return
}

//
func (s *MerchantAccountService) Get(mid int64) (has bool, info data.MerchantAccount) {
	info = data.MerchantAccount{}
	has, _ = data.Db.Where("mch_id = ?", mid).Get(&info)
	return
}

func (s *MerchantAccountService)GetByAccount(account string) (has bool, info data.MerchantAccount) {
	info = data.MerchantAccount{}
	has, _ = data.Db.Where(fmt.Sprintf("account='%s'", account)).Get(&info)
	return
}

//
func (s *MerchantAccountService) IsAccountRepeat(mchID int64, name string) (has bool) {
	var w = fmt.Sprintf("account = '%s'", name)
	if mchID > 0 {
		w += fmt.Sprintf(" AND mch_id != %d", mchID)
	}
	total, _ := data.Db.Where(w).Count(new(data.MerchantAccount))
	return total > 0
}

func (s *MerchantAccountService) IsMchAccountRepeat(mchID int64, account string) (has bool) {
	var w = fmt.Sprintf("account = '%s'", account)
	if mchID > 0 {
		w += fmt.Sprintf("mch_id != %d", mchID)
	}
	total, _ := data.Db.Where(w).Count(new(data.MerchantAccount))
	return total > 0
}


//
func (s *MerchantAccountService) generateParam(keyword string) (w string) {
	w = "1=1"
	if len(keyword) > 0 {
		w += fmt.Sprintf(" AND name LIKE '%%%s%%'", keyword)
	}
	return
}
