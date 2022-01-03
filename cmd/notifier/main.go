package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type (

	// HookMessage is the message we receive from Alertmanager
	HookMessage struct {
		Version           string            `json:"version"`
		GroupKey          string            `json:"groupKey"`
		Status            string            `json:"status"`
		Receiver          string            `json:"receiver"`
		GroupLabels       map[string]string `json:"groupLabels"`
		CommonLabels      map[string]string `json:"commonLabels"`
		CommonAnnotations map[string]string `json:"commonAnnotations"`
		ExternalURL       string            `json:"externalURL"`
		Alerts            []Alert           `json:"alerts"`
	}

	// Alert is a single alert.
	Alert struct {
		Status       string            `json:"status"`
		Labels       map[string]string `json:"labels"`
		Annotations  map[string]string `json:"annotations"`
		StartsAt     string            `json:"startsAt,omitempty"`
		EndsAt       string            `json:"EndsAt,omitempty"`
		GeneratorURL string            `json:"generatorURL"`
	}

	// just an example alert store. in a real hook, you would do something useful
	alertStore struct {
		sync.Mutex
		w1firing bool
		w2firing bool
		w2Pin    rpio.Pin
		w1Pin    rpio.Pin
	}
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen for webhook")
	flag.Parse()

	rpio.Open()
	defer rpio.Close()

	// Pins 18/19 were chosen intentionally because there are only 2 PWM channels on the Pi
	// and these pins are each on one of those separate channels
	w1Pin := rpio.Pin(19)
	w1Pin.Mode(rpio.Pwm)
	w1Pin.Freq(64000)
	w1Pin.DutyCycleWithPwmMode(0, 32, rpio.Balanced)

	w2Pin := rpio.Pin(18)
	w2Pin.Mode(rpio.Pwm)
	w2Pin.Freq(64000)
	w2Pin.DutyCycleWithPwmMode(0, 32, rpio.Balanced)

	s := &alertStore{
		w1Pin: w1Pin,
		w2Pin: w2Pin,
	}

	go s.ledhandler()

	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/alerts", s.postHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func (s *alertStore) postHandler(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var m HookMessage
	if err := dec.Decode(&m); err != nil {
		log.Printf("error decoding message: %v", err)
		http.Error(w, "invalid request body", 400)
		return
	}

	s.Lock()
	defer s.Unlock()

	log.Println("Received Notification:", m)

	for _, a := range m.Alerts {
		for k, v := range a.Labels {
			if k == "alertname" {
				firing := false
				switch a.Status {
				case "firing":
					firing = true
				case "resolved":
					firing = false
				}
				switch v {
				case "not-tracking-time-w1":
					log.Println("Setting W1 LED:", a.Status)
					s.w1firing = firing
				case "not-tracking-time-w2":
					log.Println("Setting W2 LED:", a.Status)
					s.w2firing = firing
				}
			}
		}
	}
}

func (s *alertStore) ledhandler() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Lock()
			if s.w1firing {
				s.w1Pin.DutyCycleWithPwmMode(4, 32, rpio.Balanced)
			} else {
				s.w1Pin.DutyCycleWithPwmMode(0, 32, rpio.Balanced)
			}
			if s.w2firing {
				s.w2Pin.DutyCycleWithPwmMode(4, 32, rpio.Balanced)
			} else {
				s.w2Pin.DutyCycleWithPwmMode(0, 32, rpio.Balanced)
			}
			s.Unlock()
		}
	}
}
