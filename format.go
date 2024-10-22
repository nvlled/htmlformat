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
	"golang.org/x/net/html/atom"
)

var __ = struct{}{}

var inlineSet = map[atom.Atom]struct{}{
	atom.A:        __,
	atom.Abbr:     __,
	atom.Acronym:  __,
	atom.Button:   __,
	atom.Br:       __,
	atom.Big:      __,
	atom.Bdo:      __,
	atom.B:        __,
	atom.Cite:     __,
	atom.Code:     __,
	atom.Dfn:      __,
	atom.I:        __,
	atom.Em:       __,
	atom.Img:      __,
	atom.Input:    __,
	atom.Kbd:      __,
	atom.Label:    __,
	atom.Map:      __,
	atom.Object:   __,
	atom.Output:   __,
	atom.Tt:       __,
	atom.Time:     __,
	atom.Samp:     __,
	atom.Script:   __,
	atom.Style:    __,
	atom.Select:   __,
	atom.Small:    __,
	atom.Span:     __,
	atom.Strong:   __,
	atom.Sub:      __,
	atom.Sup:      __,
	atom.Strike:   __,
	atom.Textarea: __,
}

var voidSet = map[atom.Atom]struct{}{
	atom.Area:    __,
	atom.Base:    __,
	atom.Br:      __,
	atom.Col:     __,
	atom.Command: __,
	atom.Embed:   __,
	atom.Hr:      __,
	atom.Img:     __,
	atom.Input:   __,
	atom.Keygen:  __,
	atom.Link:    __,
	atom.Meta:    __,
	atom.Param:   __,
	atom.Source:  __,
	atom.Track:   __,
	atom.Wbr:     __,
}

func isInline(node *html.Node) bool {
	if node == nil {
		return false
	}
	if node.Type == html.TextNode {
		return true
	}
	_, yes := inlineSet[node.DataAtom]
	return yes
}

func isVoid(node *html.Node) bool {
	if node == nil {
		return false
	}
	_, yes := voidSet[node.DataAtom]
	return yes
}

func String(htmlStr string) string {
	buf := bytes.NewBufferString("")
	format(htmlStr, buf)
	return buf.String()
}

func initToken(z *html.Tokenizer, tt html.TokenType, t *html.Token) *html.Token {
	t.Type = tt
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

func createNode(token *html.Token, ntype html.NodeType, parent *html.Node) *html.Node {
	node := &html.Node{}
	node.Type = ntype
	node.DataAtom = token.DataAtom
	node.Data = token.Data
	node.Attr = token.Attr
	if parent != nil {
		parent.AppendChild(node)
	}
	return node
}

func isNotSpace(r rune) bool { return !unicode.IsSpace(r) }

func format(htmlStr string, w io.Writer) {
	z := html.NewTokenizer(strings.NewReader(htmlStr))

	depth := 0
	var parent *html.Node
	var node *html.Node
	var token *html.Token = new(html.Token)

loop:
	for {
		tt := z.Next()
		indent := strings.Repeat(" ", depth*2)
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				break loop
			} else {
				panic(z.Err())
			}

		case html.TextToken:
			node = createNode(initToken(z, tt, token), html.TextNode, nil)
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
			node = createNode(initToken(z, tt, token), html.ElementNode, parent)
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
			node = parent
			parent = node.Parent
			depth--
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
			}

		case html.CommentToken, html.DoctypeToken:
			fmt.Fprintf(w, "%s%s%s", indent, string(z.Raw()), "\n")
		}
	}
}

func Write(htmlStr string, w io.Writer) {
	// this is so dumb, why isn't there an option
	// to set document root, and why does it add <html><body> by default.
	// solution 1: fork the html package so that I could set the document root
	// solution 2: use treesitter instead
	// solution 3: build partial node tree from the tokenizer
	// solution 4: option to remove <html><body> from output

	// I'll go with 3
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
			// TODO: do not trim, just find a whitespace

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
			a = b + 1
		}
	}
}
