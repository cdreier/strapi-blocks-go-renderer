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

type ParagraphRenderer interface {
	RenderParagraph(Block) string
}
type TextRenderer interface {
	RenderText(Block) string
}
type ListRenderer interface {
	RenderList(Block) string
}
type ListItemRenderer interface {
	RenderListItem(Block) string
}
type HeadingRenderer interface {
	RenderHeading(Block) string
}
type LinkRenderer interface {
	RenderLink(Block) string
}
type ImageRenderer interface {
	RenderImage(Block) string
}
type QuoteRenderer interface {
	RenderQuote(Block) string
}
type CodeRenderer interface {
	RenderCode(Block) string
}

type Renderer struct {
	ParagraphRenderer ParagraphRenderer
	TextRenderer      TextRenderer
	ListRenderer      ListRenderer
	ListItemRenderer  ListItemRenderer
	HeadingRenderer   HeadingRenderer
	LinkRenderer      LinkRenderer
	ImageRenderer     ImageRenderer
	QuoteRenderer     QuoteRenderer
	CodeRenderer      CodeRenderer
}

func New() *Renderer {
	r := &Renderer{}
	r.ParagraphRenderer = r
	r.TextRenderer = r
	r.ListRenderer = r
	r.ListItemRenderer = r
	r.HeadingRenderer = r
	r.LinkRenderer = r
	r.ImageRenderer = r
	r.QuoteRenderer = r
	r.CodeRenderer = r

	return r
}

func (r *Renderer) Render(blocks []Block) string {
	out := r.internalRender(blocks)
	return gohtml.Format(out)
}

func Render(blocks []Block) string {
	r := New()
	return r.Render(blocks)
}

func (r *Renderer) internalRender(blocks []Block) string {
	out := strings.Builder{}
	for _, block := range blocks {
		out.WriteString(r.renderBlock(block))
	}
	return out.String()
}

func (r *Renderer) renderBlock(b Block) string {
	switch b.Type {
	case BlockTypeParagraph:
		return r.ParagraphRenderer.RenderParagraph(b)
	case BlockTypeText:
		return r.TextRenderer.RenderText(b)
	case BlockTypeList:
		return r.ListRenderer.RenderList(b)
	case BlockTypeListItem:
		return r.ListItemRenderer.RenderListItem(b)
	case BlockTypeHeading:
		return r.HeadingRenderer.RenderHeading(b)
	case BlockTypeLink:
		return r.LinkRenderer.RenderLink(b)
	case BlockTypeImage:
		return r.ImageRenderer.RenderImage(b)
	case BlockTypeQuote:
		return r.QuoteRenderer.RenderQuote(b)
	case BlockTypeCode:
		return r.CodeRenderer.RenderCode(b)
	}
	return "unsupported block type"
}

func (b Block) EmptyText() bool {
	return b.Type == BlockTypeText && (b.Text == nil || (b.Text != nil && *b.Text == ""))
}

func (r *Renderer) RenderParagraph(b Block) string {
	if len(b.Children) == 1 && b.Children[0].EmptyText() {
		return "<br />"
	}
	return fmt.Sprintf("<p>%s</p>", r.internalRender(b.Children))
}

func (r *Renderer) RenderText(b Block) string {
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

func (r *Renderer) RenderList(b Block) string {
	if b.Format != nil && *b.Format == string(ListFormatUnordered) {
		return fmt.Sprintf("<ul>%s</ul>", r.internalRender(b.Children))
	}
	if b.Format != nil && *b.Format == string(ListFormatOrdered) {
		return fmt.Sprintf("<ol>%s</ol>", r.internalRender(b.Children))
	}
	return "unsupported list"
}
func (r *Renderer) RenderListItem(b Block) string {
	return fmt.Sprintf("<li>%s</li>", r.internalRender(b.Children))
}
func (r *Renderer) RenderHeading(b Block) string {
	if b.Level == nil {
		return *b.Text
	}
	switch *b.Level {
	case 1:
		return fmt.Sprintf("<h1>%s</h1>", r.internalRender(b.Children))
	case 2:
		return fmt.Sprintf("<h2>%s</h2>", r.internalRender(b.Children))
	case 3:
		return fmt.Sprintf("<h3>%s</h3>", r.internalRender(b.Children))
	case 4:
		return fmt.Sprintf("<h4>%s</h4>", r.internalRender(b.Children))
	case 5:
		return fmt.Sprintf("<h5>%s</h5>", r.internalRender(b.Children))
	case 6:
		return fmt.Sprintf("<h6>%s</h6>", r.internalRender(b.Children))
	}

	return *b.Text
}

func (r *Renderer) RenderImage(b Block) string {
	if b.Image == nil {
		return "missing image"
	}
	return fmt.Sprintf("<img src=%q alt=%q />", b.Image.URL, b.Image.AlternativeText)
}

func (r *Renderer) RenderCode(b Block) string {
	// TODO: there is a "language" attribute - react renderer also ignore it
	return fmt.Sprintf("<pre><code>%s</code></pre>", r.internalRender(b.Children))
}

func (r *Renderer) RenderQuote(b Block) string {
	return fmt.Sprintf("<blockquote>%s</blockquote>", r.internalRender(b.Children))
}

func (r *Renderer) RenderLink(b Block) string {
	url := "#"
	if b.URL != nil {
		url = *b.URL
	}

	return fmt.Sprintf(`<a href=%q>%s</a>`, url, r.internalRender(b.Children))
}
