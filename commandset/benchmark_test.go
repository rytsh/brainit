package commandset_test

import (
	"strings"
	"testing"

	"github.com/rytsh/brainit"
	"github.com/rytsh/brainit/commandset"
)

func BenchmarkPiNumber(b *testing.B) {
	piProgram := `>  +++++ +++++ +++++ (15 digits)

		[<+>>>>>>>>++++++++++<<<<<<<-]>+++++[<+++++++++>-]+>>>>>>+[<<+++[>>[-<]<[>]<-]>>
		[>+>]<[<]>]>[[->>>>+<<<<]>>>+++>-]<[<<<<]<<<<<<<<+[->>>>>>>>>>>>[<+[->>>>+<<<<]>
		>>>>]<<<<[>>>>>[<<<<+>>>>-]<<<<<-[<<++++++++++>>-]>>>[<<[<+<<+>>>-]<[>+<-]<++<<+
		>>>>>>-]<<[-]<<-<[->>+<-[>>>]>[[<+>-]>+>>]<<<<<]>[-]>+<<<-[>>+<<-]<]<<<<+>>>>>>>
		>[-]>[<<<+>>>-]<<++++++++++<[->>+<-[>>>]>[[<+>-]>+>>]<<<<<]>[-]>+>[<<+<+>>>-]<<<
		<+<+>>[-[-[-[-[-[-[-[-[-<->[-<+<->>]]]]]]]]]]<[+++++[<<<++++++++<++++++++>>>>-]<
		<<<+<->>>>[>+<<<+++++++++<->>>-]<<<<<[>>+<<-]+<[->-<]>[>>.<<<<[+.[-]]>>-]>[>>.<<
		-]>[-]>[-]>>>[>>[<<<<<<<<+>>>>>>>>-]<<-]]>>[-]<<<[-]<<<<<<<<]++++++++++`

	myInterpreter := brainit.NewInterpreter()
	myInterpreter.AddCommandSet(commandset.Brainfuck)
	// Add dummy output command
	myInterpreter.AddCommand('.', func(i *brainit.Interpreter) error {
		_ = i.GetValue()
		return nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myInterpreter.ClearMemory()
		reader := strings.NewReader(piProgram)
		myInterpreter.Interpret(reader)
	}
}

func BenchmarkHelloWorld(b *testing.B) {
	helloProgram := `++++++++++[>+>+++>+++++++>++++++++++<<<<-]>>>++.>+.+++++++..+++.<<++.>+++++++++++++++.>.+++.------.--------.<<+.<`

	myInterpreter := brainit.NewInterpreter()
	myInterpreter.AddCommandSet(commandset.Brainfuck)
	myInterpreter.AddCommand('.', func(i *brainit.Interpreter) error {
		_ = i.GetValue()
		return nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myInterpreter.ClearMemory()
		reader := strings.NewReader(helloProgram)
		myInterpreter.Interpret(reader)
	}
}

func BenchmarkMultiply(b *testing.B) {
	multiplyProgram := `++.[>++++<-]>.`

	myInterpreter := brainit.NewInterpreter()
	myInterpreter.AddCommandSet(commandset.Brainfuck)
	myInterpreter.AddCommand('.', func(i *brainit.Interpreter) error {
		_ = i.GetValue()
		return nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myInterpreter.ClearMemory()
		reader := strings.NewReader(multiplyProgram)
		myInterpreter.Interpret(reader)
	}
}
