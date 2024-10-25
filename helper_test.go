package htmlformat

import (
	"testing"
)

func TestCollapseLeftWhitespace(t *testing.T) {
	cases := [][]string{
		{"", ""},
		{" ", " "},
		{"\t", " "},
		{"\n", "\n"},
		{"\n\t ", "\n"},
		{"x", "x"},
		{" x", " x"},
		{"\tx", " x"},
		{"\nx", "\nx"},
		{"xy ", "xy "},
		{"  xy ", " xy "},
		{"\t xy ", " xy "},
		{"\t xy ", " xy "},
		{"  \n xy ", "\nxy "},
	}

	for _, c := range cases {
		arg, expected := c[0], c[1]
		actual := collapseLeftWhitespace(arg)
		if actual != expected {
			t.Errorf("expected=%q, got=%q", expected, actual)
		}
	}

}

func TestCollapseRightWhitespace(t *testing.T) {
	cases := [][]string{
		{"", ""},
		{" ", " "},
		{"\t", " "},
		{"\n", "\n"},
		{"\n\t ", "\n"},
		{"x", "x"},
		{" x ", " x "},
		{"x\n", "x\n"},
		{"xy ", "xy "},
		{"xy        ", "xy "},
		{"xy  \t\t  ", "xy "},
		{"\nxy\t\n\t\t  ", "\nxy\n"},
		{" xy  ", " xy "},
		{" xy\t", " xy "},
		{" xy\t ", " xy "},
		{"\nxy\n", "\nxy\n"},
	}

	for _, c := range cases {
		arg, expected := c[0], c[1]
		actual := collapseRightWhitespace(arg)
		if actual != expected {
			t.Errorf("expected=%q, got=%q", expected, actual)
		}
	}

}

func TestCollapseWhitespace(t *testing.T) {
	cases := [][]string{
		{"", ""},
		{" ", " "},
		{"\t", " "},
		{"\n", "\n"},
		{"\n\t ", "\n"},
		{"x", "x"},
		{"  x  ", " x "},
		{"\nx\n", "\nx\n"},
		{"\txy ", " xy "},
		{"\txy        ", " xy "},
		{"\txy  \t\t  ", " xy "},
		{" \t\t  xy \t\t\n", " xy\n"},
	}

	for _, c := range cases {
		arg, expected := c[0], c[1]
		actual := collapseWhitespace(arg)
		if actual != expected {
			t.Errorf("expected=%q, got=%q", expected, actual)
		}
	}

}

func TestGetLines(t *testing.T) {
	input := `
1111
2222
3333

4444
5555`

	expected := []string{
		"\n",
		"1111\n",
		"2222\n",
		"3333\n",
		"\n",
		"4444\n",
		"5555",
	}

	lineNum := 0
	for line := range getLines(input) {
		if lineNum >= len(expected) {
			t.Errorf("got more lines than expected")
			break
		}
		if expected[lineNum] != line {
			t.Errorf("expected=%q at line %d, got=%q", expected[lineNum], lineNum+1, line)
		}
		lineNum++
	}

}

func TestDedent(t *testing.T) {
	input := `
	function foo(x)
		if x % 2 == 0 then
			if blah then
				return 2
			else
				return 3
			end
			return 1
		end
		return 0
	end
	`
	expected := `
function foo(x)
    if x % 2 == 0 then
        if blah then
            return 2
        else
            return 3
        end
        return 1
    end
    return 0
end
	`
	actual := dedent(input)

	if actual != expected {
		t.Errorf("unexpectd output")
		println("--------------[ expected ]--------------")
		println(expected)
		println("--------------[  actual  ]--------------")
		println(actual)
	}
}
