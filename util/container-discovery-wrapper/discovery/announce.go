package discovery

import (
	"github.com/dghubble/sling"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"fmt"
	"time"
)

type serviceStruct struct {
	ID string
	Service string
	Tags []string
	Address string
	TaggedAddresses map[string]string
	Port uint16
}

type request struct {
	Node string
	Address string
	Service serviceStruct
}

func AnnounceService(
	DiscoveryServer, SCInstanceIdentifier string,
	ServiceName, ServiceInstanceId string,
	ServiceExposedPort uint16,
) bool {
	if ServiceInstanceId == "" {
		ServiceInstanceId = os.Getenv("HOSTNAME")
	}

	println("Registering the service into: ", DiscoveryServer, "with the identifier: ", ServiceInstanceId)

	realExposedPort, _ := strconv.Atoi(
		os.Getenv("DOCKER_CONTAINER_PORT_" + strconv.Itoa(int(ServiceExposedPort)) + "_REAL"),
	)

	if realExposedPort == 0 {
		realExposedPort = 666 //todo: hacky as f***, will figure it out later how to register services with no
				      //      ports exposed
	}

	return doAnnounceServiceRequest(DiscoveryServer, &request{
		Node: os.Getenv("DOCKER_HOST_NAME"),
		Address: os.Getenv("DOCKER_HOST_IP"),
		Service: serviceStruct{
			ID: ServiceInstanceId,
			Service: ServiceName,
			Address: os.Getenv("DOCKER_HOST_IP"),
			Port: uint16(realExposedPort),
			Tags: []string{"sc-instance-identifier::" + SCInstanceIdentifier},
			TaggedAddresses: map[string]string {},
		},
	})
}

func doAnnounceServiceRequest(DiscoveryServer string, reqBody *request) bool {
	req, err := sling.New().Base(DiscoveryServer).Put("/v1/catalog/register").BodyJSON(reqBody).Request()

	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	response, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	b, _ := ioutil.ReadAll(response.Body)

	if string(b) == "No cluster leader" { //Probably the Consul server didn't fully start yet.
		println("Consul is up but no leader yet, I will sleep and retry after 5 seconds...")
		time.Sleep(5 * time.Second)
		return doAnnounceServiceRequest(DiscoveryServer, reqBody)
	} else if string(b) != "true" {
		println("Error announcing the service: ")
		println(string(b))
		fmt.Printf("%+v", reqBody)
		os.Exit(1)
	}

	return true
}