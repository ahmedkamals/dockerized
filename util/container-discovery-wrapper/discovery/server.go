package discovery

import (
	"github.com/dghubble/sling"
	"io/ioutil"
	"net/http"
)

func ServerIsUp(serverUrl string) bool {
	req, err := sling.New().Base(serverUrl).Get("/v1/status/leader").Request()

	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	response, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	b, _ := ioutil.ReadAll(response.Body)

	return string(b) != ""
}
