package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mechanoid/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegistrationTestSuite struct {
	suite.Suite
	app *gin.Engine
}

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

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRegistrationTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrationTestSuite))
}
