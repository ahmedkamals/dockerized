package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"rocket-labs/container-discovery-wrapper/discovery"
	"os/exec"
	//"bytes"
	"time"
	"strings"
	"io/ioutil"
	"io"
	"bufio"
	"os/signal"
	"syscall"
	"github.com/dghubble/sling"
	"net/http"
	"errors"
	"rocket-labs/container-discovery-wrapper/dockerports"
)

func main() {
	var opts struct {
		// Discovery endpoint, for example: 'http://discovery-server.qa.sellercenter.net:8500/', the plan is
		// to use a Consul server that can be used for both Service Discovery and a centralized Key-Value store.
		DiscoveryServer string `required:"true" long:"discovery-server" description:"Discovery endpoint, for example: 'http://discovery-server.qa.sellercenter.net:8500/', the plan is to use a Consul server that can be used for both Service Disovery and a centralized Key-Value store."`

		// Command to run before announcing the service, the process out of this command will be monitored and
		// terminated/killed whenever the application get a SIGTERM/SIGKILL.
		Command string `required:"true" long:"command" description:"Command to run before announcing the service, the process out of this command will be monitored and terminated/killed whenever the application get a SIGTERM/SIGKILL"`

		// Unique identifier for the Seller Center installation, for example: 'lazada-sg-live',
		// 'sandbox-seller-[UUID]', etc...
		SCInstanceIdentifier string `required:"true" long:"sc-instance-identifier" description:"Unique identifier for the Seller Center installation, for example: 'lazada-sg-live', 'sandbox-seller-[UUID]', etc..."`

		// Service name, you can have multiple services with the same name. Examples for value:
		// 'store_relational', 'store_cache', 'frontend_gui', etc..
		ServiceName string `required:"true" long:"service-name" description:"Service name, you can have multiple services with the same name. Examples for value: 'store_relational', 'store_cache', 'frontend_gui', etc.."`

		// Port number that's exposed from the container.
		ServiceExposedPort uint16 `required:"true" long:"service-exposed-port" description:"Port number that's exposed from the container."`

		// Service Instance Identifier, this must be unique across the Seller Center instance. Default value is
		// the Docker container identifier.
		ServiceInstanceId string `long:"service-instance-id" description:"Service Identifier, this must be unqiue accross the Seller Center instance. Default value is the Docker container identifier."`

		// If supplied, the service is only announced after a matching text is found in the process Stdout
		AnnounceAfterOutputText string `long:"announce-after-output-text" description:"If supplied, the service is only announced after a matching text is found in the process stdoutput."`

		// If supplied, the wrapper will not parse Port mappings from /docker.sock, only used for dev/testing
		MockServicePortMapping bool `long:"mock-service-port-mapping" description:"DONT USE! If supplied, the wrapper will not parse Port mappings from /docker.sock, only used for dev/testing."`

		// If supplied, the service is only ran and announced after the listed services are up and could be
		// fetched from the discovery server.
		ServiceDependencies []string `long:"service-dependency" description:"If supplied, the service is only ran and announced after the listed services are up and could be fetched from the discovery server."`
	}
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	if os.Getenv("IS_AWS_ECS") == "1" {
		req, err := sling.New().Get("http://169.254.169.254/latest/meta-data/local-ipv4").Request()
		if err != nil {
			panic(err)
		}
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		ipAddress, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		if string(ipAddress) == "" {
			panic(errors.New("Couldn't figure out the ECS instance IP address."))
		}
		os.Setenv("DOCKER_HOST_NAME", string(ipAddress))
		os.Setenv("DOCKER_HOST_IP", string(ipAddress))
	}

	if opts.MockServicePortMapping == false {
		dockerports.ParseAndSetenv()
	}

	if discovery.ServerIsUp(opts.DiscoveryServer) == false {
		println("Couldn't connect to the Discovery Server, make sure", opts.DiscoveryServer, "is up and running")
		os.Exit(1)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-sigc

		println("De-registering the service before quiting...")
		discovery.DeregisterService(
			opts.DiscoveryServer,
			opts.ServiceInstanceId,
		)
		println("Hasta la vista, baby")
		os.Exit(0)
	}()

	discovery.BlockForDependencies(
		opts.DiscoveryServer,
		opts.SCInstanceIdentifier,
		opts.ServiceDependencies,
	)

	cmdParts := strings.Split(opts.Command, "| |")
	cmdExecutable := cmdParts[0]
	cmdArgs := cmdParts[1:]

	//o, err := exec.Command(cmdExecutable, cmdArgs...).CombinedOutput()
	//
	//println(string(o))
	//if err != nil {
	//	println(err.Error())
	//}
	//os.Exit(0)

	cmd := exec.Command(cmdExecutable, cmdArgs...)
	//cmd := exec.Command("/bin/bash", "/entrypoint.sh mysqld --sql_mode=\"\"")
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	announced := false
	if opts.AnnounceAfterOutputText != "" {
		scanner := bufio.NewScanner(cmdStdout)
		go func() {
			for scanner.Scan() {
				t := scanner.Text()
				triggerFound := strings.Contains(t, opts.AnnounceAfterOutputText)
				println(t)
				if triggerFound {
					announced = discovery.AnnounceService(
						opts.DiscoveryServer,
						opts.SCInstanceIdentifier,
						opts.ServiceName,
						opts.ServiceInstanceId,
						opts.ServiceExposedPort,
					)
				}
			}
		}()

		errScanner := bufio.NewScanner(cmdStderr)
		go func() {
			for errScanner.Scan() {
				t := errScanner.Text()
				triggerFound := strings.Contains(t, opts.AnnounceAfterOutputText)
				println("Stderr:", t)
				if triggerFound {
					announced = discovery.AnnounceService(
						opts.DiscoveryServer,
						opts.SCInstanceIdentifier,
						opts.ServiceName,
						opts.ServiceInstanceId,
						opts.ServiceExposedPort,
					)
				}
			}
		}()
	} else {
		go io.Copy(os.Stdout, cmdStdout)
		go io.Copy(os.Stderr, cmdStderr)

		ticker := time.NewTicker(time.Second)
		go func(ticker *time.Ticker) {
			seconds := 0
			for _ = range ticker.C {
				if seconds++; seconds >= 5 {
					announced = discovery.AnnounceService(
						opts.DiscoveryServer,
						opts.SCInstanceIdentifier,
						opts.ServiceName,
						opts.ServiceInstanceId,
						opts.ServiceExposedPort,
					)
					break
				}
			}
		}(ticker)
	}

	err = cmd.Start()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", opts.ServiceDependencies)

	err = cmd.Wait()

	if err != nil {
		out, _ := ioutil.ReadAll(cmdStdout)
		println(string(out))
		out, _ = ioutil.ReadAll(cmdStderr)
		println(string(out))
		panic(err)
	}

	if announced == false {
		println("Service was not announced! Maybe the process exited too fast.")
		out, _ := ioutil.ReadAll(cmdStdout)
		println(string(out))
		out, _ = ioutil.ReadAll(cmdStderr)
		println(string(out))
		os.Exit(1)
	}

	<- make(chan int)
	//Receive signals and relay it to the main process
}
