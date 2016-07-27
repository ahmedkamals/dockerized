package discovery

import (
	"os"
	"github.com/dghubble/sling"
	"net/http"
	"io/ioutil"
	"fmt"
)

func DeregisterService(
	DiscoveryServer, ServiceInstanceId string,
) {
	if ServiceInstanceId == "" {
		ServiceInstanceId = os.Getenv("HOSTNAME")
	}

	reqBody := struct{
		Datacenter string
		Node string
		ServiceID string
	} {
		Datacenter: "dc1",
		Node: os.Getenv("DOCKER_HOST_NAME"),
		ServiceID: ServiceInstanceId,
	}

	req, err := sling.New().Base(DiscoveryServer).Put("/v1/catalog/deregister").BodyJSON(reqBody).Request()

	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	response, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	b, _ := ioutil.ReadAll(response.Body)

	if string(b) != "true" {
		println("Error deregistering the service: ")
		println(string(b))
		fmt.Printf("%+v", reqBody)
		os.Exit(1)
	}
}
