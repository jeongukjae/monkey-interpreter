package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

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
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		out.WriteString(fmt.Sprintf("%+v\n", tok))
	}
	return out.String()
}
