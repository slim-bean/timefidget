package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"timefidget/pkg/model"
)

type Server struct {
}

func NewServer(port int) (*Server, error) {
	s := &Server{}

	http.HandleFunc("/push", s.pushHandler)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", "0.0.0.0", port), nil)

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
		ID: p.ID,
		X:  xf,
		Y:  yf,
		Z:  zf,
	}

	model.LogPosition(d.ID, d.X, d.Y, d.Z)

	//level.Info(util.Logger).Log("msg", fmt.Sprintf("AccelDTO: %+v", d))
}
