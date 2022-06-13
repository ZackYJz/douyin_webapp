package service

import (
	"go_webapp/global"
	"go_webapp/internal/dao"

	"github.com/gin-gonic/gin"
)

type Service struct {
	c   *gin.Context
	dao *dao.Dao
}

func New(c *gin.Context) Service {
	svc := Service{c: c}
	svc.dao = dao.New(global.DBEngine)
	return svc
}
