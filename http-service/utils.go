package http_service

import (
	"io/ioutil"
	"net"
	"net/http"

	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func getRemoteIPv4(r *http.Request) net.IP {
	ip := net.ParseIP(r.URL.Query().Get("ip")).To4()
	if ip == nil {
		ip = net.ParseIP(r.RemoteAddr).To4()
	}
	return ip
}
func gbkToUTF8(value string) string {
	reader := transform.NewReader(strings.NewReader(value), simplifiedchinese.GBK.NewDecoder())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return ""
	}
	return string(data)
}
