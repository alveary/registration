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

	assert.True(suite.T(), doc.Exists(".test-registration-form"), "Should have a form on it rendered")
	assert.True(suite.T(), doc.Exists("[name=firstname]"), "Should have a field for firstname")
	assert.True(suite.T(), doc.Exists("[name=lastname]"), "Should have a field for firstname")
}

func (suite *RegistrationTestSuite) TestCreatingANewRegistration() {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`lastname=hoppe&email=falk.hoppe@innoq.com`))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	w := httptest.NewRecorder()

	suite.app.ServeHTTP(w, req)
	doc, err := goquery.NewDocumentFromReader(w.Body)

	if err != nil {
		log.Fatal(err)
	}

	assert.True(suite.T(), doc.Find(".test-registration-form [name=firstname]").Parent().HasClass("error"), "Should have a form on it rendered")
	assert.False(suite.T(), doc.Find(".test-registration-form [name=lastname]").Parent().HasClass("error"), "Should have a form on it rendered")
	assert.False(suite.T(), doc.Find(".test-registration-form [name=email]").Parent().HasClass("error"), "Should have a form on it rendered")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRegistrationTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrationTestSuite))
}
