package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/shumkovdenis/currency-server/rate"
)

type Server struct {
	service rate.Service
}

func Start(port int, service rate.Service) {
	s := &Server{service}

	http.HandleFunc("/last", s.last)
	http.HandleFunc("/rates", s.rates)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) last(w http.ResponseWriter, r *http.Request) {
	err := s.service.Error()
	if err != nil {
		log.Error("Rate service error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rate, err := s.service.LastRate()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(rate)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)
}

func (s *Server) rates(w http.ResponseWriter, r *http.Request) {
	err := s.service.Error()
	if err != nil {
		log.Error("Rate service error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rates, err := s.service.Rates()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(rates)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)
}
