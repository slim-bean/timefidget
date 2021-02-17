package main

import (
	"bytes"
	"flag"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/cortexproject/cortex/pkg/util"
	"github.com/cortexproject/cortex/pkg/util/flagext"
	gklog "github.com/go-kit/kit/log"
	"github.com/go-logfmt/logfmt"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/promtail/api"
	"github.com/grafana/loki/pkg/promtail/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"

	fs_model "timefidget/pkg/model"
)

func main() {

	from := flag.String("from", "", "Start Time RFC339Nano 2006-01-02T15:04:05.999999999Z07:00")
	to := flag.String("to", "", "End Time RFC339Nano 2006-01-02T15:04:05.999999999Z07:00")
	project := flag.String("project", "", "source datasource config")
	flag.Parse()

	u, err := url.Parse("http://localhost:" + strconv.FormatInt(int64(5100), 10) + "/loki/api/v1/push")
	if err != nil {
		panic(err)
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
	logger := gklog.NewLogfmtLogger(gklog.NewSyncWriter(os.Stderr))
	c, err := client.New(prometheus.DefaultRegisterer, cfg, logger)
	if err != nil {
		panic(err)
	}

	f := mustParse(*from)
	t := mustParse(*to)
	ct := f
	buf := &bytes.Buffer{}
	enc := logfmt.NewEncoder(buf)
	err = enc.EncodeKeyvals("level", "info", "type", "sub", fs_model.ProjectName, *project)
	if err != nil {
		panic(err)
	}
	line := buf.String()
	for ct.Before(t) {
		e := api.Entry{
			Labels: model.LabelSet{
				"job":  "timefidget",
				"type": "sub",
			},
			Entry: logproto.Entry{
				Timestamp: ct,
				Line:      line,
			},
		}
		log.Println(e)
		c.Chan() <- e
		ct = ct.Add(5 * time.Second)
		time.Sleep(50 * time.Millisecond)
	}

}

func mustParse(t string) time.Time {

	ret, err := time.Parse(time.RFC3339Nano, t)

	if err != nil {
		log.Fatalf("Unable to parse time %v", err)
	}

	return ret
}
