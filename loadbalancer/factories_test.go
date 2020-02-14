package loadbalancer

import (
	"reflect"
	"testing"
)

func TestNewTarget(t *testing.T) {
	observed, err := NewTarget("http://testhost:8000")
	if err != nil {
		t.Error(err.Error())
	}
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf((*Target)(nil)).Elem()
	if reflect.PtrTo(observedType).Implements(expectedType) {
		t.Error("Expected", expectedType, "got", observedType)
	}
}
func TestNeTargetInvalidUrl(t *testing.T) {
	_, err := NewTarget("\x80testhost:8000")
	if err == nil {
		t.Error("Url validation is broken")
		return
	}
}

func TestNew(t *testing.T) {
	observed, _ := New(RoundRobin)
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf((*TargetBucket)(nil)).Elem()
	if reflect.PtrTo(observedType).Implements(expectedType) {
		t.Error("Expected", expectedType, "got", observedType)
	}
}

func TestNewInvalidAlgorithm(t *testing.T) {
	observed, err := New("invalid")
	if err == nil {
		t.Error("Algorithm check is broken")
	}
	if observed != nil {
		t.Error("Expected nil")
	}
}
