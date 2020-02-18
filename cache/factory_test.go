package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	observed := New(InfiniteMapStrategy, 1*time.Second, 1024, 1024)
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf((*Cache)(nil)).Elem()
	if reflect.PtrTo(observedType).Implements(expectedType) {
		t.Error("Expected", expectedType, "got", observedType)
	}
}
