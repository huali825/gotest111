package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"goworkwebook/webook003/internal/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
}
