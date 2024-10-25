package htmlformat

import (
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

type nodePool struct {
	data []*html.Node
}

func (np *nodePool) get() *html.Node {
	if len(np.data) > 0 {
		node := np.data[0]
		np.data = np.data[1:]
		return node
	}
	return &html.Node{}
}

func (np *nodePool) free(node *html.Node) {
	if node.Parent != nil {
		node.Parent.RemoveChild(node)
	}
	np.data = append(np.data, node)
}
