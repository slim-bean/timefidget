package fidgserver

import (
	"flag"
	"os"

	"github.com/cortexproject/cortex/pkg/util/flagext"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/grafana/loki/pkg/cfg"
	"github.com/grafana/loki/pkg/loki"

	"timefidget/pkg/lokiembed"
	"timefidget/pkg/server"
	"timefidget/pkg/util"
)

type Config struct {
	loki.Config `yaml:",inline"`
	EmbedLoki   bool   `yaml:"embed_loki"`
	Port        int    `yaml:"port"`
	configFile  string `yaml:"-"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&c.configFile, "config.file", "", "yaml file to load")
	c.Config.RegisterFlags(f)
}

// Clone takes advantage of pass-by-value semantics to return a distinct *Config.
// This is primarily used to parse a different flag set without mutating the original *Config.
func (c *Config) Clone() flagext.Registerer {
	return func(c Config) *Config {
		return &c
	}(*c)
}

type fidgserver struct {
	server *server.Server
}

func NewFidgserver() (*fidgserver, error) {
	var config Config

	if err := cfg.Parse(&config); err != nil {
		return nil, err
	}

	var logger log.Logger

	if config.EmbedLoki {
		lokiembed.RunLoki(config.Config)
		lw, err := lokiembed.NewLogWriter(config.Config)
		if err != nil {
			return nil, err
		}
		logger = log.NewLogfmtLogger(log.NewSyncWriter(lw))
	} else {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	}

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
