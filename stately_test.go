package stately

import (
	"reflect"
	"sort"
	"testing"
)

type TestStater struct {
	Value int
	Stately
}

func TestDefineState(t *testing.T) {
	sm := NewStateMachine("the-beginning")
	sm.DefineState("the-end")

	if len(sm.States) != 2 {
		t.Fatal("Expected 2 statss to be defined")
	}

	result := make([]string, 0, len(sm.States))
	for k := range sm.States {
		result = append(result, k)
	}

	expected := []string{"the-beginning", "the-end"}

	sort.Strings(result)
	sort.Strings(expected)
	if !reflect.DeepEqual(result, expected) {
		t.Fatal("Expected last state to be 'the-end'")
	}
}

func TestDefineEventAndTransition(t *testing.T) {
	input := &TestStater{Value: 1}

	sm := NewStateMachine("start")
	sm.DefineState("end")
	sm.DefineEvent("go-to-end").
		To("end").
		From("start").
		Do(func(x interface{}) error {
		v := x.(*TestStater)
		v.Value = v.Value * 2
		return nil
	})

	err := sm.Trigger("go-to-end", input)

	result := input.GetState()
	if err != nil {
		t.Fatalf("Wasn't expecting an error: %s", err)
	}

	if result != "end" {
		t.Errorf("Expectd input to be in state 'end' not %s", result)
	}
}
