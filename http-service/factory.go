package http_service

import (
	"fmt"
	"net/http"
)

var mapping = map[string]func(path string) (http.Handler, error){
	"ipip.net": makeIPIPNetHandler,
	"qqwry":    makeQQWryHandler,
}

func GetService(name, path string) (http.Handler, error) {
	if handle, ok := mapping[name]; ok {
		return handle(path)
	}
	return nil, fmt.Errorf("service %s not found", name)
}
