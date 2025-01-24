package streamtagparser

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// TagParser Non-concurrency safe, a TagParser can only be used for one stream
type TagParser struct {
	needParsed []string

	tagTotalBuffer   strings.Builder
	tagAttrBuffer    strings.Builder
	tagContentBuffer strings.Builder
	tagEndBuffer     strings.Builder

	inTag        bool
	inTagName    bool
	inAttr       bool
	inTagContent bool

	currentTagName string
}

func NewTagParser(needParsed ...string) *TagParser {
	return &TagParser{
		needParsed: needParsed,
	}
}

func (p *TagParser) initStatus() {
	p.inTag = false
	p.inTagName = false
	p.inAttr = false
	p.inTagContent = false

	p.currentTagName = ""

	p.tagTotalBuffer.Reset()
	p.tagAttrBuffer.Reset()
	p.tagContentBuffer.Reset()
	p.tagEndBuffer.Reset()
}

func (p *TagParser) Parse(streamStr string) (tagsData []*TagStreamData) {
	if streamStr == "" {
		return nil
	}

	tagStr := streamStr

	if !p.inTag {
		i := strings.Index(streamStr, "<")
		if i == -1 {
			tagsData = append(tagsData, &TagStreamData{
				Type: TagStreamTypeText,
				Text: streamStr,
			})
			return
		} else {
			if rawText := streamStr[:i]; rawText != "" {
				tagsData = append(tagsData, &TagStreamData{
					Type: TagStreamTypeText,
					Text: rawText,
				})
			}
			tagStr = streamStr[i:]
		}
	}

	for _, r := range tagStr {
		tagsData = append(tagsData, p.parseRune(r)...)
	}

	return p.mergeStreams(tagsData)
}

func (p *TagParser) ParseDone() (tagsData []*TagStreamData) {
	if (p.inTagName || p.inAttr) && p.tagTotalBuffer.Len() > 0 {
		tagsData = append(tagsData, NewTextTagStreamData(p.tagTotalBuffer.String()))
	}
	if p.inTagContent {
		content := p.tagContentBuffer.String()
		if p.tagEndBuffer.Len() > 0 {
			content += p.tagEndBuffer.String()
		}
		tagsData = append(tagsData, NewEndTagStreamData(p.currentTagName, p.parseAttr(), content))
	}
	p.initStatus()
	return
}

func (p *TagParser) mergeStreams(ss []*TagStreamData) (list []*TagStreamData) {
	for _, s := range ss {
		if s == nil {
			continue
		}
		var last *TagStreamData
		if len(list) > 0 {
			last = list[len(list)-1]
		}
		if last == nil {
			list = append(list, s)
			continue
		}
		if last.Type == TagStreamTypeText && s.Type == last.Type {
			last.Text += s.Text
			continue
		}
		if last.Type == TagStreamTypeContent && s.Type == last.Type {
			last.Content += s.Content
			continue
		}
		list = append(list, s)
	}
	return
}

func (p *TagParser) parseRune(r rune) (tags []*TagStreamData) {
	switch {
	case r == '<' && !p.inTag:
		p.inTag = true
		p.inTagName = true
		if p.tagTotalBuffer.Len() > 0 {
			tags = append(tags, NewTextTagStreamData(p.tagTotalBuffer.String()))
		}
		p.tagTotalBuffer.Reset()
		p.tagTotalBuffer.WriteRune(r)
		return
	case p.inTagName:
		if r == '<' {
			tags = append(tags, NewTextTagStreamData(p.tagTotalBuffer.String()))
			p.tagTotalBuffer.Reset()
			p.tagTotalBuffer.WriteRune(r)
			return
		}
		if r == ' ' {
			p.inTagName = false
			p.inAttr = true
			p.parseCurrentTagName()
		}
		if r == '>' {
			p.inTagContent = true
			p.inTagName = false
			p.parseCurrentTagName()
			p.tagTotalBuffer.WriteRune(r)
			if !p.isNeedParseTag() {
				tags = append(tags, NewTextTagStreamData(p.tagTotalBuffer.String()))
				p.initStatus()
			} else {
				tags = append(tags, NewStartTagStreamData(p.currentTagName, nil))
			}
			return
		}
		if !p.tagPrefixMatch() {
			tags = append(tags, NewTextTagStreamData(p.tagTotalBuffer.String()+string(r)))
			p.initStatus()
			return
		}
		p.tagTotalBuffer.WriteRune(r)
		return
	case p.inAttr:
		if r == '>' {
			p.inTagContent = true
			p.inAttr = false
			p.tagTotalBuffer.WriteRune(r)
			tags = append(tags, NewStartTagStreamData(p.currentTagName, p.parseAttr()))
			return
		}
		p.tagAttrBuffer.WriteRune(r)
		p.tagTotalBuffer.WriteRune(r)
		// 防止ai输出错误
		if p.tagAttrBuffer.Len() > 500 {
			tags = append(tags, NewTextTagStreamData(p.tagTotalBuffer.String()))
			p.initStatus()
			return
		}
		return
	case p.inTagContent:
		var content string
		if r == '<' {
			if p.tagEndBuffer.Len() == 0 {
				p.tagEndBuffer.WriteRune(r)
			} else {
				content = p.tagEndBuffer.String()
				p.tagEndBuffer.Reset()
				p.tagEndBuffer.WriteRune(r)
			}
		}
		if p.tagEndBuffer.Len() > 0 && content == "" {
			if r != '<' {
				p.tagEndBuffer.WriteRune(r)
			}
			if p.currentTagParseEnd() {
				tags = append(
					tags,
					NewEndTagStreamData(
						p.currentTagName,
						p.parseAttr(),
						p.tagContentBuffer.String(),
					),
				)
				p.initStatus()
				return
			}
			if p.tagSuffixMatch() {
				p.tagTotalBuffer.WriteRune(r)
				return
			}
			if !p.tagSuffixMatch() {
				content = p.tagEndBuffer.String()
				p.tagEndBuffer.Reset()
			}
		}

		p.tagTotalBuffer.WriteRune(r)
		p.tagContentBuffer.WriteString(cmp.Or(content, string(r)))
		tags = append(tags, NewContentTagStreamData(p.currentTagName, cmp.Or(content, string(r))))
		return
	default:
		tags = append(tags, NewTextTagStreamData(string(r)))
		return
	}
}

func (p *TagParser) tagPrefixMatch() bool {
	for _, tag := range p.needParsed {
		if strings.HasPrefix("<"+tag, p.tagTotalBuffer.String()) {
			return true
		}
	}
	return false
}

func (p *TagParser) tagSuffixMatch() bool {
	return strings.HasPrefix("</"+p.currentTagName, p.tagEndBuffer.String())
}

func (p *TagParser) parseCurrentTagName() {
	if p.currentTagName != "" {
		return
	}
	if p.tagTotalBuffer.Len() <= 1 {
		return
	}
	p.currentTagName = p.tagTotalBuffer.String()[1:]
}

func (p *TagParser) isNeedParseTag() bool {
	return slices.Contains(p.needParsed, p.currentTagName)
}

func (p *TagParser) currentTagParseEnd() bool {
	return fmt.Sprintf("</%s>", p.currentTagName) == p.tagEndBuffer.String()
}

func (p *TagParser) parseAttr() (tags []TagAttr) {
	if p.tagAttrBuffer.Len() == 0 {
		return nil
	}
	attrs := strings.Split(p.tagAttrBuffer.String(), " ")
	for _, attr := range attrs {
		kv := strings.Split(attr, "=")
		if len(kv) != 2 {
			continue
		}
		tags = append(tags, TagAttr{
			Name:  strings.TrimSuffix(strings.TrimPrefix(kv[0], `"`), `"`),
			Value: strings.TrimSuffix(strings.TrimPrefix(kv[1], `"`), `"`),
		})
	}
	return
}
