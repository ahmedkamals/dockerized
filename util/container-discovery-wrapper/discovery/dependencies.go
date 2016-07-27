package discovery

import (
	"github.com/dghubble/sling"
	"io/ioutil"
	"net/http"
	"time"
	"os"
	"strings"
	"encoding/json"
	"strconv"
)

func BlockForDependencies(serverUrl, SCInstanceIdentifier string, dependencies []string) {
	if len(dependencies) == 0 {
		return
	}

	query := struct {
		Tag string `url:"tag"`
	}{
		Tag: "sc-instance-identifier::" + SCInstanceIdentifier,
	}

	//Refactor to use blocking requests in concurrent go routines
	for _, dependencyServiceName := range dependencies {
		for {
			req, err := sling.New().Base(serverUrl).Get("/v1/catalog/service/" + dependencyServiceName).QueryStruct(query).Request()

			if err != nil {
				panic(err)
			}

			response, err := http.DefaultClient.Do(req)

			if err != nil {
				panic(err)
			}

			b, _ := ioutil.ReadAll(response.Body)
			dependencyFound := string(b) != "[]" && string(b) != "No cluster leader"

			if dependencyFound == false {
				println("Dependency not met yet:", dependencyServiceName, ". I will sleep and retry after 2 seconds...")
				time.Sleep(2 * time.Second)
			} else {
				var services []map[string]interface{}
				json.Unmarshal(b, &services)
				//
				//println(services[0]["ServiceAddress"].(string))
				//println(strconv.FormatFloat(services[0]["ServicePort"].(float64), 'f', 0, 64))
				println("SERVICE_"+strings.ToUpper(strings.Replace(dependencyServiceName, "-", "_", -1))+"_ADDRESS_IP")

				os.Setenv(
					"SERVICE_"+strings.ToUpper(strings.Replace(dependencyServiceName, "-", "_", -1))+"_ADDRESS_IP",
					services[0]["ServiceAddress"].(string),
				)

				os.Setenv(
					"SERVICE_"+strings.ToUpper(strings.Replace(dependencyServiceName, "-", "_", -1))+"_ADDRESS_PORT",
					strconv.FormatFloat(services[0]["ServicePort"].(float64), 'f', 0, 64),
				)

				break
			}
		}
	}
}
