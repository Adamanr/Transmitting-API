package config

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetMessageTokenCheck Заготовка для установки Header на маршруты
func GetMessageTokenCheck(context *gin.Context) {
	if token := context.GetHeader("X-Master-Token"); token != "Dwk123mFl23G71falvcvmc" {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": errors.New("не прошёл проверку")})
	}
}

func SetHeaderPost(req *http.Request) {
	req.Header = http.Header{
		"Content-Type":   {"application/json"},
		"X-Master-Token": {"Dwk123mFl23G71falvcvmc"},
		"X-Auth-Token":   {"a1aad449-a79c-4983-b203-90a102a35e9b"},
	}
}
