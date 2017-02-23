package pkg

import (
	"leaf/log"
	"net/http"
	"strconv"
)

// TODO:
//STATUS,TPS..ETC
type Result struct {
	result map[string]string
}

func Test(ag ApiGeter) {
	var api Api

	rs := &Result{
		result: make(map[string]string),
	}

	// TODO:now is circle
	//NOTIFY INTERRUPT
	for {
		ag(&api)

		request, err := api.Request()
		if err != nil {
			log.Error(err)
		}
		client := http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Error(err)
		}

		if response != nil && api.ret == strconv.Itoa(response.StatusCode) {
			rs.result["SUCESS"] += api.URL
		} else {
			rs.result["FAIL"] += api.URL
		}
	}
}
