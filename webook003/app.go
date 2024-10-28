package main

import (
	"github.com/gin-gonic/gin"
	"goworkwebook/webook003/internal/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
