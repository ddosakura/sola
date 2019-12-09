package sola

import (
	"reflect"
	"testing"
)

func Test_context_Store(t *testing.T) {
	c1 := newContext()
	c1.Set("a", 1)
	c1.Set("b", 2)
	c2 := c1.Shadow()
	c2.Set("b", 3)
	c2.Set("c", 4)
	c3 := c2.Shadow()
	c3.Set("c", 5)
	c3.Set("d", 6)
	tests := []struct {
		name    string
		context Context
		wantS   map[string]interface{}
	}{
		{
			name:    "base test",
			context: c3,
			wantS: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": 5,
				"d": 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.context
			if gotS := c.Store(); !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("context.Store() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
