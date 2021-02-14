package lokiembed

import (
	"net/url"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/cortexproject/cortex/pkg/util"
	"github.com/cortexproject/cortex/pkg/util/flagext"
	"github.com/cortexproject/cortex/pkg/util/log"
	"github.com/go-kit/kit/log/level"
	"github.com/grafana/loki/pkg/loki"
	"github.com/grafana/loki/pkg/util/validation"
	"github.com/prometheus/common/model"
	"github.com/prometheus/common/version"
	"github.com/weaveworks/common/logging"

	util_log "github.com/cortexproject/cortex/pkg/util/log"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/promtail/api"
	"github.com/grafana/loki/pkg/promtail/client"
	"github.com/prometheus/client_golang/prometheus"
)

type LogWriter struct {
	client client.Client
}

func NewLogWriter(config loki.Config) (*LogWriter, error) {
	u, err := url.Parse("http://localhost:" + strconv.FormatInt(int64(config.Server.HTTPListenPort), 10) + "/loki/api/v1/push")
	if err != nil {
		return nil, err
	}
	cfg := client.Config{
		URL: flagext.URLValue{
			URL: u,
		},
		BatchWait: client.BatchWait,
		BatchSize: client.BatchSize,
		BackoffConfig: util.BackoffConfig{
			MinBackoff: client.MinBackoff,
			MaxBackoff: client.MaxBackoff,
			MaxRetries: client.MaxRetries,
		},
		Timeout: client.Timeout,
	}
	c, err := client.New(prometheus.DefaultRegisterer, cfg, log.Logger)
	if err != nil {
		return nil, err
	}
	w := &LogWriter{
		client: c,
	}
	return w, nil
}

func (l *LogWriter) Stop() {
	l.client.Stop()
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	level.Info(util_log.Logger).Log("msg", "sending to Loki")
	//d := logfmt.NewDecoder(p)
	//var project string
	//for d.ScanKeyval() {
	//	if string(d.Key()) == fidgmodel.ProjectName {
	//		project = string(d.Value())
	//	}
	//}
	str := string(p)
	e := api.Entry{
		Labels: model.LabelSet{
			"job": "timefidget",
			//fidgmodel.ProjectName: model.LabelValue(project),
		},
		Entry: logproto.Entry{
			Timestamp: time.Now(),
			Line:      str,
		},
	}
	l.client.Chan() <- e
	level.Info(util_log.Logger).Log("msg", "sent to Loki", "entry", str)
	return len(p), nil
}

func RunLoki(config loki.Config, wg *sync.WaitGroup) {

	// This global is set to the config passed into the last call to `NewOverrides`. If we don't
	// call it atleast once, the defaults are set to an empty struct.
	// We call it with the flag values so that the config file unmarshalling only overrides the values set in the config.
	validation.SetDefaultLimitsForYAMLUnmarshalling(config.LimitsConfig)

	// Init the logger which will honor the log level set in config.Server
	if reflect.DeepEqual(&config.Server.LogLevel, &logging.Level{}) {
		level.Error(util_log.Logger).Log("msg", "invalid log level")
		os.Exit(1)
	}
	util_log.InitLogger(&config.Server)

	// Validate the config once both the config file has been loaded
	// and CLI flags parsed.
	err := config.Validate(util_log.Logger)
	if err != nil {
		level.Error(util_log.Logger).Log("msg", "validating config", "err", err.Error())
		os.Exit(1)
	}

	// Start Loki
	t, err := loki.New(config)
	util_log.CheckFatal("initialising loki", err)

	level.Info(util_log.Logger).Log("msg", "Starting Loki", "version", version.Info())

	go func() {
		wg.Add(1)
		defer wg.Done()
		err = t.Run()
		util_log.CheckFatal("running loki", err)
	}()
}
