package htmlformat

import (
	"fmt"
	"strings"
	"testing"
)

type Data struct {
	input          string
	expectedOutput string
}

func test(t *testing.T, data Data) {
	output := String(data.input)
	output = strings.TrimSpace(output)
	expected := strings.TrimSpace(data.expectedOutput)

	if expected != output {
		printComparison(expected, output)
		t.Fatalf("unexpected output")
	}
}

func printComparison(expected, actual string) {
	s := fmt.Sprintf("\n------[expected]------ \n%s\n------[ actual ]------\n%s\n----------------------\n", expected, actual)
	println(s)
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

	/*
		    template:
					test(t, Data{
						input: `
				    `,
						expectedOutput: `
				    `,
					})
	*/
}
