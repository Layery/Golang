package server

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"skill-order/server/dao"
	"skill-order/server/model"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
}

func NewOrderController() *OrderController {
	return &OrderController{}
}

func (o *OrderController) SkillV1(c *gin.Context) {
	// 通过数据库实现的秒杀, 使用for update 倒是没超卖, 但是并发上不去,
	// 实测qps
	//	1. 判断库存, 如果<=0, 则返回下单失败, 库存不足
	//	2. 如果库存>0, 则下单, 减库存

	// goodsID, ok := c.Params.Get("goods_id")
	// if !ok {
	// 	c.JSON(http.StatusOK, gin.H{"error": "无效的商品id"})
	// 	return
	// }

	goodsID := "1"
	goodsNumer := 1

	tx := dao.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"msg": err})
		}
	}()

	var goods model.SkillGoods
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", goodsID).
		First(&goods)

	if goods.Stock <= 0 {
		panic("售罄了")
	}

	// 扣商品库存
	tx.Model(&model.SkillGoods{}).Where("id = ?", goodsID).Updates(map[string]interface{}{
		"stock": gorm.Expr("stock-?", goodsNumer),
	})

	// 生成订单,
	order := model.Order{
		OrderSn:     genOrderSn(),
		GoodsNumber: goodsNumer,
		GoodsId:     goods.Id,
	}
	tx.Create(&order)

	tx.Commit()

	c.JSON(200, gin.H{"msg": "success"})

}

func (o *OrderController) SkillV2(c *gin.Context) {
	// 使用乐观锁, 不使用for update测试一下
	// 实际测试下来, 如果不加事务, 那么会出现数据错乱的情况,
	// 加了事务, 虽然不会错乱, 但是qps极其的低, 才20几个
	// 并发场景下, 很多请求拿不到正确的版本号, 所以无法正常下单

	goodsID := 1
	goodsNumer := 1

	dao.DB.Transaction(func(tx *gorm.DB) error {
		goods := &model.SkillGoods{}
		tx.Model(goods).Where("id = ? ", goodsID).First(goods)

		if goods.Stock <= 0 {
			c.JSON(400, gin.H{"msg": "售罄了"})
			return errors.New("售罄了")
		}

		rows := tx.Model(goods).Where("updated_at = ?", goods.UpdatedAt).Updates(map[string]interface{}{
			"stock": gorm.Expr("stock-?", goodsNumer),
		}).RowsAffected
		if rows <= 0 {
			c.JSON(400, gin.H{"msg": "扣减库存失败"})
			return errors.New("扣减库存失败")
		}

		// 生成订单,
		order := model.Order{
			OrderSn:     genOrderSn(),
			GoodsNumber: goodsNumer,
			GoodsId:     goods.Id,
		}

		tx.Create(&order)
		c.JSON(200, gin.H{"msg": "success"})
		return nil
	})
}

func genOrderSn() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d", timestamp)
}
