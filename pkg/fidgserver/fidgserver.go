package fidgserver

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"timefidget/pkg/server"
	"timefidget/pkg/util"
)

type fidgserver struct {
	server *server.Server
}

func NewFidgserver() (*fidgserver, error) {

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	util.Logger = logger

	s, err := server.NewServer()
	if err != nil {
		return nil, err
	}

	return &fidgserver{
		server: s,
	}, nil
}
