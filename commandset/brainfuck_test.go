package commandset_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rytsh/brainit"
	"github.com/rytsh/brainit/commandset"
)

type TResult uint

const (
	NUM TResult = iota
	PP
)

type rest struct {
	num string
	pp  string
}

type test struct {
	str    string
	result string
	typ    TResult
	name   string
}

var result = rest{}

func testManipulation(myInterpreter *brainit.Interpreter) {
	myInterpreter.AddCommand('.', func(i *brainit.Interpreter) error {
		result.num += fmt.Sprintf("%v", i.GetValue())
		result.pp += fmt.Sprintf("%c", i.GetValue())
		return nil
	})
}

func TestInterpreter_RUNTIME(t *testing.T) {
	tests := []test{
		{
			name:   "Small test",
			str:    `+++.[---].>+++[---.]++.`,
			result: `3002`,
			typ:    NUM,
		},
		{
			name:   "Inner loop test",
			str:    `+++.[>++>++[<+>-.].++[<+>-]<<-.]`,
			result: `3100210011000`,
			typ:    NUM,
		},
		{
			name:   "Basic sub/plus test",
			str:    `+.+..-.`,
			result: `1221`,
			typ:    NUM,
		},
		{
			name:   "Multiply test",
			str:    `++.[>++++<-]>.`,
			result: `28`,
			typ:    NUM,
		},
		{
			name:   "Inner stack",
			str:    `++.[>++++<--[+++++]]>.`,
			result: `24`,
			typ:    NUM,
		},
		{
			name:   "Multiple inner stack",
			str:    `++.[>++++<--[+++++]>[--.]]>.`,
			result: `2200`,
			typ:    NUM,
		},
		{
			name: "Pi number",
			str: `>  +++++ +++++ +++++ (15 digits)
	
				[<+>>>>>>>>++++++++++<<<<<<<-]>+++++[<+++++++++>-]+>>>>>>+[<<+++[>>[-<]<[>]<-]>>
				[>+>]<[<]>]>[[->>>>+<<<<]>>>+++>-]<[<<<<]<<<<<<<<+[->>>>>>>>>>>>[<+[->>>>+<<<<]>
				>>>>]<<<<[>>>>>[<<<<+>>>>-]<<<<<-[<<++++++++++>>-]>>>[<<[<+<<+>>>-]<[>+<-]<++<<+
				>>>>>>-]<<[-]<<-<[->>+<-[>>>]>[[<+>-]>+>>]<<<<<]>[-]>+<<<-[>>+<<-]<]<<<<+>>>>>>>
				>[-]>[<<<+>>>-]<<++++++++++<[->>+<-[>>>]>[[<+>-]>+>>]<<<<<]>[-]>+>[<<+<+>>>-]<<<
				<+<+>>[-[-[-[-[-[-[-[-[-<->[-<+<->>]]]]]]]]]]<[+++++[<<<++++++++<++++++++>>>>-]<
				<<<+<->>>>[>+<<<+++++++++<->>>-]<<<<<[>>+<<-]+<[->-<]>[>>.<<<<[+.[-]]>>-]>[>>.<<
				-]>[-]>[-]>>>[>>[<<<<<<<<+>>>>>>>>-]<<-]]>>[-]<<<[-]<<<<<<<<]++++++++++
				`,
			result: `3.14159265358979`,
			typ:    PP,
		},
		{
			name:   "Hello World",
			str:    `++++++++++[>+>+++>+++++++>++++++++++<<<<-]>>>++.>+.+++++++..+++.<<++.>+++++++++++++++.>.+++.------.--------.<<+.<`,
			result: `Hello World!`,
			typ:    PP,
		},
	}

	// get interpreter
	myInterpreter := brainit.NewInterpreter()
	// add brainfuck commands
	myInterpreter.AddCommandSet(commandset.Brainfuck)
	testManipulation(myInterpreter)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result = rest{}
			reader := strings.NewReader(tt.str)

			myInterpreter.ClearMemory()
			if err := myInterpreter.Interpret(reader); err != nil {
				t.Error(err)
			}

			resString := "Result:[%v] Wants:[%v]"
			failed := false
			if tt.typ == NUM {
				resString = fmt.Sprintf(resString, result.num, tt.result)
				if result.num != tt.result {
					failed = true
				}
			} else if tt.typ == PP {
				resString = fmt.Sprintf(resString, result.pp, tt.result)
				if strings.Compare(result.pp, tt.result) != 0 {
					fmt.Println(strings.Compare(result.pp, tt.result))
					failed = true
				}
			}

			if failed {
				t.Error("\nFailed Test - ", resString)
			}
		})
	}
}
