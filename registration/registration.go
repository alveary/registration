package registration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alveary/overseer-client/ask"
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

		factory, err := ask.ForService("user-factory")

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
