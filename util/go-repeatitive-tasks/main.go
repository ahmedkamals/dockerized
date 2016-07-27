package main

import (
	"time"
	"os/signal"
	"syscall"
	"os"
	"os/exec"
	"bytes"
)

func main() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				var out bytes.Buffer
				cmd := exec.Command(
					"php",
					"/ak/projects/www/entity/live/codebase/cli/index.php",
					"--env=live",
					"--module=cli",
					"--controller=crontab",
				)
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					println("Failed to invoke the cron controller: ", err.Error())
				} else {
					//println(out.String())
				}
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	println("Triggering crontab controller every minute...")
	<-sigc
	close(quit)

	println("Bye Bye!")
}
