package transition

import (
	"fmt"
)

// Transition is a struct, embed it in your struct to enable state machine for the struct
type Transition[T any] struct {
	State string
}

// SetState set state to Stater, just set, won't save it into database
func (transition *Transition[T]) SetState(name string) {
	transition.State = name
}

// GetState get current state from
func (transition Transition[T]) GetState() string {
	return transition.State
}

// Stater is a interface including methods `GetState`, `SetState`
type Stater[T any] interface {
	SetState(name string)
	GetState() string
}

// New initialize a new StateMachine that hold states, events definitions
func New[T any](value T) *StateMachine[T] {
	return &StateMachine[T]{
		states: map[string]*State[T]{},
		events: map[string]*Event[T]{},
	}
}

// StateMachine a struct that hold states, events definitions
type StateMachine[T any] struct {
	initialState string
	states       map[string]*State[T]
	events       map[string]*Event[T]
}

// Initial define the initial state
func (sm *StateMachine[T]) Initial(name string) *StateMachine[T] {
	sm.initialState = name
	return sm
}

// State define a state
func (sm *StateMachine[T]) State(name string) *State[T] {
	state := &State[T]{Name: name}
	sm.states[name] = state
	return state
}

// Event define an event
func (sm *StateMachine[T]) Event(name string) *Event[T] {
	event := &Event[T]{Name: name}
	sm.events[name] = event
	return event
}

// Trigger trigger an event
func (sm *StateMachine[T]) Trigger(name string, value Stater[T]) error {
	stateWas := value.GetState()

	if stateWas == "" {
		stateWas = sm.initialState
		value.SetState(sm.initialState)
	}

	if event := sm.events[name]; event != nil {
		var matchedTransitions []*EventTransition[T]
		for _, transition := range event.transitions {
			var validFrom = len(transition.froms) == 0
			if len(transition.froms) > 0 {
				for _, from := range transition.froms {
					if from == stateWas {
						validFrom = true
					}
				}
			}

			if validFrom {
				matchedTransitions = append(matchedTransitions, transition)
			}
		}

		if len(matchedTransitions) == 1 {
			transition := matchedTransitions[0]

			// State: exit
			if state, ok := sm.states[stateWas]; ok {
				for _, exit := range state.exits {
					if err := exit(value); err != nil {
						return err
					}
				}
			}

			// Transition: before
			for _, before := range transition.befores {
				if err := before(value); err != nil {
					return err
				}
			}

			value.SetState(transition.to)

			// State: enter
			if state, ok := sm.states[transition.to]; ok {
				for _, enter := range state.enters {
					if err := enter(value); err != nil {
						value.SetState(stateWas)
						return err
					}
				}
			}

			// Transition: after
			for _, after := range transition.afters {
				if err := after(value); err != nil {
					value.SetState(stateWas)
					return err
				}
			}

			return nil
		}
	}
	return fmt.Errorf("failed to perform event %s from state %s", name, stateWas)
}

// State contains State information, including enter, exit hooks
type State[T any] struct {
	Name   string
	enters []func(value Stater[T]) error
	exits  []func(value Stater[T]) error
}

// Enter register an enter hook for State
func (state *State[T]) Enter(fc func(value Stater[T]) error) *State[T] {
	state.enters = append(state.enters, fc)
	return state
}

// Exit register an exit hook for State
func (state *State[T]) Exit(fc func(value Stater[T]) error) *State[T] {
	state.exits = append(state.exits, fc)
	return state
}

// Event contains Event information, including transition hooks
type Event[T any] struct {
	Name        string
	transitions []*EventTransition[T]
}

// To define EventTransition of go to a state
func (event *Event[T]) To(name string) *EventTransition[T] {
	transition := &EventTransition[T]{to: name}
	event.transitions = append(event.transitions, transition)
	return transition
}

// EventTransition hold event's to/froms states, also including befores, afters hooks
type EventTransition[T any] struct {
	to      string
	froms   []string
	befores []func(value Stater[T]) error
	afters  []func(value Stater[T]) error
}

// From used to define from states
func (transition *EventTransition[T]) From(states ...string) *EventTransition[T] {
	transition.froms = states
	return transition
}

// Before register before hooks
func (transition *EventTransition[T]) Before(fc func(value Stater[T]) error) *EventTransition[T] {
	transition.befores = append(transition.befores, fc)
	return transition
}

// After register after hooks
func (transition *EventTransition[T]) After(fc func(value Stater[T]) error) *EventTransition[T] {
	transition.afters = append(transition.afters, fc)
	return transition
}
