package brainit

import (
	"strings"
	"testing"
)

var i = NewInterpreter()

func init() {
	i.AddCommand('+', func(i *Interpreter) error {
		i.SetValue(i.GetValue() + 1)
		return nil
	})

	i.AddCommand('>', func(i *Interpreter) error {
		i.Next()
		return nil
	})

	i.AddCommand('<', func(i *Interpreter) error {
		i.Prev()
		return nil
	})

	i.AddLoopCommand('(', ')', func(i *Interpreter) error {
		if i.GetValue() == 0 {
			return ErrExit
		}
		return nil
	})
}

func TestInterpreter_ClearMemory(t *testing.T) {
	test := strings.NewReader("++(>++>)>")
	if err := i.Interpret(test); err != nil {
		t.Error(err)
	}

	if i.currentRec != i.recCode.GetFront() || i.currentRec != i.recCode.GetBack() {
		t.Error("recCode clear failed!")
	}

	i.ClearMemory()
	if i.currentMem != i.memory.GetFront() || i.currentMem != i.memory.GetBack() {
		t.Error("memory clear failed!")
	}
}

func TestInterpreter_SKIPSTACK(t *testing.T) {
	i.ClearMemory()
	test := strings.NewReader("(>++>)<")
	if err := i.Interpret(test); err != nil {
		t.Error(err)
	}

	if got := i.GetValue(); got != rune(0) {
		t.Error("Mistake in reader!, got: ", got)
	}

	i.ClearMemory()

	test = strings.NewReader("+++[>++>++[<+>-]++[<+>-]<<-]")
	if err := i.Interpret(test); err != nil {
		t.Error(err)
	}
	if got := i.GetValue(); got != rune(3) {
		t.Error("Mistake in reader!, got: ", got)
	}

	i.ClearMemory()
}

func TestInterpreter_NotPossible(t *testing.T) {
	i.ClearMemory()
	test := strings.NewReader("(>++>))<")
	if err := i.Interpret(test); err != nil {
		t.Error(err)
	}

	if got := i.GetValue(); got != rune(0) {
		t.Error("Mistake in reader!, got: ", got)
	}

	i.ClearMemory()
}
