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

	if i.recCode.Current != i.recCode.Front || i.recCode.Current != i.recCode.Back {
		t.Error("recCode clear failed!")
	}

	i.ClearMemory()
	if i.memory.Current != i.memory.Front || i.memory.Current != i.memory.Back {
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
