package htmlformat

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func ExampleFormat() {
	html := `<p   id = "x">hello</p > `
	fmt.Println(Format(html))
	// Output:
	// <p id="x">hello</p>
}

func ExampleWrite() {
	// Write output directly to stdout
	html := `<p   id = "x">hello to stdout</p > `
	Write(html, os.Stdout)
	// Output:
	// <p id="x">hello to stdout</p>
}

func ExampleWrite_second() {
	// Write output directly to a file
	file, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())
	html := `<p   id = "x">hello to file</p > `
	Write(html, file)
	file.Sync()
	file.Close()

	bytes, err := os.ReadFile(file.Name())
	fmt.Println(string(bytes))
	// Output:
	// <p id="x">hello to file</p>
}

func TestFormat(t *testing.T) {
	test(t, Data{
		input: `
  <!DOCTYPE html>
<div id=test  >
<input   checked >
</div>
    `,
		expectedOutput: `
<!DOCTYPE html>
<div id="test">
    <input checked>
</div>
    `,
	})

	test(t, Data{
		input: `
 <h1 >heading</h1 >
 <h1>  heading with spaces </h1>
 <h1>
    multiline heading
 </h1>

	   <h2> longer subheading with surrounding spaces </h2>

	   <h2>

	   	longer subheading with surrounding newlines

	   </h2>

	   <h3>

	   	subheading with weird spacing</h3>
    `,
		expectedOutput: `
<h1>heading</h1>
<h1> heading with spaces </h1>
<h1>
    multiline heading
</h1>
<h2> longer subheading with surrounding spaces </h2>
<h2>
    longer subheading with surrounding newlines
</h2>
<h3>
    subheading with weird spacing
</h3>
    `,
	})

	/*
	 */

	test(t, Data{
		input: `
    <!--a-->
    <div><!--b--></div>

         <!--
          aaaa
            bbbb
          cccc
         -->
    `,
		expectedOutput: `
<!--a-->
<div><!--b--></div>
<!--
aaaa
  bbbb
cccc
-->
    `,
	})

	test(t, Data{
		input: `
						<body>
						<div id="site-menu-container"><ul id="site-menu"><li><a class="selected"href="/">/top/</a></li><li><a class=""href="/?feed=new">/new/</a></li><li><a class=""href="/?feed=best">/best/</a></li><li><a class=""href="/?feed=ask">/ask/</a></li><li><a class=""href="/?feed=show">/show/</a></li><li><a class=""href="/?feed=job">/job/</a></li></ul></div><div id="site-nav"><div id="site-logo"><a href="/">^</a></div><a id="site-name"href="/">sitename</a><div id="site-nav-spacing"></div></div><div id="wrapper"></div></body>
    `,
		expectedOutput: `
<body>
    <div id="site-menu-container">
        <ul id="site-menu">
            <li>
                <a class="selected" href="/">/top/</a>
            </li>
            <li>
                <a class href="/?feed=new">/new/</a>
            </li>
            <li>
                <a class href="/?feed=best">/best/</a>
            </li>
            <li>
                <a class href="/?feed=ask">/ask/</a>
            </li>
            <li>
                <a class href="/?feed=show">/show/</a>
            </li>
            <li>
                <a class href="/?feed=job">/job/</a>
            </li>
        </ul>
    </div>
    <div id="site-nav">
        <div id="site-logo">
            <a href="/">^</a>
        </div>
        <a id="site-name" href="/">sitename</a>
        <div id="site-nav-spacing"></div>
    </div>
    <div id="wrapper"></div>
</body>
    `,
	})

	test(t, Data{
		input: `
			   	<p>
		  Whitespaces around, between, in, after texts/nodes
		  should be <i>maintained</i>. The period from
		  the last sentence should not be separated with
		  a space, but this one <b>is</b>   .

			   	  <em>    This has a leading space</em>

			   	      <em>This has a trailing space   </em>

			   	  Here's a [<a>link</a>] inside a pair of brackets without spaces.

			   	  Here's one [    <a>  link </a>   ] have erratic spaces.

		  <a>
		      This sentence is a whole link.
		  </a></p>
    `,
		expectedOutput: `
<p>
    Whitespaces around, between, in, after texts/nodes
    should be <i>maintained</i>. The period from
    the last sentence should not be separated with
    a space, but this one <b>is</b> .
    <em> This has a leading space</em>
    <em>This has a trailing space </em>
    Here's a [<a>link</a>] inside a pair of brackets without spaces.

    Here's one [ <a> link </a> ] have erratic spaces.
    <a>
        This sentence is a whole link.
    </a>
</p>
    `,
	})

	test(t, Data{
		input: `
<p>
	<em><a>blah</a></em>foo
	bar <em>baz</em>
	<a><b><i>x</i></b></a>
</p>
	`,
		expectedOutput: `
<p>
    <em><a>blah</a></em>foo
    bar <em>baz</em>
    <a><b><i>x</i></b></a>
</p>
	`,
	})

	test(t, Data{
		input: `
<a><b><i>x</i></b></a><a><b><i>x</i></b></a><a><b><i>x</i></b></a><a><b><i>x</i></b></a>
	`,
		expectedOutput: `
<a><b><i>x</i></b></a><a><b><i>x</i></b></a><a><b><i>x</i></b></a><a><b><i>x</i></b></a>
	`,
	})

	test(t, Data{
		input: `
<p>
  >>>>
		&lt;script&gt;alert(&#34;hey&#34;)&lt;/script&gt;
	<script> console.log(1 > 2); </script>
</p>
	`,
		expectedOutput: `
<p>
    >>>>
    &lt;script&gt;alert(&#34;hey&#34;)&lt;/script&gt;
    <script>
        console.log(1 > 2);
    </script>
</p>
	`,
	})

	test(t, Data{
		input: `
<pre> aaa
        bbb
    ccc
</pre>
    <pre>
        <code>
        aaa
            bbb
                ccc
        </code>
    </pre>
    <pre><code> one
    two
        three</code></pre>
	`,
		expectedOutput: `
<pre> aaa
        bbb
    ccc
</pre>
<pre>
        <code>
        aaa
            bbb
                ccc
        </code>
    </pre>
<pre><code> one
    two
        three</code></pre>
	`,
	})

	/* template
	test(t, Data{
		input: `
	`,
		expectedOutput: `
	`,
		})
	*/
}

type Data struct {
	input          string
	expectedOutput string
}

func test(t *testing.T, data Data) {
	output := Format(data.input)
	output = strings.TrimSpace(output)
	expected := strings.TrimSpace(data.expectedOutput)

	if expected != output {
		printComparison(expected, output)
		t.Error("unexpected output")
	}
}

func showTrailingSpaces(s string) string {
	var buf bytes.Buffer

	runes := []rune(s)
	var i int
	for i = len(runes) - 1; i >= 0; i-- {
		if isNotSpace(runes[i]) {
			break
		}
	}

	trailing := i
	i = 0
	for ; i < len(runes); i++ {
		c := runes[i]
		if c == ' ' && i >= trailing {
			buf.WriteRune('␣')
		} else if c == '\t' {
			buf.WriteString("↦   ")
		} else {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}
func showTrailingSpacesByLine(s string) string {
	var buf bytes.Buffer
	for line := range getLines(s) {
		buf.WriteString(showTrailingSpaces(line))
	}

	return buf.String()
}

func printComparison(expected, actual string) {
	expected = showTrailingSpacesByLine(expected)
	actual = showTrailingSpacesByLine(actual)
	line := strings.Repeat("=~", 12)
	fmt.Fprintf(os.Stderr, "\n%s[expected]%s\n%s\n%s[ actual ]%s\n%s\n%s%s%s\n", line, line, expected, line, line, actual, line, line, line)
}
