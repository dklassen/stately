package stately

import "fmt"

type StateMachine struct {
	InitialState string
	States       map[string]*State
	Events       map[string]*Event
}

type State struct {
	Name string
}

type Event struct {
	Name        string
	Transitions []*Transition
}

type Transition struct {
	to    string
	froms []string
	do    func(value interface{}) error
}

type Stately struct {
	State string
}

func NewStateMachine(initialState string) StateMachine {
	state := &State{Name: initialState}
	return StateMachine{InitialState: initialState,
		States: map[string]*State{initialState: state},
		Events: map[string]*Event{}}
}

func (s *Stately) GetState() string {
	return s.State
}

func (s *Stately) SetState(name string) {
	s.State = name
}

type Statelier interface {
	GetState() string
	SetState(name string)
}

func (sm *StateMachine) DefineState(stateName string) *State {
	state := &State{Name: stateName}
	sm.States[state.Name] = state
	return state
}

func (sm *StateMachine) DefineEvent(eventName string) *Event {
	event := &Event{Name: eventName}
	sm.Events[event.Name] = event
	return event
}

func (e *Event) To(stateName string) *Transition {
	transition := &Transition{to: stateName}
	e.Transitions = append(e.Transitions, transition)
	return transition
}

func (t *Transition) From(states ...string) *Transition {
	t.froms = states
	return t
}

func (t *Transition) Do(do func(value interface{}) error) {
	t.do = do
}

func (t *Transition) ValidFrom(fromState string) bool {
	for _, from := range t.froms {
		if from == fromState {
			return true
		}
	}
	return false
}

func (e *Event) filterTransitions(toState, fromState string) []*Transition {
	validTransitions := []*Transition{}
	for _, transition := range e.Transitions {
		matchedTransition := transition.ValidFrom(fromState)
		if matchedTransition {
			validTransitions = append(validTransitions, transition)
		}
	}
	return validTransitions
}

func (sm *StateMachine) Trigger(toState string, target Statelier) error {
	fromState := target.GetState()

	if fromState == "" {
		fromState = sm.InitialState
		target.SetState(sm.InitialState)
	}

	event, ok := sm.Events[toState]
	if !ok {
		return fmt.Errorf("Invalid next state: %s", toState)
	}

	validTransitions := event.filterTransitions(toState, fromState)
	if len(validTransitions) == 0 {
		return fmt.Errorf("No valid transition found from current state: %s", fromState)
	}

	if len(validTransitions) > 1 {
		return fmt.Errorf("More than one transition found for current state %s to next state %s", fromState, toState)
	}

	transition := validTransitions[0]
	if err := transition.do(target); err != nil {
		return err
	}
	target.SetState(transition.to)

	return nil
}
