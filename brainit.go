// Package brainit a customizable rune interpreter.
//
// Source code and other details for the project are available at GitHub:
//
//   https://github.com/rytsh/brainit
//
package brainit

import (
	"bufio"
	"errors"
	"io"

	"github.com/rytsh/casset"
)

// ErrExit should throw for exiting loop.
var ErrExit = errors.New("brainit: exit loop")

var errSkip = errors.New("brainit: not possible case")

type beginEnd uint

const (
	begin beginEnd = iota
	end
)

// Exec is function signature, usable for add new ones.
type Exec func(*Interpreter) error

// LoopKey to record new new loop runes.
type LoopKey struct {
	Begin rune
	End   rune
}

// Cset for commands sets.
type Cset struct {
	LoopCommands map[LoopKey]Exec
	Commands     map[rune]Exec
}

type stackLoop struct {
	begin      casset.IElement
	end        casset.IElement
	stackLoop  []*stackLoop
	stackUpper *stackLoop
}

// Interpreter is main struct hold all memory, code, executors and stacks.
type Interpreter struct {
	memory       casset.IMemory
	recCode      casset.IMemory
	stackCurrent *stackLoop
	executor     map[rune]Exec
	loopKeys     []LoopKey
}

// Init is initialize a Interpreter.
func (i *Interpreter) Init() *Interpreter {
	i.memory = casset.NewMemory(rune(0))
	i.recCode = casset.NewMemory(rune(0))
	i.executor = make(map[rune]Exec)

	return i
}

// AddCommand for record new function to specific key.
func (i *Interpreter) AddCommand(key rune, fn Exec) {
	i.executor[key] = fn
}

// AddLoopCommand for record new function to specific loop keys.
func (i *Interpreter) AddLoopCommand(begin rune, end rune, fn Exec) {
	i.executor[begin] = fn
	i.executor[end] = fn
	i.loopKeys = append(i.loopKeys, LoopKey{
		Begin: begin,
		End:   end,
	})
}

// AddCommandSet for record new command sets.
func (i *Interpreter) AddCommandSet(c Cset) {
	for key, value := range c.Commands {
		i.AddCommand(key, value)
	}

	for key, value := range c.LoopCommands {
		i.AddLoopCommand(key.Begin, key.End, value)
	}
}

// GetValue is return current value in a current memory.
func (i *Interpreter) GetValue() rune {
	return i.memory.GetCurrent().GetValue().(rune)
}

// SetValue is setting a value in a current memory.
func (i *Interpreter) SetValue(v rune) {
	i.memory.GetCurrent().SetValue(v)
}

func (i *Interpreter) exec(key rune) error {
	return i.executor[key](i)
}

func (i *Interpreter) addLoop() *stackLoop {
	loop := new(stackLoop)
	loop.stackLoop = make([]*stackLoop, 0)
	loop.begin = i.recCode.GetCurrent()

	if i.stackCurrent != nil {
		loop.stackUpper = i.stackCurrent
		i.stackCurrent.stackLoop = append(i.stackCurrent.stackLoop, loop)
	}

	return loop
}

// Next is move to next memory area.
func (i *Interpreter) Next() {
	i.memory.SetCurrent(i.memory.GetCurrent().Next(rune(0)))
}

// Prev is move to previous memory area.
func (i *Interpreter) Prev() {
	i.memory.SetCurrent(i.memory.GetCurrent().Prev(rune(0)))
}

// ClearMemory erase all mem and set zero of current pointer.
func (i *Interpreter) ClearMemory() {
	i.memory.Remove(nil, nil)
	i.memory.Init(casset.NewElement(rune(0)))
}

// clearCodeMemory erase record code memory. Usable in internal.
func (i *Interpreter) clearCodeMemory() {
	i.recCode.Remove(nil, i.recCode.GetBack().GetPrevElement())
}

func (i *Interpreter) checkBeginEnd(key rune, check beginEnd) (res bool) {
	for _, v := range i.loopKeys {
		if check == begin && v.Begin == key {
			res = true
			break
		} else if check == end && v.End == key {
			res = true
			break
		}
	}

	return
}

// checkLoop to initialize or change stack.
func (i *Interpreter) changeStack(key rune) (*stackLoop, error) {
	var changedStack *stackLoop

	if i.checkBeginEnd(key, begin) {
		// check end to get info not in a loop
		if i.stackCurrent == nil || i.stackCurrent.end == nil {
			// add a new loop and go inside
			changedStack = i.addLoop()
		} else if i.stackCurrent.begin != i.recCode.GetCurrent() {
			// find own stack
			// for nested loops
			for _, v := range i.stackCurrent.stackLoop {
				if v.begin == i.recCode.GetCurrent() {
					changedStack = v
					break
				}
			}
		}
	} else if i.checkBeginEnd(key, end) {
		if i.stackCurrent != nil {
			// add end info to current stack
			if i.stackCurrent.end == nil {
				i.stackCurrent.end = i.recCode.GetCurrent()
			}
			changedStack = i.stackCurrent.stackUpper
		} else {
			// should not possible
			return nil, errSkip
		}
	}

	return changedStack, nil
}

// Interpret is a runner.
// if reader get error other than io.EOF, return it.
func (i *Interpreter) Interpret(r io.Reader) error {
	reader := bufio.NewReader(r)

	gotoLoopEnd := false

	var loopStack *stackLoop
	var changeStack *stackLoop

	for {
		// read a rune
		key, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// check keyword is exist
		if _, exist := i.executor[key]; !exist {
			continue
		}

		// record code
		if i.stackCurrent != nil {
			i.recCode.SetCurrent(i.recCode.GetBack().Next(key))
		} else {
			i.recCode.GetCurrent().SetValue(key)
			if i.recCode.GetLen().Cmp(1) == 1 {
				i.clearCodeMemory()
			}
		}

		// inner loop for stacks
		for {
			key = i.recCode.GetCurrent().GetValue().(rune)

			if changeStack, err = i.changeStack(key); err != nil {
				break
			}

			if i.checkBeginEnd(key, begin) {
				i.stackCurrent = changeStack
			}

			if gotoLoopEnd {
				if i.checkBeginEnd(key, end) {
					i.stackCurrent = changeStack
				}

				if loopStack == i.stackCurrent && i.checkBeginEnd(key, end) {
					gotoLoopEnd = false
					loopStack = nil
				}

				break
			}

			// run function
			if err = i.exec(key); err != nil {
				if err == ErrExit {
					if i.checkBeginEnd(key, begin) {
						if i.stackCurrent == nil || i.stackCurrent.end == nil {
							// go to end of stack
							gotoLoopEnd = true
							loopStack = i.stackCurrent.stackUpper
							break
						}
						// inside of stuck
						// jump to end
						i.recCode.SetCurrent(i.stackCurrent.end)
					}
					// end of the stuck exit
					i.stackCurrent = i.stackCurrent.stackUpper

					// upper stack and cursor is in the end of inner stuck
				}
			}

			// check not in the loop read new keys
			if i.stackCurrent == nil || i.stackCurrent.end == nil {
				break
			}

			// return loop begin
			if err == nil && i.checkBeginEnd(key, end) {
				// loop again
				i.recCode.SetCurrent(i.stackCurrent.begin.GetNextElement())
				continue
			}

			i.recCode.SetCurrent(i.recCode.GetCurrent().GetNextElement())
		}
	}

	return nil
}

// NewInterpreter is helps to initialize new interpreter and return it.
func NewInterpreter() *Interpreter {
	i := new(Interpreter).Init()
	return i
}
