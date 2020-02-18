package context

import (
	"github.com/Lyt99/iop-statistics/config"
	"github.com/Lyt99/iop-statistics/model"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Success: true,
		Data:    data,
		Error:   "",
	})
}

func SuccessWithCache(c *gin.Context, key string, data interface{}) {
	resp := response{
		Success: true,
		Data:    data,
		Error:   "",
	}

	jString, err := jsoniter.Marshal(resp)
	if err != nil {
		Error(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	_ = model.RedisClient.Set(key, jString, 0)

	c.Data(http.StatusOK, gin.MIMEJSON, jString)
}

func TryResponseCache(c *gin.Context, key string) bool {
	jString, err := model.RedisClient.Get(key).Bytes()

	if err != nil {
		return false
	}

	c.Data(http.StatusOK, gin.MIMEJSON, jString)
	return true
}

func Error(c *gin.Context, status int, msg string, err error) {
	ret := response{
		Success: false,
		Data:    nil,
		Error:   msg,
	}
	if config.GlobalConfig.EnableDebug && err != nil {
		ret.Error += ": " + err.Error()
	}

	c.JSON(status, ret)
}
