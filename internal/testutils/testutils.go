package testutils

import (
	"os"
	"testing"
	"time"
)

func Date1(t *testing.T) time.Time {
	out, err := time.Parse("2006-01-02", "1993-11-23")
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func Setenv(t *testing.T, key, newValue string) (cleanup func()) {
	oldValue, existed := os.LookupEnv(key)
	err := os.Setenv(key, newValue)
	if err != nil {
		t.Fatal(err)
	}
	return func() {
		var err error
		if existed {
			err = os.Setenv(key, oldValue)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Unsetenv(t *testing.T, key string) (cleanup func()) {
	oldValue, existed := os.LookupEnv(key)
	if !existed {
		return func() {}
	}
	err := os.Unsetenv(key)
	if err != nil {
		t.Fatal(err)
	}
	return func() {
		err := os.Setenv(key, oldValue)
		if err != nil {
			t.Fatal(err)
		}
	}
}
