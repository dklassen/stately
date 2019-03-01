# Stately
This is not a repo for use. This was a learning project for Go. Inspired heavily from https://github.com/qor/transition.

A simple state machine with the barebones for transitioning states:

```golang
import (
  "fmt"
  "stately"
)

type Poll struct {
  Name string
  Question string
  Answers []string
  stately.Stately
}

// create and set the initial state to blank
sm := stately.New("blank")

sm.DefineState("questions_filled")
sm.DefineEvent("collect_questions")\
    .To("questions_filled")\
    .From("blank")
    .Do(func(obj interface{})error{
      // Sadly, You still need to do the proper type conversions
      poll := obj.(*Poll)'
      ...
      return err
     })

poll := &Poll{Name: "bananas"}
// Trigger the event to transition the poll from on state to another
if err := sm.Trigger("collect_questions", poll, "Collecting a question for poll"); err != nil {
  fmt.Printf("Successfully triggered poll building"
}

```

