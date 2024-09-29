package blocks

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed blocks_out.json
var testInput []byte

func TestBlock_Render(t *testing.T) {

	var blocks []Block
	json.Unmarshal(testInput, &blocks)

	out := Render(blocks)

	assert.Equal(t, `<p>
  this is normal text
</p>
<p>
  this is text with
  <strong>
    bold
  </strong>
  and
  <em>
    italic
  </em>
  and
  <u>
    underlined
  </u>
  or even
  <del>
    striked
  </del>
</p>
<p>
  and multiple
  <u>
    <em>
      <strong>
        modifers at once
      </strong>
    </em>
  </u>
</p>
<p>
  we also have some
  <code>
    code
  </code>
</p>
<p>
  and
  <a href="http://asdf.de">
    links
  </a>
</p>
<h1>
  now titles: header 1
</h1>
<h2>
  header 2
</h2>
<h3>
  header 3
</h3>
<img src="http://localhost:1337/uploads/cdreier_gopher_small_a32e6e2b51.jpg" alt="cdreier_gopher_small.jpg" />
<blockquote>
  this does support block quotes
</blockquote><pre><code>func andCodeBlocks() string {
  return "with multilines"
}</code></pre>
<ul>
  <li>
    this is unorderlist
  </li>
  <li>
    list 1
  </li>
  <li>
    list 2
  </li>
  <ul>
    <li>
      sublist 1
    </li>
    <li>
      sublist 2
    </li>
  </ul>
  <li>
    list 3
  </li>
</ul>
<br />
<p>
  and ordererd
</p>
<br />
<ol>
  <li>
    one
  </li>
  <li>
    two
  </li>
  <ol>
    <li>
      two.a
    </li>
  </ol>
  <li>
    three
  </li>
</ol>
<br />`, out)

}
