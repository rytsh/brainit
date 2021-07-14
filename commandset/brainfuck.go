// Package commandset is include predefined interpreter commands.
package commandset

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rytsh/brainit"
)

// Brainfuck is totaly brainfuck command set for brainit.
var Brainfuck = brainit.Cset{
	LoopCommands: map[brainit.LoopKey]brainit.Exec{
		{
			Begin: '[',
			End:   ']',
		}: func(i *brainit.Interpreter) error {
			if i.GetValue() == 0 {
				return brainit.ErrExit
			}
			return nil
		},
	},
	Commands: map[rune]brainit.Exec{
		'>': func(i *brainit.Interpreter) error {
			i.Next()
			return nil
		},
		'<': func(i *brainit.Interpreter) error {
			i.Prev()
			return nil
		},
		'+': func(i *brainit.Interpreter) error {
			i.SetValue(i.GetValue() + 1)
			return nil
		},
		'-': func(i *brainit.Interpreter) error {
			i.SetValue(i.GetValue() - 1)
			return nil
		},
		'.': func(i *brainit.Interpreter) error {
			fmt.Printf("%c", i.GetValue())
			return nil
		},
		',': func(i *brainit.Interpreter) error {
			in := bufio.NewReader(os.Stdin)
			r, _, _ := in.ReadRune()
			i.SetValue(r)
			return nil
		},
	},
}
