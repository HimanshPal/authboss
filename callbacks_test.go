package authboss

import (
	"errors"
	"testing"
)

func TestCallbacks(t *testing.T) {
	afterCalled := false
	beforeCalled := false
	c := NewCallbacks()

	c.Before(EventRegister, func(ctx *Context) (bool, error) {
		beforeCalled = true
		return false, nil
	})
	c.After(EventRegister, func(ctx *Context) error {
		afterCalled = true
		return nil
	})

	if beforeCalled || afterCalled {
		t.Error("Neither should be called.")
	}

	stopped, err := c.FireBefore(EventRegister, NewContext())
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if stopped {
		t.Error("It should not have been stopped.")
	}

	if !beforeCalled {
		t.Error("Expected before to have been called.")
	}
	if afterCalled {
		t.Error("Expected after not to be called.")
	}

	c.FireAfter(EventRegister, NewContext())
	if !afterCalled {
		t.Error("Expected after to be called.")
	}
}

func TestCallbacksInterrupt(t *testing.T) {
	before1 := false
	before2 := false
	c := NewCallbacks()

	c.Before(EventRegister, func(ctx *Context) (bool, error) {
		before1 = true
		return true, nil
	})
	c.Before(EventRegister, func(ctx *Context) (bool, error) {
		before2 = true
		return false, nil
	})

	stopped, err := c.FireBefore(EventRegister, NewContext())
	if err != nil {
		t.Error(err)
	}
	if !stopped {
		t.Error("It was not stopped.")
	}

	if !before1 {
		t.Error("Before1 should have been called.")
	}
	if before2 {
		t.Error("Before2 should not have been called.")
	}
}

func TestCallbacksErrors(t *testing.T) {
	before1 := false
	before2 := false
	c := NewCallbacks()

	errValue := errors.New("Problem occured.")

	c.Before(EventRegister, func(ctx *Context) (bool, error) {
		before1 = true
		return false, errValue
	})
	c.Before(EventRegister, func(ctx *Context) (bool, error) {
		before2 = true
		return false, nil
	})

	stopped, err := c.FireBefore(EventRegister, NewContext())
	if err != errValue {
		t.Error("Expected an error to come back.")
	}
	if stopped {
		t.Error("It should not have been stopped.")
	}

	if !before1 {
		t.Error("Before1 should have been called.")
	}
	if before2 {
		t.Error("Before2 should not have been called.")
	}
}
