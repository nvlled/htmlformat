package htmlformat

import (
	"bytes"
	"iter"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func getLines(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		a := 0
		for a < len(s) {
			b := a
			for ; b < len(s); b++ {
				if s[b] == '\n' {
					b++
					break
				}
			}
			yield(s[a:b])
			a = b
		}
	}
}

func dedent(s string) string {
	var buf bytes.Buffer
	indentSize := -1

	if strings.Count(s, "\n") <= 1 {
		return s
	}

	for line := range getLines(s) {
		if indentSize < 0 {
			numSpaces := 0
			index := 0
			for i, c := range line {
				index = i
				if !unicode.IsSpace(c) {
					indentSize = numSpaces
					break
				}
				if c == '\t' {
					numSpaces += 4
				} else {
					numSpaces += 1
				}
			}

			if indentSize < 0 {
				buf.WriteString(line)
			} else {
				buf.WriteString(line[index:])
			}

			continue
		}

		numSpaces := 0
		index := 0
		for i, c := range line {
			index = i
			if !unicode.IsSpace(c) {
				break
			}
			if c == '\t' {
				numSpaces += 4
			} else {
				numSpaces += 1
			}
		}
		indent := max(numSpaces-indentSize, 0)

		buf.WriteString(strings.Repeat(" ", indent))
		buf.WriteString(line[index:])
	}

	return buf.String()
}

func collapseLeftWhitespace(s string) string {
	newlineFound := false
	index := -1

	for i := 0; i < len(s); i++ {
		c := s[i]
		if !unicode.IsSpace(rune(c)) {
			break
		}
		newlineFound = newlineFound || c == '\n'
		index = i
	}

	if index < 0 {
		return s
	}

	space := " "
	if newlineFound {
		space = "\n"
	}

	return space + s[index+1:]
}

func collapseRightWhitespace(s string) string {
	newlineFound := false
	index := -1

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if !unicode.IsSpace(rune(c)) {
			break
		}
		newlineFound = newlineFound || c == '\n'
		index = i
	}

	if index < 0 {
		return s
	}

	space := " "
	if newlineFound {
		space = "\n"
	}

	return s[0:index] + space
}

func collapseWhitespace(s string) string {
	s = collapseLeftWhitespace(s)
	s = collapseRightWhitespace(s)
	return s
}

func isNotSpace(r rune) bool { return !unicode.IsSpace(r) }

func startsWithNewLine(node *html.Node) bool {
	if node == nil {
		return false
	}

	if node.Type == html.TextNode {
		return strings.HasPrefix(node.Data, "\n")
	}
	if node.Type == html.ElementNode {
		return startsWithNewLine(node.LastChild)
	}
	return true
}

func endsWithNewLine(node *html.Node) bool {
	if node == nil {
		return false
	}

	if node.Type == html.TextNode {
		return strings.HasSuffix(node.Data, "\n")
	}
	if node.Type == html.ElementNode {
		return endsWithNewLine(node.LastChild)
	}
	return true
}
