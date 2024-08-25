// Package brainit a customizable rune interpreter.
//
// Source code and other details for the project are available at GitHub:
//
//	https://github.com/rytsh/brainit
package brainit

import (
	"bufio"
	"errors"
	"io"
	"math/big"

	"github.com/rytsh/casset"
)

// ErrExit should throw for exiting loop.
var ErrExit = errors.New("brainit: exit loop")

var errSkip = errors.New("brainit: not possible case")

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
	begin      casset.IElement[rune]
	end        casset.IElement[rune]
	stackLoop  []*stackLoop
	stackUpper *stackLoop
}

// Interpreter is main struct hold all memory, code, executors and stacks.
type Interpreter struct {
	memory       casset.IMemory[rune]
	recCode      casset.IMemory[rune]
	currentMem   casset.IElement[rune]
	currentRec   casset.IElement[rune]
	stackCurrent *stackLoop
	executor     map[rune]Exec
	loopKeys     []LoopKey
}

// Init is initialize a Interpreter.
func (i *Interpreter) Init() *Interpreter {
	i.memory = casset.NewMemory[rune]()
	i.recCode = casset.NewMemory[rune]()

	i.executor = make(map[rune]Exec)
	i.currentMem = i.memory.GetFront()
	i.currentRec = i.recCode.GetFront()

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
	return i.currentMem.GetValue()
}

// SetValue is setting a value in a current memory.
func (i *Interpreter) SetValue(v rune) {
	i.currentMem.SetValue(v)
}

func (i *Interpreter) exec(key rune) error {
	return i.executor[key](i)
}

func (i *Interpreter) addLoop() *stackLoop {
	loop := new(stackLoop)
	loop.stackLoop = make([]*stackLoop, 0)
	loop.begin = i.currentRec

	if i.stackCurrent != nil {
		loop.stackUpper = i.stackCurrent
		i.stackCurrent.stackLoop = append(i.stackCurrent.stackLoop, loop)
	}

	return loop
}

// Next is move to next memory area.
func (i *Interpreter) Next() {
	i.currentMem = i.currentMem.Next(rune(0))
}

// Prev is move to previous memory area.
func (i *Interpreter) Prev() {
	i.currentMem = i.currentMem.Prev(rune(0))
}

// ClearMemory erase all mem and set zero of current pointer.
func (i *Interpreter) ClearMemory() {
	i.memory.Clear()
	i.currentMem = i.memory.GetFront()
}

// clearCodeMemory erase record code memory. Usable in internal.
func (i *Interpreter) clearCodeMemory() {
	lastValue := i.recCode.GetBack().GetValue()
	i.recCode.Clear().GetFront().SetValue(lastValue)
	i.currentRec = i.recCode.GetFront()
}

func (i *Interpreter) checkBegin(key rune) bool {
	for _, v := range i.loopKeys {
		if v.Begin == key {
			return true
		}
	}

	return false
}

func (i *Interpreter) checkEnd(key rune) bool {
	for _, v := range i.loopKeys {
		if v.End == key {
			return true
		}
	}

	return false
}

// checkLoop to initialize or change stack.
func (i *Interpreter) changeStack(key rune) (*stackLoop, error) {
	var changedStack *stackLoop

	if i.checkBegin(key) {
		// check end to get info not in a loop
		if i.stackCurrent == nil || i.stackCurrent.end == nil {
			// add a new loop and go inside
			changedStack = i.addLoop()
		} else if i.stackCurrent.begin != i.currentRec {
			// find own stack
			// for nested loops
			for _, v := range i.stackCurrent.stackLoop {
				if v.begin == i.currentRec {
					changedStack = v
					break
				}
			}
		}
	} else if i.checkEnd(key) {
		if i.stackCurrent != nil {
			// add end info to current stack
			if i.stackCurrent.end == nil {
				i.stackCurrent.end = i.currentRec
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
			i.currentRec = i.recCode.GetBack().Next(key)
		} else {
			i.currentRec.SetValue(key)
			if i.recCode.GetLen().Cmp(big.NewInt(1)) == 1 {
				i.clearCodeMemory()
			}
		}

		// inner loop for stacks
		for {
			key = i.currentRec.GetValue()

			if changeStack, err = i.changeStack(key); err != nil {
				break
			}

			if i.checkBegin(key) {
				i.stackCurrent = changeStack
			}

			if gotoLoopEnd {
				if i.checkEnd(key) {
					i.stackCurrent = changeStack
				}

				if loopStack == i.stackCurrent && i.checkEnd(key) {
					gotoLoopEnd = false
					loopStack = nil
				}

				break
			}

			// run function
			if err = i.exec(key); err != nil {
				if err == ErrExit {
					if i.checkBegin(key) {
						if i.stackCurrent == nil || i.stackCurrent.end == nil {
							// go to end of stack
							gotoLoopEnd = true
							loopStack = i.stackCurrent.stackUpper
							break
						}
						// inside of stuck
						// jump to end
						i.currentRec = i.stackCurrent.end
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
			if err == nil && i.checkEnd(key) {
				// loop again
				i.currentRec = i.stackCurrent.begin.GetNextElement()
				continue
			}

			i.currentRec = i.currentRec.GetNextElement()
		}
	}

	return nil
}

// NewInterpreter is helps to initialize new interpreter and return it.
func NewInterpreter() *Interpreter {
	i := new(Interpreter).Init()
	return i
}
