package main

import (
	"bytes"
	"flag"
	"github.com/go-logfmt/logfmt"
	"github.com/grafana/loki-client-go/pkg/backoff"
	"github.com/grafana/loki-client-go/pkg/urlutil"
	"github.com/prometheus/common/model"
	"log"
	"net/url"
	"time"

	client "github.com/grafana/loki-client-go/loki"

	fs_model "github.com/slim-bean/timefidget/pkg/model"
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
		URL: urlutil.URLValue{
			URL: u,
		},
		BatchWait: client.BatchWait,
		BatchSize: client.BatchSize,
		BackoffConfig: backoff.BackoffConfig{
			MinBackoff: client.MinBackoff,
			MaxBackoff: client.MaxBackoff,
			MaxRetries: client.MaxRetries,
		},
		Timeout: client.Timeout,
	}
	c, err := client.New(cfg)
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
		log.Printf("Entry: %v %s %s\n", ct, lbls, line)
		if *write {
			err := c.Handle(lbls, ct, line)
			if err != nil {
				log.Println("error sending log:", err)
			}
		}
		ct = ct.Add(5 * time.Second)
		time.Sleep(1 * time.Millisecond)
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
