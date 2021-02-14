package server

import "net/http"

type Config struct {
}

type server struct {
}

func NewServer() (*server, error) {
	s := &server{}

	http.HandleFunc("/push", s.pushHandler)
	go http.ListenAndServe(":8080", nil)

	return s, nil
}


func (s *server) pushHandler(w http.ResponseWriter, req *http.Request) {

}
