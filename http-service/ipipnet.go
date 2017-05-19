package http_service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"time"

	"../geoip-seeker/ipip.net"
)

type IPIPNetHandler struct {
	seeker *ipip_net.IPSeeker
}

func makeIPIPNetHandler(path string) (http.Handler, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	path = strings.ToLower(path)

	mux := http.NewServeMux()
	handler := new(IPIPNetHandler)
	switch {
	case strings.HasSuffix(path, ".dat"):
		handler.seeker, err = ipip_net.NewDAT(data)
		if err != nil {
			return nil, err
		}
	case strings.HasSuffix(path, ".datx"):
		handler.seeker, err = ipip_net.NewDATX(data)
		if err != nil {
			return nil, err
		}
	}

	mux.HandleFunc("/", handler.query)
	mux.HandleFunc("/status", handler.status)
	return mux, nil
}

func (handler *IPIPNetHandler) query(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ip := getRemoteIPv4(r)
	location, err := handler.seeker.LookupByIP(ip)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(err.Error()))
		return
	}
	if r.Header.Get("Content-Type") == "text/plain" {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(location.String()))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(location)
	rw.Write(payload)
}

func (handler *IPIPNetHandler) status(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	status := new(struct {
		RecordCount int       `json:"record_count"`
		PublishDate time.Time `json:"publish_date"`
	})
	status.RecordCount = handler.seeker.RecordCount()
	status.PublishDate = handler.seeker.PublishDate()

	rw.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(status)
	rw.Write(payload)
}
