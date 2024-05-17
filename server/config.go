package server

import (
	"github.com/noona-hq/app-template/db"
	"github.com/noona-hq/app-template/logger"
	"github.com/noona-hq/app-template/services/noona"
)

type Config struct {
	Noona  noona.Config
	Logger logger.Config
	DB     db.Config
	// Store can either be memory or mongodb
	Store string `default:"mongodb"`
}
