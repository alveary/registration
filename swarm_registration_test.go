package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-martini/martini"
	"github.com/mechanoid/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

//
type RegistrationTestSuite struct {
	suite.Suite
	app *martini.ClassicMartini
}

//
func (suite *RegistrationTestSuite) SetupTest() {
	suite.app = AppEngine()
}

func (suite *RegistrationTestSuite) TestIfFormContainsRegistrationFields() {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	suite.app.ServeHTTP(w, req)
	doc, err := goquery.NewDocumentFromReader(w.Body)

	if err != nil {
		log.Fatal(err)
	}

	assert.True(suite.T(), doc.Exists("[name=email]"), "Should have a field for email")
	assert.True(suite.T(), doc.Exists("[name=password]"), "Should have a field for password")
}

func (suite *RegistrationTestSuite) TestCreatingANewRegistration() {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`email=falk.hoppe@innoq.com`))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	w := httptest.NewRecorder()

	suite.app.ServeHTTP(w, req)
	doc, err := goquery.NewDocumentFromReader(w.Body)

	if err != nil {
		log.Fatal(err)
	}

	assert.False(suite.T(), doc.Find(".test-registration-form [name=email]").Parent().HasClass("error"), "Shouldn't have an error, because email was given")
	assert.True(suite.T(), doc.Find(".test-registration-form [name=password]").Parent().HasClass("error"), "Should have an error, because no password given")
}

func (suite *RegistrationTestSuite) TestRegistrationWithInvalidMailAndPassword() {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`email=falk.hoppe@innoq&password=123`))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	w := httptest.NewRecorder()

	suite.app.ServeHTTP(w, req)
	doc, err := goquery.NewDocumentFromReader(w.Body)

	if err != nil {
		log.Fatal(err)
	}

	assert.True(suite.T(), doc.Find(".test-registration-form [name=email]").Parent().HasClass("error"), "Shouldn't have an error, because email is corrupt")
	assert.True(suite.T(), doc.Find(".test-registration-form [name=password]").Parent().HasClass("error"), "Should have an error, because password is too short")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRegistrationTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrationTestSuite))
}
