package cmd

import (
	"github.com/gin-gonic/gin"
	"sync"
	"test/mainworker/structs"
	"testing"
)

func TestGlobalStorage_GetUsers(t *testing.T) {
	type fields struct {
		Users map[string]*structs.User
		Mu    *sync.RWMutex
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GlobalStorage{
				Users: tt.fields.Users,
				Mu:    tt.fields.Mu,
			}
			gs.GetUsers(tt.args.ctx)
		})
	}
}
