package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "

const MONKEY_FACE = `
 .--.  .-"      "-.  .--.
/ .. \/  .-. .-.   \/ .. \
| | '|  /   Y   \  |' |  |
| \  \  \ 0 | 0 /  /  /  |
\ '- ,\.-"""""""-./, -' /
 ''-' /_   ^ ^   _\ '-''
     |  \._   _./  |
     \   \ '~' /   /
      '._ '-=-' _.'
         '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	inChan := make(chan string)
	outChan := make(chan string)

	go StartChannel(inChan, outChan)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		inChan <- scanner.Text()
		output := <-outChan
		io.WriteString(out, output)
	}
}

func StartChannel(in chan string, out chan string) {
	env := object.NewEnvironment()

	for {
		line := <-in

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			out <- printParseErrors(p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			out <- evaluated.Inspect() + "\n"
		} else {
			out <- ""
		}
	}
}

func printParseErrors(errors []string) string {
	var out bytes.Buffer
	out.WriteString(MONKEY_FACE)
	out.WriteString("Woops! We ran into some monkey business here!\n")
	out.WriteString("parser errors:\n")
	for _, msg := range errors {
		out.WriteString("\t" + msg + "\n")
	}
	return out.String()
}
