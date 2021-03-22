package main

import (
	"bytes"
	"flag"
	"log"
	"net/url"
	"os"
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
	write := flag.Bool("write", false, "send output to Loki, false sends to stderr to review, true writes to Loki and stderr")
	versionLabel := flag.String("version", "", "Set a value for a `version` label, used to add entries and avoid out of order errors")
	typeLabelVal := flag.String("typeLabelVal", "sub", "Set a value for the `type` label, default `sub` is used by dashboards for subtracting errors, useful for testing, requires write=true to send to Loki")
	flag.Parse()

	u, err := url.Parse("https://loki-personal.edjusted.com/loki/api/v1/push")
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
	err = enc.EncodeKeyvals("level", "info", "type", *typeLabelVal, fs_model.ProjectName, *project)
	if err != nil {
		panic(err)
	}
	line := buf.String()
	for ct.Before(t) {
		lbls := model.LabelSet{
			"job":  "timefidget",
			"type": model.LabelValue(*typeLabelVal),
		}
		if *versionLabel != "" {
			lbls["version"] = model.LabelValue(*versionLabel)
		}
		e := api.Entry{
			Labels: lbls,
			Entry: logproto.Entry{
				Timestamp: ct,
				Line:      line,
			},
		}
		log.Println(e)
		if *write {
			c.Chan() <- e
		}
		ct = ct.Add(5 * time.Second)
		time.Sleep(25 * time.Millisecond)
	}
	c.Stop()

}

func mustParse(t string) time.Time {

	ret, err := time.Parse(time.RFC3339Nano, t)

	if err != nil {
		log.Fatalf("Unable to parse time %v", err)
	}

	return ret
}
