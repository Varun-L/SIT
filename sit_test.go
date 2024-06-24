package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

// godogsCtxKey is the key used to store the available godogs in the context.Context.
type godogsCtxKey struct{}

func thereAreGodogs(ctx context.Context, available int) (context.Context, error) {
	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func iEat(ctx context.Context, num int) (context.Context, error) {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return ctx, errors.New("there are no godogs available")
	}

	if available < num {
		return ctx, fmt.Errorf("you cannot eat %d godogs, there are %d available", num, available)
	}

	available -= num

	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func thereShouldBeRemaining(ctx context.Context, remaining int) error {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return errors.New("there are no godogs available")
	}

	return assertExpectedAndActual(
		assert.Equal, available, remaining,
		"Expected %d godogs to be remaining, but there is %d", remaining, available,
	)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	sc.Step(`^I eat (\d+)$`, iEat)
	sc.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}

// assertExpectedAndActual is a helper function to allow the step function to call
// assertion functions where you want to compare an expected and an actual value.
func assertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, expected, actual, msgAndArgs...)
	return t.err
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

// asserter is used to be able to retrieve the error reported by the called assertion
type asserter struct {
	err error
}

// Errorf is used by the called assertion to report an error
func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}
