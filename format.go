package htmlformat

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Format returns a cleanly formatted html string.
func Format(htmlStr string) string {
	buf := bytes.NewBufferString("")
	Write(htmlStr, buf)
	return buf.String()
}

func Pipe(r io.Reader, w io.Writer) {
	z := html.NewTokenizer(r)
	depth := 0
	pool := &nodePool{}

	// if preDepth is greater than 1, then there's
	// at least currently one opened <pre>
	preDepth := 0

	var parent *html.Node = new(html.Node)
	var token *html.Token = new(html.Token)

	createNode := func(token *html.Token, ntype html.NodeType, parent *html.Node) *html.Node {
		node := pool.get()
		node.Type = ntype
		node.DataAtom = token.DataAtom
		node.Attr = token.Attr
		node.Data = token.Data
		if parent != nil {
			parent.AppendChild(node)
		}
		return node
	}

	initToken := func(tt html.TokenType, t *html.Token) *html.Token {
		t.Type = tt
		t.Attr = t.Attr[:0]
		switch tt {
		case html.CommentToken, html.DoctypeToken:
			t.Data = string(z.Text())
		case html.TextToken:
			// Use Raw() here instead of Text(),
			// otherwise escaped characters will be unescaped!!
			t.Data = string(z.Raw())
		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken:
			name, moreAttr := z.TagName()
			for moreAttr {
				var key, val []byte
				key, val, moreAttr = z.TagAttr()
				t.Attr = append(t.Attr, html.Attribute{
					Key: atom.String(key),
					Val: string(val),
				})
			}
			if a := atom.Lookup(name); a != 0 {
				t.DataAtom, t.Data = a, a.String()
			} else {
				t.DataAtom, t.Data = 0, string(name)
			}
		}
		return t
	}

loop:
	for {
		tt := z.Next()
		indent := strings.Repeat(" ", depth*4)

		// note: node.NextSibling will be always null
		// since this is "one-pass" iteration and
		// so the code won't know in advance what
		// the following nodes are.

		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break loop
			} else {
				panic(z.Err())
			}

		case html.TextToken:
			node := createNode(initToken(tt, token), html.TextNode, parent)
			shouldDedent := false

			if node.Parent != nil {
				switch node.Parent.Data {
				case "script", "style":
					shouldDedent = true
				}
			}

			if preDepth > 0 {
				fmt.Fprintf(w, "%s", node.Data)
			} else if shouldDedent {
				node.Data = collapseWhitespace("\n" + dedent(node.Data) + "\n")
				for line := range getLines(node.Data) {
					if strings.ContainsFunc(line, isNotSpace) {
						fmt.Fprintf(w, "%s", indent)
					}
					fmt.Fprintf(w, "%s", line)
				}
			} else {
				node.Data = collapseWhitespace(node.Data)
				lineno := 0
				for line := range getLines(node.Data) {
					if len(line) == 0 {
						continue
					}
					if lineno > 0 {
						if strings.ContainsFunc(line, isNotSpace) {
							fmt.Fprintf(w, "%s", indent)
						}
						fmt.Fprintf(w, "%s", strings.TrimLeft(line, "\t "))
					} else {
						fmt.Fprintf(w, "%s", line)
					}

					lineno++
				}
			}

		case html.SelfClosingTagToken, html.StartTagToken:
			node := createNode(initToken(tt, token), html.ElementNode, parent)

			if tt != html.SelfClosingTagToken && !isVoid(node) {
				parent = node
				depth++
			}

			if node.DataAtom == atom.Pre {
				preDepth++
			}

			if preDepth <= 0 {
				if !endsWithNewLine(node.PrevSibling) {
					if isBlock(node) || isBlock(node.PrevSibling) {
						ws := pool.get()
						ws.Type = html.TextNode
						ws.Data = "\n"
						node.Parent.InsertBefore(ws, node)
						fmt.Fprintf(w, "\n")
					}
				}

				if endsWithNewLine(node.PrevSibling) {
					fmt.Fprintf(w, "%s", indent)
				}
			}
			fmt.Fprintf(w, "<%s", node.Data)

			for _, attr := range node.Attr {
				if attr.Val == "" {
					fmt.Fprintf(w, ` %s`, attr.Key)
				} else {
					fmt.Fprintf(w, ` %s=%q`, attr.Key, attr.Val)
				}
			}
			fmt.Fprint(w, ">")

		case html.EndTagToken:
			node := parent
			parent = node.Parent
			if depth > 0 {
				depth--
			}
			indent := strings.Repeat(" ", depth*4)

			if !isVoid(node) {
				if preDepth <= 0 {
					if (startsWithNewLine(node.FirstChild) && !endsWithNewLine(node.LastChild)) ||
						isBlock(node.LastChild) {
						ws := pool.get()
						ws.Type = html.TextNode
						ws.Data = "\n"
						node.AppendChild(ws)
						fmt.Fprintf(w, "\n")
					}
					if endsWithNewLine(node.LastChild) {
						fmt.Fprintf(w, "%s", indent)
					}
				}
				fmt.Fprintf(w, "</%s>", node.Data)
			}

			if node.DataAtom == atom.Pre && preDepth > 0 {
				preDepth--
			}

			for c := node.FirstChild; c != nil; c = c.NextSibling {
				node.RemoveChild(c)
				pool.free(c)
			}

		case html.DoctypeToken:
			fmt.Fprintf(w, "%s", string(z.Raw()))

		case html.CommentToken:
			node := createNode(initToken(tt, token), html.TextNode, parent)
			node.Data = collapseWhitespace(dedent(node.Data))

			if parent != nil {
				lastChild := parent.LastChild
				if (lastChild != nil && endsWithNewLine(lastChild.PrevSibling)) || endsWithNewLine(parent) {
					fmt.Fprintf(w, "%s", indent)
				}
			}
			fmt.Fprint(w, "<!--")

			lineNum := 0
			for line := range getLines(node.Data) {
				if lineNum > 1 {
					fmt.Fprintf(w, "%s", indent)
				} else if strings.HasPrefix(node.Data, "\n") {
					fmt.Fprintf(w, "%s", indent)
				}

				fmt.Fprintf(w, "%s", line)
				lineNum++
			}

			if strings.HasSuffix(node.Data, "\n") {
				fmt.Fprintf(w, "%s", indent)
			}
			fmt.Fprint(w, "-->")
		}
	}
}

// Format writes the cleanly formatted html string into w.
func Write(htmlStr string, w io.Writer) {
	r := strings.NewReader(htmlStr)
	Pipe(r, w)
}
