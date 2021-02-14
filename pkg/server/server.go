package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log/level"

	"timefidget/pkg/model"
	"timefidget/pkg/util"
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
	// Declare a new Person struct.
	var p model.AccelDTO

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
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

	level.Info(util.Logger).Log("msg", fmt.Sprintf("AccelDTO: %+v", d))
}
