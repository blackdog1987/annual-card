package merchant

import (
	"fmt"

	"github.com/blackdog1987/annual-card/lib/data"
	"github.com/blackdog1987/annual-card/module"
)

type MerchantService struct{}

func init() {
	module.Merchant = &MerchantService{}
}

//
func (s *MerchantService) Count(keyword string) (count int64, err error) {
	return data.Db.Where(s.generateParam(keyword)).Count(&data.Merchant{})
}

//
func (s *MerchantService) Search(keyword string, offset, length int) (items []data.Merchant, err error) {
	items = make([]data.Merchant, 0)
	err = data.Db.Where(s.generateParam(keyword)).Limit(length, offset).OrderBy("created ASC").Find(&items)
	return
}

//
func (s *MerchantService) Get(mid int64) (has bool, info data.Merchant) {
	info = data.Merchant{}
	has, _ = data.Db.Where("mch_id = ?", mid).Get(&info)
	return
}

//
func (s *MerchantService) IsMchNameRepeat(mchID int64, name string) (has bool) {
	var w = fmt.Sprintf("account = '%s'", name)
	if mchID > 0 {
		w += fmt.Sprintf("mch_id != %d", mchID)
	}
	total, _ := data.Db.Where(w).Count(new(data.Merchant))
	return total > 0
}

//
func (s *MerchantService) generateParam(keyword string) (w string) {
	w = "1=1"
	if len(keyword) > 0 {
		w += fmt.Sprintf(" AND name LIKE '%%%s%%'", keyword)
	}
	return
}
