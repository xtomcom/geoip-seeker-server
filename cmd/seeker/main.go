package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"os"

	"../../http-service"
)

var configure *Configure

type Configure struct {
	Listen   string             `json:"listen"`
	Services map[string]Service `json:"services"`
}

type Service struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func init() {
	data, err := ioutil.ReadFile("./seeker.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	configure = new(Configure)
	if err := json.Unmarshal(data, configure); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	mux := http.NewServeMux()
	for path, service := range configure.Services {
		handler, err := http_service.GetService(service.Name, service.Path)
		if err != nil {
			fmt.Println(fmt.Sprintf(
				"register failed %s (%s) to %s",
				service.Name,
				service.Path,
				path,
			))
			continue
		}
		mux.Handle(path+"/", http.StripPrefix(path, handler))
		fmt.Println(fmt.Sprintf(
			"register %s (%s) to %s",
			service.Name,
			service.Path,
			path,
		))
	}
	fmt.Println(fmt.Sprintf("listen: http://%s", configure.Listen))
	http.ListenAndServe(configure.Listen, mux)
}
