package htmlformat

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"log"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func String(htmlStr string) string {
	buf := bytes.NewBufferString("")
	Output(htmlStr, buf)
	return buf.String()
}

func Output(htmlStr string, w io.Writer) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Fatal(err)
	}

	formatHTML(doc, w, 0)
}

func formatHTML(n *html.Node, w io.Writer, level int) {
	indent := strings.Repeat(" ", level*2)

	for c := range iterateNonEmptyChildren(n) {
		inline := isInline(c)

		if c.Type == html.ElementNode {
			if isInline(c.PrevSibling) && !isInline(c) {
				fmt.Fprint(w, "\n")
			}
			if !isInline(c.PrevSibling) || !isInline(c) {
				fmt.Fprint(w, indent)
			}

			fmt.Fprintf(w, "<%s", c.Data)
			for _, attr := range c.Attr {
				fmt.Fprintf(w, ` %s="%s"`, attr.Key, attr.Val)
			}

			fmt.Fprint(w, ">")
			if !inline {
				fmt.Fprint(w, "\n")
			}

		} else if c.Type == html.TextNode {
			if isInline(n) {
				fmt.Fprintf(w, c.Data)
			} else {
				leftSpace := unicode.IsSpace(rune(c.Data[0]))
				rightSpace := unicode.IsSpace(rune(c.Data[len(c.Data)-1]))

				if leftSpace && isInline(c.PrevSibling) {
					fmt.Fprint(w, " ")
				}

				for line := range getLines(strings.TrimSpace(c.Data)) {
					if !isInline(c.PrevSibling) || !isInline(c) {
						fmt.Fprint(w, indent)
					}
					if len(line) > 0 {
						newline := line[len(line)-1] == '\n'

						s := strings.TrimSpace(line)
						fmt.Fprint(w, s)
						if newline {
							fmt.Fprint(w, "\n")
						}
					}
				}

				if rightSpace && isInline(c.NextSibling) {
					fmt.Fprint(w, " ")
				}
			}
		} else if c.Type == html.CommentNode {
			if isInline(c) {
				fmt.Fprintf(w, "<!--%s-->", c.Data)
			} else {
				if isInline(c.PrevSibling) {
					fmt.Fprint(w, "\n")
				}
				fmt.Fprintf(w, "%s<!--%s-->\n", indent, c.Data)
			}
		} else if c.Type == html.DoctypeNode {
			fmt.Fprintf(w, "<!DOCTYPE %s", c.Data)
			for _, attr := range c.Attr {
				if attr.Key == "public" {
					fmt.Fprint(w, " PUBLIC")
				}
				fmt.Fprintf(w, ` "%v"`, attr.Val)
			}
			fmt.Fprint(w, ">\n")
		}

		formatHTML(c, w, level+1)

		if c.Type == html.ElementNode && !isVoid(c) {
			if inline {
				fmt.Fprintf(w, "</%s>", c.Data)
			} else {
				if isInline(c.LastChild) {
					fmt.Fprint(w, "\n")
				}
				fmt.Fprintf(w, "%s</%s>\n", indent, strings.TrimSpace(c.Data))
			}
		}
	}

}

// skips text nodes with only whitespaces
func iterateNonEmptyChildren(node *html.Node) iter.Seq[*html.Node] {
	return func(yield func(*html.Node) bool) {
		child := node.FirstChild
		for child != nil {
			isEmptyStr := child.Type == html.TextNode && strings.TrimSpace(child.Data) == ""

			if isEmptyStr {
				next := child.NextSibling
				node.RemoveChild(child)
				child = next
				continue
			}
			yield(child)
			child = child.NextSibling
		}
	}
}
