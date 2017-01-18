package merchant

import (
	"fmt"

	"github.com/blackdog1987/annual-card/lib/data"
	"github.com/blackdog1987/annual-card/module"
)

type MerchantStoreService struct{}

func init() {
	module.MerchantStore = &MerchantStoreService{}
}

//
func (s *MerchantStoreService) Count(mchID int64, keyword string) (count int64, err error) {
	return data.Db.Where(s.generateParam(mchID, keyword)).Count(&data.MerchantStore{})
}

//
func (s *MerchantStoreService) Search(mchId int64, keyword string, offset, length int) (items []data.MerchantStore, err error) {
	items = make([]data.MerchantStore, 0)
	err = data.Db.Where(s.generateParam(mchId, keyword)).Limit(length, offset).OrderBy("created ASC").Find(&items)
	return
}

//
func (s *MerchantStoreService) Get(mid int64) (has bool, info data.MerchantStore) {
	info = data.MerchantStore{}
	has, _ = data.Db.Where("store_id = ?", mid).Get(&info)
	return
}

func (s *MerchantStoreService) GetByAccount(account string) (has bool, info data.MerchantStore) {
	info = data.MerchantStore{}
	has, _ = data.Db.Where(fmt.Sprintf("account='%s'", account)).Get(&info)
	return
}
//
func (s *MerchantStoreService) IsAccountRepeat(storeID int64, name string) (has bool) {
	var w = fmt.Sprintf("account = '%s'", name)
	if storeID > 0 {
		w += fmt.Sprintf(" AND store_id != %d", storeID)
	}
	total, _ := data.Db.Where(w).Count(new(data.MerchantStore))
	return total > 0
}

//
func (s *MerchantStoreService) generateParam(mchID int64, keyword string) (w string) {
	w = fmt.Sprintf("mch_id = %d", mchID)
	if len(keyword) > 0 {
		w += fmt.Sprintf(" AND store_name LIKE '%%%s%%'", keyword)
	}
	return
}
