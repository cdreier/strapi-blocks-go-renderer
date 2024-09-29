package blocks

import (
	"fmt"
	"strings"

	"github.com/yosssi/gohtml"
)

type BlockType string

const BlockTypeParagraph BlockType = "paragraph"
const BlockTypeText BlockType = "text"
const BlockTypeList BlockType = "list"
const BlockTypeLink BlockType = "link"
const BlockTypeListItem BlockType = "list-item"
const BlockTypeHeading BlockType = "heading"
const BlockTypeImage BlockType = "image"
const BlockTypeQuote BlockType = "quote"
const BlockTypeCode BlockType = "code"

type ListFormat string

const ListFormatUnordered ListFormat = "unordered"
const ListFormatOrdered ListFormat = "ordered"

type Block struct {
	Type          BlockType `json:"type"`
	Children      []Block   `json:"children"`
	Text          *string   `json:"text"`
	Italic        *bool     `json:"italic"`
	Underline     *bool     `json:"underline"`
	Bold          *bool     `json:"bold"`
	StrikeThrough *bool     `json:"strikethrough"`
	Code          *bool     `json:"code"`
	Format        *string   `json:"format"`
	URL           *string   `json:"url"`
	Level         *int      `json:"level"`
	Image         *Image    `json:"image"`
}

type Image struct {
	Name            string `json:"name"`
	AlternativeText string `json:"alternativeText"`
	URL             string `json:"url"`
}

func Render(blocks []Block) string {
	out := internalRender(blocks)
	return gohtml.Format(out)
}

func internalRender(blocks []Block) string {
	out := strings.Builder{}
	for _, block := range blocks {
		out.WriteString(block.Render())
	}
	return out.String()
}

func (b Block) Render() string {
	switch b.Type {
	case BlockTypeParagraph:
		return b.RenderParagraph()
	case BlockTypeText:
		return b.RenderText()
	case BlockTypeList:
		return b.RenderList()
	case BlockTypeListItem:
		return b.RenderListItem()
	case BlockTypeHeading:
		return b.RenderHeading()
	case BlockTypeLink:
		return b.RenderLink()
	case BlockTypeImage:
		return b.RenderImage()
	case BlockTypeQuote:
		return b.RenderQuote()
	case BlockTypeCode:
		return b.RenderCode()
	}
	return "unsupported block type"
}

func (b Block) EmptyText() bool {
	return b.Type == BlockTypeText && (b.Text == nil || (b.Text != nil && *b.Text == ""))
}

func (b Block) RenderParagraph() string {
	if len(b.Children) == 1 && b.Children[0].EmptyText() {
		return "<br />"
	}
	return fmt.Sprintf("<p>%s</p>", internalRender(b.Children))
}

func (b Block) RenderText() string {
	out := *b.Text
	if b.Bold != nil && *b.Bold {
		out = fmt.Sprintf("<strong>%s</strong>", out)
	}
	if b.Italic != nil && *b.Italic {
		out = fmt.Sprintf("<em>%s</em>", out)
	}
	if b.Underline != nil && *b.Underline {
		out = fmt.Sprintf("<u>%s</u>", out)
	}
	if b.StrikeThrough != nil && *b.StrikeThrough {
		out = fmt.Sprintf("<del>%s</del>", out)
	}
	if b.Code != nil && *b.Code {
		out = fmt.Sprintf("<code>%s</code>", out)
	}
	return out
}

func (b Block) RenderList() string {
	if b.Format != nil && *b.Format == string(ListFormatUnordered) {
		return fmt.Sprintf("<ul>%s</ul>", internalRender(b.Children))
	}
	if b.Format != nil && *b.Format == string(ListFormatOrdered) {
		return fmt.Sprintf("<ol>%s</ol>", internalRender(b.Children))
	}
	return "unsupported list"
}
func (b Block) RenderListItem() string {
	return fmt.Sprintf("<li>%s</li>", internalRender(b.Children))
}
func (b Block) RenderHeading() string {
	if b.Level == nil {
		return *b.Text
	}
	switch *b.Level {
	case 1:
		return fmt.Sprintf("<h1>%s</h1>", internalRender(b.Children))
	case 2:
		return fmt.Sprintf("<h2>%s</h2>", internalRender(b.Children))
	case 3:
		return fmt.Sprintf("<h3>%s</h3>", internalRender(b.Children))
	case 4:
		return fmt.Sprintf("<h4>%s</h4>", internalRender(b.Children))
	case 5:
		return fmt.Sprintf("<h5>%s</h5>", internalRender(b.Children))
	case 6:
		return fmt.Sprintf("<h6>%s</h6>", internalRender(b.Children))
	}

	return *b.Text
}

func (b Block) RenderImage() string {
	if b.Image == nil {
		return "missing image"
	}
	return fmt.Sprintf("<img src=%q alt=%q />", b.Image.URL, b.Image.AlternativeText)
}

func (b Block) RenderCode() string {
	// TODO: there is a "language" attribute - react renderer also ignore it
	return fmt.Sprintf("<pre><code>%s</code></pre>", internalRender(b.Children))
}

func (b Block) RenderQuote() string {
	return fmt.Sprintf("<blockquote>%s</blockquote>", internalRender(b.Children))
}

func (b Block) RenderLink() string {
	url := "#"
	if b.URL != nil {
		url = *b.URL
	}

	return fmt.Sprintf(`<a href=%q>%s</a>`, url, internalRender(b.Children))
}
