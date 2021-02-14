package server

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log/level"

	"timefidget/pkg/model"
	"timefidget/pkg/util"
)

const (
	ON_MIN   = 8
	ON_MAX   = 11
	OFF_MIN  = -1
	OFF_MAX  = 1
	HALF_MIN = 5
	HALF_MAX = 8
	Z_THRESH = 5

	P0 = ""
	P1 = "Slack"
	P2 = "Loki Community"
	P3 = "Loki Ops"
	P4 = "Loki"
	P5 = "Sales"
	P6 = "1-1"
	P7 = "Management"
	P8 = "BAU"
)

type Config struct {
}

type Server struct {
}

func NewServer() (*Server, error) {
	s := &Server{}

	http.HandleFunc("/push", s.pushHandler)
	go http.ListenAndServe("0.0.0.0:8080", nil)

	return s, nil
}

func (s *Server) pushHandler(w http.ResponseWriter, req *http.Request) {

	var p model.AccelDTO

	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	xf, err := strconv.ParseFloat(p.X, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	yf, err := strconv.ParseFloat(p.Y, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	zf, err := strconv.ParseFloat(p.Z, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d := model.Accel{
		X: xf,
		Y: yf,
		Z: zf,
	}

	if math.Abs(d.Z) > Z_THRESH {
		// Off
		//level.Info(util.Logger).Log("pos", "0")
	} else if d.X > OFF_MIN && d.X < OFF_MAX && d.Y < -ON_MIN && d.Y > -ON_MAX {
		// Position 1
		level.Info(util.Logger).Log("pos", "1", model.ProjectName, P1)
	} else if d.X > HALF_MIN && d.X < HALF_MAX && d.Y < -HALF_MIN && d.Y > -HALF_MAX {
		// Position 2
		level.Info(util.Logger).Log("pos", "2", model.ProjectName, P2)
	} else if d.X > ON_MIN && d.X < ON_MAX && d.Y > OFF_MIN && d.Y < OFF_MAX {
		// Position 3
		level.Info(util.Logger).Log("pos", "3", model.ProjectName, P3)
	} else if d.X > HALF_MIN && d.X < HALF_MAX && d.Y > HALF_MIN && d.Y < HALF_MAX {
		// Position 4
		level.Info(util.Logger).Log("pos", "4", model.ProjectName, P4)
	} else if d.X > OFF_MIN && d.X < OFF_MAX && d.Y > ON_MIN && d.Y < ON_MAX {
		// Position 5
		level.Info(util.Logger).Log("pos", "5", model.ProjectName, P5)
	} else if d.X < -HALF_MIN && d.X > -HALF_MAX && d.Y > HALF_MIN && d.Y < HALF_MAX {
		// Position 6
		level.Info(util.Logger).Log("pos", "6", model.ProjectName, P6)
	} else if d.X < -ON_MIN && d.X > -ON_MAX && d.Y > OFF_MIN && d.Y < ON_MIN {
		// Position 7
		level.Info(util.Logger).Log("pos", "7", model.ProjectName, P7)
	} else if d.X < -HALF_MIN && d.X > -HALF_MAX && d.Y < -HALF_MIN && d.Y > -HALF_MAX {
		// Position 8
		level.Info(util.Logger).Log("pos", "8", model.ProjectName, P8)
	}

	//level.Info(util.Logger).Log("msg", fmt.Sprintf("AccelDTO: %+v", d))
}
