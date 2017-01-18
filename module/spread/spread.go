package spread

import (
	"fmt"

	"github.com/molibei/annual-card/lib/data"
	"github.com/molibei/annual-card/module"
	"net/url"
	"io/ioutil"
	"encoding/base64"
	"net/http"
	"bytes"
	"encoding/json"
	"errors"
)

// SpreadService .
type SpreadService struct{}

func init() {
	module.Spread = &SpreadService{}
}

// Plan .
func (s *SpreadService) Plan(planID int64) (plan *data.SpreadPlan, has bool) {
	plan = &data.SpreadPlan{}
	has, _ = data.Db.Where("sp_id = ?", planID).Get(plan)
	return
}

func (s *SpreadService) UpdatePlan(plan *data.SpreadPlan) (ok bool) {
	_, err := data.Db.Where("sp_id = ?", plan.SPID).
		Cols("name", "channel", "contact", "is_disabled", "reg_commission", "sale_commission", "qrcode").
		Update(plan)
	return err == nil
}

// PlansCount .
// 推广计划统计
func (s *SpreadService) PlansCount(state int) (count int64, err error) {
	var w = s.generatePlanParam(state)
	return data.Db.Where(w).Count(&data.SpreadPlan{})
}

// Plans .
// 推广计划列表
func (s *SpreadService) Plans(state int, offset, length int) (items []data.SpreadPlan, err error) {
	var w = s.generatePlanParam(state)
	items = make([]data.SpreadPlan, 0)
	err = data.Db.Where(w).Limit(length, offset).OrderBy("created DESC").Find(&items)
	return items, err
}

// IsPlanNameRepeat .
// 推广计划名称重复检测
func (s *SpreadService) IsPlanNameRepeat(id int64, name string) (has bool) {
	var w = fmt.Sprintf("name = '%s'", name)
	if id > 0 {
		w += fmt.Sprintf("sp_id != %d", id)
	}
	total, _ := data.Db.Where(w).Count(new(data.SpreadPlan))
	return total > 0
}

// LogsCount .
// 推广明细统计
func (s *SpreadService) LogsCount(planID, start, end int64) (count int64, err error) {
	var w = fmt.Sprintf("plan_id=%d", planID)
	if start > 0 {
		w += fmt.Sprintf(" AND created>=%d", start)
	}
	if end > 0 {
		w += fmt.Sprintf(" AND created<=%d", end)
	}
	return data.Db.Where(w).Count(&data.SpreadLogs{})
}

// Logs .
func (s *SpreadService) Logs(planID, start, end int64, offset, length int) (items []data.SpreadLogItem, err error) {
	items = make([]data.SpreadLogItem, 0)
	var w = fmt.Sprintf(" `spread_logs`.`plan_id`=%d ", planID)
	if start > 0 {
		w += fmt.Sprintf("AND `spread_logs`.`created`>=%d ", start)
	}
	if end > 0 {
		w += fmt.Sprintf("AND `spread_logs`.`created`<=%d ", end)
	}
	err = data.Db.Sql("SELECT `spread_logs`.`log_id`, `spread_logs`.`plan_id`,`spread_logs`.`uid`," +
		"`spread_logs`.`category`, `spread_logs`.`commission`, `spread_logs`.`order_total`," +
		"`spread_logs`.`created`,`consumer`.`nickname`, `consumer`.`realname`,`consumer`.`phone`" +
		"FROM `spread_logs` " +
		"LEFT JOIN `user` as consumer ON consumer.uid= spread_logs.uid " +
		"WHERE" + w +
		"ORDER BY `spread_logs`.`created` DESC " +
		"LIMIT ?, ?", offset, length).
		Find(&items)
	return
}

func (s *SpreadService) QRCode(rid int64, category string) (qrcode string, err error) {
	var info = &data.SpreadAccess{}
	has, err := data.Db.Where(fmt.Sprintf("relation_id = %d AND category = '%s'", rid, category)).Get(info)
	if !has {
		// 不存在 需要生成
		info.Category = category
		info.RelationID = rid
		affected, err := data.Db.InsertOne(info)
		if affected == 0 || err != nil {
			return "", err
		}
	}
	if info.QRCode == "" || category == `user` {
		// 获取access_token
		tk, err := module.Wechat.GetAccessToken()
		if err != nil {
			return "", err
		}
		var (
			b string
		)
		sceneId := fmt.Sprintf("%d", info.AccessID)
		switch category {
		case "channel":
			b = `{"action_name":"QR_LIMIT_SCENE","action_info":{"scene":{"scene_id":` + sceneId + `}}}`
		case "user":
			b = `{"expire_seconds": 2592000, "action_name": "QR_SCENE", "action_info": {"scene": {"scene_id":` + sceneId + `}}}`
		}
		uri := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=" + tk
		client := http.DefaultClient
		resp, err := client.Post(uri, "application/json;charset=utf-8", bytes.NewBuffer([]byte(b)))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		var ticket struct {
			Ticket string `json:"ticket"`
		}
		json.NewDecoder(resp.Body).Decode(&ticket)
		if ticket.Ticket == "" {
			return qrcode, errors.New("request ticket failed.")
		}
		resp, _ = client.Get("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(ticket.Ticket))
		defer resp.Body.Close()
		head, _ := ioutil.ReadAll(resp.Body)
		info.QRCode = base64.StdEncoding.EncodeToString(head)
		_, err = data.Db.Where("access_id = ?", info.AccessID).Update(info)
		if err != nil {
			return "", err
		}
	}
	return info.QRCode, nil
}

// 处理销售提成
func (s *SpreadService) Commission (rid, total, uid, orderId int64) (ok bool) {
	var info = &data.SpreadAccess{}
	if has , _ := data.Db.Where("access_id = ?", rid).Get(info); !has {
		return
	}
	switch info.Category {
	case "channel":
		// 渠道
		plan, has :=module.Spread.Plan(info.RelationID)
		if has && plan.IsDisabled == 0 {
			spreadLogs := &data.SpreadLogs{
				PlanID:plan.SPID,
				UID:uid,
				Category:2,
				OrderTotal:total,
				Commission:total * plan.SaleCommission / 100,
			}
			data.Db.Insert(spreadLogs)
			return true
		}
	case "user":
		// 普通用户推广
		if has, userInfo := module.User.GetInfoByUID(info.RelationID); has {
			_, sale := module.Conf.GetInt64("SPREAD_SALE")
			userInfo.PointsBalance += sale
			userInfo.PointsEarning += sale
			earning := &data.PointsLog{
				UID:userInfo.UID,
				RelationUID:userInfo.UID,
				RelationLogID:orderId,
				FriendlyIntro:"推广销售获得积分",
				Total:sale,
				Type:1,
			}
			data.Db.Where("uid = ?", userInfo.UID).Update(userInfo)
			data.Db.Insert(earning)
			return true
		}
	}
	return
}

func (s *SpreadService) generatePlanParam(state int) (w string) {
	w = "1=1"
	if state != -1 {
		w += fmt.Sprintf(" AND is_disabled = %d", state)
	}
	return
}
