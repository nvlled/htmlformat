package htmlformat

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func String(htmlStr string) string {
	buf := bytes.NewBufferString("")
	Output(htmlStr, buf)
	return buf.String()
}

func Output(htmlStr string, w io.Writer) {
	z := html.NewTokenizer(strings.NewReader(htmlStr))
	depth := 0
	pool := &nodePool{}

	var parent *html.Node = new(html.Node)
	var token *html.Token = new(html.Token)

	createNode := func(token *html.Token, ntype html.NodeType, parent *html.Node) *html.Node {
		node := pool.get()
		node.Type = ntype
		node.DataAtom = token.DataAtom
		node.Data = token.Data
		node.Attr = token.Attr
		if parent != nil {
			parent.AppendChild(node)
		}
		return node
	}

	initToken := func(tt html.TokenType, t *html.Token) *html.Token {
		t.Type = tt
		t.Attr = t.Attr[:0]
		switch tt {
		case html.TextToken, html.CommentToken, html.DoctypeToken:
			t.Data = string(z.Text())
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
		indent := strings.Repeat(" ", depth*2)

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
			node := createNode(initToken(tt, token), html.TextNode, nil)
			if !strings.ContainsFunc(node.Data, isNotSpace) {
				// skip text nodes with whitespace only
				continue
			}

			if parent != nil {
				parent.AppendChild(node)
			}

			if isInline(node.Parent) {
				fmt.Fprintf(w, node.Data)
			} else {
				leftSpace := unicode.IsSpace(rune(node.Data[0]))

				if leftSpace && isInline(node.PrevSibling) {
					fmt.Fprint(w, " ")
				}

				for line := range getLines(strings.TrimSpace(node.Data)) {
					if !isInline(node.PrevSibling) || !isInline(node) {
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
			}

		case html.SelfClosingTagToken, html.StartTagToken:
			node := createNode(initToken(tt, token), html.ElementNode, parent)

			if tt == html.StartTagToken {
				parent = node
				depth++
			}

			if isInline(node) && node.PrevSibling != nil {
				p := node.PrevSibling
				rightSpace := unicode.IsSpace(rune(p.Data[len(p.Data)-1]))
				if rightSpace {
					fmt.Fprint(w, " ")
				}
			}
			if isInline(node.PrevSibling) && !isInline(node) {
				fmt.Fprint(w, "\n")
			}
			if !isInline(node.PrevSibling) || !isInline(node) {
				fmt.Fprint(w, indent)
			}

			fmt.Fprintf(w, "<%s", node.Data)

			for _, attr := range node.Attr {
				fmt.Fprintf(w, ` %s="%s"`, attr.Key, attr.Val)
			}

			fmt.Fprint(w, ">")
			if !isInline(node) {
				fmt.Fprint(w, "\n")
			}

		case html.EndTagToken:
			node := parent
			parent = node.Parent
			if depth > 0 {
				depth--
			}
			indent := strings.Repeat(" ", depth*2)

			if !isVoid(node) {
				if isInline(node) {
					fmt.Fprintf(w, "</%s>", node.Data)
				} else {
					if isInline(node.LastChild) {
						fmt.Fprint(w, "\n")
					}
					fmt.Fprintf(w, "%s</%s>\n", indent, strings.TrimSpace(node.Data))
				}
			}
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				node.RemoveChild(c)
				pool.free(c)
			}

		case html.CommentToken, html.DoctypeToken:
			fmt.Fprintf(w, "%s%s%s", indent, string(z.Raw()), "\n")
		}
	}
}
