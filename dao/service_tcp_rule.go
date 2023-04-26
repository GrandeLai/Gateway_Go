package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceID int64 `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port      int   `json:"port" gorm:"column:port" description:"端口	"`
}

func (t *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (t *TcpRule) Find(c *gin.Context, tx *gorm.DB, search *TcpRule) (*TcpRule, error) {
	model := &TcpRule{}
	err := tx.Where(search).Find(model).Error
	return model, err
}

func (t *TcpRule) TcpRuleIsExist(tx *gorm.DB, search *TcpRule) bool {
	var count int64
	tx.Where(search).Count(&count)
	return count > 0
}

func (t *TcpRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *TcpRule) Update(c *gin.Context, tx *gorm.DB) error {
	return tx.Where("id = ?", t.ID).Updates(&t).Error
}

func (t *TcpRule) ListByServiceID(c *gin.Context, tx *gorm.DB, serviceID int64) ([]TcpRule, int64, error) {
	var list []TcpRule
	var count int64
	query := tx.Table(t.TableName()).Select("*")
	query = query.Where("service_id=?", serviceID)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}
