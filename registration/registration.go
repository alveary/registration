package registration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alveary/overseer/service"
	"github.com/jamieomatthews/validation"
	"github.com/martini-contrib/binding"
)

// Registration is the minimal information a user has to provide
type Registration struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Validate password of a given registration
func (registration Registration) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := validation.NewValidation(&errors, registration)

	v.Validate(&registration.Email).Classify("email-class").Email()
	v.Validate(&registration.Password).Key("password").MinLength(9)

	return *v.Errors.(*binding.Errors)
}

func ForService(serviceName string) (retrieved service.Service, err error) {
	overseerRoot := os.Getenv("OVERSEER_ROOT")
	responsechan := make(chan *http.Response)
	errorchan := make(chan error)
	defer func() {
		close(responsechan)
		close(errorchan)
	}()

	go func() {
		resp, err := http.Get(overseerRoot + "/" + serviceName)
		if err != nil {
			errorchan <- err
			return
		}
		if resp.StatusCode > 299 {
			errorchan <- fmt.Errorf("Request Error: %s", resp.Status)
			return
		}

		responsechan <- resp
	}()

	select {
	case resp := <-responsechan:
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)
		dec.Decode(&retrieved)

		return retrieved, nil
	case err = <-errorchan:
		return retrieved, err
	case <-time.After(time.Second * 3):
		return retrieved, fmt.Errorf("Timeout of service registry call")
	}
}

// RequestRegistration ...
func (registration Registration) RequestRegistration() (target string, err error) {
	failure := make(chan error)
	success := make(chan string)
	done := make(chan bool)

	defer func() {
		close(failure)
		close(success)
		close(done)
	}()

	go func(done chan bool) {
		json, err := json.Marshal(registration)
		if err != nil {
			failure <- err
			return
		}

		factory, err := ForService("user-factory")

		if err != nil {
			failure <- err
			return
		}

		// TODO: handle target
		_, requestErr := http.Post(factory.Root, "application/json", bytes.NewBuffer(json))

		select {
		case <-done:
			return
		default:
			if requestErr != nil {
				failure <- requestErr
				return
			}

			success <- "target"
		}
	}(done)

	select {
	case result := <-success:
		return result, nil
	case exception := <-failure:
		return "", exception
	case <-time.After(time.Second * 3):
		done <- true
		return "", errors.New("timeout")
	}
}
