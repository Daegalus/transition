# Transition

Transition is a [Golang](http://golang.org/) [*state machine*](https://en.wikipedia.org/wiki/Finite-state_machine) implementation.

NOTE: This fork removes all GORM/database backed functionality including dependencies. The goal here is to have a clean and useful FSM implementation, and nothing more.

[![GoDoc](https://godoc.org/github.com/qor/transition?status.svg)](https://godoc.org/github.com/qor/transition)

# Usage

### Enable Transition for your struct

Embed `transition.Transition` into your struct, it will enable the state machine feature for the struct:

```go
import "github.com/qor/transition"

type Order struct {
  ID uint
  transition.Transition
}
```

### Define States and Events

```go
var OrderStateMachine = transition.New(&Order{})

// Define initial state
OrderStateMachine.Initial("draft")

// Define a State
OrderStateMachine.State("checkout")

// Define another State and what to do when entering and exiting that state.
OrderStateMachine.State("paid").Enter(func(order interface{}) error {
  // To get order object use 'order.(*Order)'
  // business logic here
  return
}).Exit(func(order interface{}) error {
  // business logic here
  return
})

// Define more States
OrderStateMachine.State("cancelled")
OrderStateMachine.State("paid_cancelled")


// Define an Event
OrderStateMachine.Event("checkout").To("checkout").From("draft")

// Define another event and what to do before and after performing the transition.
OrderStateMachine.Event("paid").To("paid").From("checkout").Before(func(order interface{}) error {
  // business logic here
  return nil
}).After(func(order interface{}) error {
  // business logic here
  return nil
})

// Different state transitions for one event
cancellEvent := OrderStateMachine.Event("cancel")
cancellEvent.To("cancelled").From("draft", "checkout")
cancellEvent.To("paid_cancelled").From("paid").After(func(order interface{}) error {
  // Refund
}})
```

### Trigger an Event

```go
// func (*StateMachine) Trigger(name string, value Stater) error
OrderStatemachine.Trigger("paid", &order)

OrderStatemachine.Trigger("cancel", &order)
// order's state will be changed to cancelled if current state is "draft"
// order's state will be changed to paid_cancelled if current state is "paid"
```

### Get/Set State

```go
var order Order

// Get Current State
order.GetState()

// Set State
order.SetState("finished") // this will only update order's state
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
