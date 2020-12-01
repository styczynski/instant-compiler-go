package parser

import (
	"bufio"
	"bytes"
	"strings"
	"time"

	"github.com/alecthomas/participle/v2"
)

func tryInsertingBracket(done chan struct{}, parser *participle.Parser, bracket string, input []byte, lineNo int, pos int, goPast bool) ([]byte, bool, int, int) {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	curLineNo := 1
	contentLines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if curLineNo == lineNo {
			if goPast {
				for i, ch := range line {
					if ch == ' ' && i >= pos {
						pos = i
						break
					}
				}
			} else {
				if pos-1 > 0 {
					pos--
				} else {
					return input, false, lineNo, pos
				}
				for i := pos; i>=0; i-- {
					if line[i] == ' ' {
						pos = i
						break
					}
				}
			}
			line = line[:pos] + bracket + line[pos:]
		}
		contentLines = append(contentLines, line)
		curLineNo++
	}

	newContent := []byte(strings.Join(contentLines, "\n"))
	err := parser.Parse("", bytes.NewReader(newContent), &LatteProgram{})
	if err == nil {
		return newContent, true, 10000, 10000
	}

	parserError := err.(participle.Error)
	//fmt.Printf(string(newContent))
	//fmt.Printf("GOT IMPROVEMENT L: %d->%d C: %d->%d\n", lineNo, parserError.Position().Line, pos, parserError.Position().Column)
	if parserError.Position().Line > lineNo || (parserError.Position().Line == lineNo && parserError.Position().Column > pos + 3) {
		t1, _, l1, c1 := tryInsertingBracket(done, parser, bracket, newContent, parserError.Position().Line, parserError.Position().Column, true)
		t2, _, l2, c2 := tryInsertingBracket(done, parser, bracket, newContent, parserError.Position().Line, parserError.Position().Column, false)
		if l1 > l2 || (l1 == l2 && c1 > c2) {
			return t1, true, l1, c1
		}
		return t2, true, l2, c2
	}
	return input, false, lineNo, pos
}


func tryInsertingBrackets(parser *participle.Parser, input []byte, lineNo int, pos int) string {
	brackets := []string{ ")", "(", "{", "}" }
	maxL := lineNo
	maxC := pos
	bestBracket := ""
	for _, bracket := range brackets {
		for _, goPast := range []bool{ true, false } {
			done := make(chan struct{})
			go func() {
				_, isOk, l, c := tryInsertingBracket(done, parser, bracket, input, lineNo, pos, goPast)
				if isOk {
					if (maxL < l) || (maxL == l && maxC < pos-3) {
						maxL = l
						maxC = c
						bestBracket = bracket
					}
				}
				close(done)
			}()

			ok := false
			select {
			case <-done:
				ok = true
			case <-time.After(5 * time.Second):
			}

			if !ok {
				return ""
			}
		}
	}

	return bestBracket
}
