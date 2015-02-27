package announce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func requestServiceAnnouncement(overseerRoot string, service []byte) {
	go func() {
		available := true

		for available {

			checkchan := make(chan bool)
			errorchan := make(chan error)
			defer func() {
				close(checkchan)
				close(errorchan)
			}()

			go func() {
				resp, err := http.Post(overseerRoot, "application/json", bytes.NewBuffer(service))
				if err != nil {
					errorchan <- err
					return
				}
				if resp.StatusCode > 299 {
					errorchan <- fmt.Errorf("Request Error: %s", resp.Status)
					return
				}

				checkchan <- true
			}() // end request routine

			select {
			case <-checkchan:
				// JUST GO ON
			case <-errorchan:
			case <-time.After(time.Second * 3):
			}

			time.Sleep(10 * time.Minute) // update registration once a minute

		} // end for loop

	}() // end loop routine
}

// NewService provides a method to attach a new Service to the overseer stack
func NewService(serviceName string, serviceRoot string, aliveResource string) {
	service, _ := json.Marshal(struct {
		serviceName   string
		serviceRoot   string
		aliveResource string
	}{
		serviceName,
		serviceRoot,
		aliveResource,
	})

	overseerRoot := os.Getenv("OVERSEER_ROOT")

	if overseerRoot == "" {
		fmt.Println("OVERSEER_ROOT is not set:")
		return
	}

	requestServiceAnnouncement(overseerRoot, service)
}
