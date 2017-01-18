package annual

import (
	"fmt"

	"github.com/molibei/annual-card/lib/data"
	"github.com/molibei/annual-card/module"
)

type AnnualCardService struct{}

func init() {
	module.AnnualCard = &AnnualCardService{}
}

func (s *AnnualCardService) Count(isActive int, planID, start, end int64, name, idcard string) (count int64, err error) {
	return data.Db.Where(s.generateParam(isActive, planID, start, end, name, idcard)).Count(&data.AnnualCard{})
}
func (s *AnnualCardService) Search(isActive int, planID, start, end int64,name,idcard string, offset, length int) (items []data.AnnualCard, err error) {
	items = make([]data.AnnualCard, 0)
	err = data.Db.Where(s.generateParam(isActive, planID, start, end, name, idcard)).Limit(length, offset).OrderBy("created DESC").Find(&items)
	return
}

// 绑定统计
func (s *AnnualCardService) CountBind(uid int64, active int) (total int64, err error) {
	return data.Db.Where("relation_uid = ? AND is_active = ?", uid, active).Count(&data.AnnualCard{})
}
// 绑定查询
func (s *AnnualCardService) SearchBind(uid int64, active, offset, length int) (items []data.AnnualCard, err error) {
	items = make([]data.AnnualCard, 0)
	stmt := data.Db.Where("relation_uid = ? and is_delete=0", uid)
	if active != -1 {
		stmt.And("is_active = ?", active)
	}
	err = stmt.Limit(length, offset).OrderBy("is_active DESC, updated DESC").Find(&items)
	return
}

func (s *AnnualCardService) Usages(uid, mchID int64) (items []data.CardUsageInfo, err error) {
	items = make([]data.CardUsageInfo, 0)
	sqlStr := `SELECT card_id,
	   card_name,
	   card_no,
	   bind_name,
       bind_idcard,
       bind_contact,
       bind_headimg,
       (SELECT count(*)
  FROM annual_card_usage_log
 WHERE annual_card_usage_log.card_id= annual_card.card_id
   AND annual_card_usage_log.mch_id= ?) as usage_num
  FROM annual_card
 WHERE annual_card.relation_uid= ?
   AND annual_card.is_active= ?
   AND annual_card.is_delete=0`
	err = data.Db.Sql(sqlStr, mchID, uid, 1).Find(&items)
	return
}

func (s *AnnualCardService) UsagesByCardID(cardID int64) (info data.CardUsageInfo, has bool) {
	info = data.CardUsageInfo{}
	sqlStr := `SELECT card_id,
	   card_name,
	   card_no,
	   bind_name,
       bind_idcard,
       bind_contact,
       bind_headimg,
       (
SELECT count(*)
  FROM annual_card_usage_log
 WHERE annual_card_usage_log.card_id= annual_card.card_id
   ) as usage_num
  FROM annual_card
 WHERE annual_card.card_id= ?`
	has, _ = data.Db.Sql(sqlStr, cardID).Get(&info)
	return
}

//
func (s *AnnualCardService) PlanCount() (count int64, err error) {
	return data.Db.Count(&data.AnnualCardPlan{})
}

//
func (s *AnnualCardService) PlanSearch(offset, length int) (items []data.AnnualCardPlan, err error) {
	items = make([]data.AnnualCardPlan, 0)
	err = data.Db.Sql(`SELECT
  cp_id,
  channel,
  expired_start,
  expired_stop,
  create_num,
  created,
  updated,
  is_disabled,
  (SELECT count(*)
   FROM annual_card
   WHERE annual_card.plan_id=annual_card_plan.cp_id
      AND annual_card.is_active = 1 AND annual_card.is_delete = 0) as active_num
FROM annual_card_plan
ORDER BY created DESC
LIMIT ?, ?`, offset, length).Find(&items)
	return
}

func (s *AnnualCardService) CardByCardNo(cardNo string) (info *data.AnnualCard, has bool) {
	info = &data.AnnualCard{}
	has, _ = data.Db.Where("card_no = ?", cardNo).Get(info)
	return
}

func (s *AnnualCardService) CardByID(cardID int64) (info *data.AnnualCard, has bool) {
	info = &data.AnnualCard{}
	has, _ = data.Db.Where("card_id = ?", cardID).Get(info)
	return
}

func (s *AnnualCardService) IsExistIDCard(cardID int64, idCard string) (info *data.AnnualCard, has bool) {
	info = &data.AnnualCard{}
	has, _ = data.Db.Where("card_id != ? AND bind_idcard = ?", cardID, idCard).Get(info)
	return
}

//
func (s *AnnualCardService) generateParam(isActive int, planID, start, end int64, name, idcard string) (w string) {
	w = "is_delete=0"
	if isActive != -1 {
		w += fmt.Sprintf(" AND is_active = %d", isActive)
	}
	if start > 0 {
		w += fmt.Sprintf(" AND updated >= %d", start)
	}
	if end > 0 {
		w += fmt.Sprintf(" AND updated <= %d", end)
	}
	if planID != 0 {
		w += fmt.Sprintf(" AND plan_id = %d", planID)
	}
	if len(name) > 0 {
		w += fmt.Sprintf(" AND bind_name LIKE '%%%s%%'", name)
	}
	if len(idcard) > 0 {
		w += fmt.Sprintf(" AND bind_idcard LIKE '%%%s%%'", idcard)
	}
	return
}
