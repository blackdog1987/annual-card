package order

import (
	"fmt"

	"github.com/blackdog1987/annual-card/lib/data"
	"github.com/blackdog1987/annual-card/module"
)

type OrderService struct{}

func init() {
	module.Order = &OrderService{}
}

//
func (s *OrderService) Count(start, end int64) (count int64, err error) {
	return data.Db.Where(s.generateParam(start, end)).Count(&data.OrderLog{})
}

//
func (s *OrderService) Search(start, end int64, offset, length int) (items []data.OrderLog, err error) {
	items = make([]data.OrderLog, 0)
	err = data.Db.Where(s.generateParam(start, end)).Limit(length, offset).OrderBy("created DESC").Find(&items)
	return
}

func (s *OrderService) Coupon(uid int64) (info *data.CouponLog, has bool) {
	info = &data.CouponLog{}
	has, _ = data.Db.Where("uid = ? AND is_usage = ?", uid, 0).Get(info)
	return
}

//
func (s *OrderService) generateParam(start, end int64) (w string) {
	w = "is_pay=1"
	if start > 0 {
		w += fmt.Sprintf(" AND updated >= %d", start)
	}
	if end > 0 {
		w += fmt.Sprintf(" AND updated <= %d", end)
	}
	return
}
