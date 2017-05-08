package http_service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"../geoip-seeker/qqwry"
)

type QQWayHandler struct {
	seeker *qqwry.IPSeeker
}

func makeQQWryHandler(path string) (http.Handler, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	path = strings.ToLower(path)
	mux := http.NewServeMux()
	handler := new(QQWayHandler)
	handler.seeker = qqwry.New(data)
	mux.HandleFunc("/", handler.query)
	mux.HandleFunc("/status", handler.status)
	return mux, nil
}

func (h *QQWayHandler) query(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ip := getRemoteIPv4(r)
	location, err := h.seeker.LookupByIP(ip)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(err.Error()))
		return
	}
	if r.Header.Get("Content-Type") == "text/plain" {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(gbkToUTF8(location.String())))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	location.Country = gbkToUTF8(location.Country)
	location.Area = gbkToUTF8(location.Area)
	payload, _ := json.Marshal(location)
	rw.Write(payload)
}

func (h *QQWayHandler) status(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	status := new(struct {
		RecordCount int    `json:"record_count"`
		Version     string `json:"publish_date"`
	})
	status.RecordCount = h.seeker.RecordCount()
	status.Version = gbkToUTF8(h.seeker.Version())

	rw.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(status)
	rw.Write(payload)
}
