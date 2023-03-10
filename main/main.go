package main

import (
	"github.com/LexaTRex/timetravelDB/api"
	"github.com/LexaTRex/timetravelDB/utils"
)

func main() {
	utils.Debug("starting TTQL application...")
	utils.DEBUG = true
	api.Api()
}
