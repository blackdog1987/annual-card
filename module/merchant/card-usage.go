package merchant

import (
	"fmt"

	"github.com/blackdog1987/annual-card/lib/data"
	"github.com/blackdog1987/annual-card/module"
)

type CardUsageService struct{}

func init() {
	module.CardUsage = &CardUsageService{}
}

//
func (s *CardUsageService) Count(mchID, storeID, start, end int64, isStore bool) (count int64, err error) {
	return data.Db.Where(s.generateParam(mchID, storeID, start, end, isStore)).Count(&data.AnnualCardUsageLog{})
}

//
func (s *CardUsageService) Search(mchID, storeID, start, end int64, isStore bool, offset, length int) (items []data.AnnualCardUsageLog, err error) {
	items = make([]data.AnnualCardUsageLog, 0)
	err = data.Db.Table("annual_card_usage_log").
		Join("LEFT", "annual_card", "annual_card.card_id = annual_card_usage_log.card_id").
		Join("LEFT", "merchant_store", "annual_card_usage_log.store_id=merchant_store.store_id").
		Where(s.generateParam(mchID, storeID, start, end, isStore)).
		Limit(length, offset).
		OrderBy("usage_time DESC").
		Find(&items)
	return
}

//
func (s *CardUsageService) generateParam(mchID, storeID, start, end int64, isStore bool) (w string) {
	w = "1=1"
	if isStore {
		if storeID > 0 {
			w += fmt.Sprintf(" AND annual_card_usage_log.store_id = %d", storeID)
		}
	} else {
		if mchID > 0 {
			w += fmt.Sprintf(" AND annual_card_usage_log.mch_id = %d", mchID)
		}
	}
	if start > 0 {
		w += fmt.Sprintf(" AND annual_card_usage_log.usage_time >= %d", start)
	}
	if end > 0 {
		w += fmt.Sprintf(" AND annual_card_usage_log.usage_time <= %d", end)
	}
	return
}
