package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nrtkbb/go-MEL/lexer"
	"github.com/nrtkbb/go-MEL/token"
)

const prompt = ">> "

// Start は in を解析して out に書き出すループを実行する
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
