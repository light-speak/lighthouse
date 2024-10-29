package service

import (
	_ "user/repo"

	"github.com/light-speak/lighthouse/handler"
)


func StartService() {
	handler.StartService()
}
