package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
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

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		fmt.Fprint(out, StartStep(line))
	}
}

func StartStep(input string) string {
	var out bytes.Buffer

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return printParseErrors(p.Errors())
	}

	evaluated := evaluator.Eval(program)
	out.WriteString(evaluated.Inspect())
	out.WriteString("\n")
	return out.String()
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
