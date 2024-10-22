package htmlformat

import (
	"testing"
)

func TestFormatSimple(t *testing.T) {
	output := String(`
	<body>
	<div>
		<!--comment----->
		<h1 class="greeting">hello</h1>
		<p>
		x<em>111</em><em>222</em>x
		</p>
		<p>
		x<em>111 </em><em>222</em>x
		</p>
		<p>
		x<em>111</em><em> 222</em>x
		</p>
		<p>
		x<em>111</em>y <em>222</em>x
		</p>
		<p>
		x<em>111</em> y<em>222</em>x
		</p>
		<p>
		x<em>111</em> y <br/> x
		</p>
	</div>
	<p>oop</p>
	</body>
	<!DOCTYPE blah>
	`)
	println("----------------")
	println(output)
}

// func TestTokenizer(t *testing.T) {
// 	huh(`
// 	<div>
// 		<h1 class="greeting">hello</h1>
// 	</div>
// 	`, nil)
// }
