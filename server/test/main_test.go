package main_test

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/ory/dockertest/v3"
	"testing"
	"time"
)

var pool *dockertest.Pool
var mock *dockertest.Resource
var server *dockertest.Resource
var testContextInstance testContext

type testContext struct {
	Endpoint *string
}

func TestFeatures(t *testing.T) {
	spawnStack(t)
	time.Sleep(time.Second * 5) //Just to wait container init

	endpoint := fmt.Sprintf("http://172.17.0.1:%s", server.GetPort("8000/tcp"))
	testContextInstance.Endpoint = &endpoint

	// Initialize Godog
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			// Authentication
			s.Step(`^I have a JWT token for my (regular|admin) user$`, testContextInstance.iHaveAJWTTokenForMyUser)
			s.Step(`^I clear the current JWT token$`, testContextInstance.iClearTheCurrentJWTToken)

			// Database
			s.Step(`^I reset the database$`, testContextInstance.resetDatabase)

			// HTTP requests
			s.Step(`^I send a "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, testContextInstance.iSendARequest)
			s.Step(`^I send a "(GET|POST|PUT|DELETE)" request to "([^"]*)" with payload$`, testContextInstance.iSendARequestWithPayload)
			s.Step(`^the response code should be (\d+)$`, testContextInstance.theResponseCodeShouldBe)
			s.Step(`^the response should have (\d+) items in path "([^"]*)"$`, testContextInstance.theResponseShouldHaveItemsCountInPath)
			s.Step(`^I save the "([^"]*)" header as "([^"]*)" for suite`, testContextInstance.iSaveHeaderForSuite)
			s.Step(`^I save the value "([^"]*)" in path "([^"]*)" as "([^"]*)" for suite`, testContextInstance.iSaveTheValueInPath)
			s.Step(`^the response should have the following content in path "([^"]*)"$`, testContextInstance.theResponseShouldHaveTheFollowingContentInPath)
			s.Step(`^the response should have the following items in path "([^"]*)"$`, testContextInstance.theResponseShouldHaveTheFollowingItemsInPath)
			s.Step(`^the response should not have the key "([^"]*)" in path "([^"]*)"$`, testContextInstance.theResponseShouldNotHaveTheKeyInPath)
			s.Step(`^the response should have the key "([^"]*)" in path "([^"]*)"$`, testContextInstance.theResponseShouldHaveTheKeyInPath)

			s.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				return testContextInstance.resetDatabase(ctx)
			})
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	// If something wrong occurs, write it down there
	if suite.Run() != 0 {
		destroyStack()
		t.Fatal("non-zero status returned, failed to run feature tests")
	}

	// Tests ended, destroy our stack.
	destroyStack()
}
