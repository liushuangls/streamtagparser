package streamtagparser

import (
	"fmt"
	"testing"
)

type parserTest struct {
	items     []parserTestItem
	doneDatas []*TagStreamData
}

type parserTestItem struct {
	input        string
	expectedTags []*TagStreamData
}

var longText = `"article"="The Model Context Protocol is rapidly evolving. This page outlines our current thinking on key priorities and future direction for the first half of 2025, though these may change significantly as the project develops.The Model Context Protocol is rapidly evolving. This page outlines our current thinking on key priorities and future direction for the first half of 2025, though these may change significantly as the project develops. The Model Context Protocol is rapidly evolving. This page outlines our current thinking on key priorities and future direction for the first half of 2025, though these may change significantly as the project develops.The Model Context Protocol is rapidly evolving. This page outlines our current thinking on key priorities and future direction for the first half of 2025, though these may change significantly as the project develops. The Model Context Protocol is rapidly evolving. This page outlines our current thinking on key priorities and future direction for the first half of 2025, though these may change significantly as the project develops.The Model Context Protocol is rapidly evolving. This page outlines our current thinking on key priorities and future direction for the first half of 2025, though these may change significantly as the project develops."`

func TestTagParser(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		testData := parserTest{
			items: []parserTestItem{
				{
					input: "hello",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "hello"},
					},
				},
				{
					input: "world <A",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "world "},
					},
				},
				{
					input:        "rtifact",
					expectedTags: nil,
				},
				{
					input: ` "id"=1>`,
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeStart, TagName: "Artifact", Attrs: []TagAttr{{Name: "id", Value: "1"}}},
					},
				},
				{
					input: "local a",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "local a"},
					},
				},
				{
					input: "=1",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "=1"},
					},
				},
				{
					input:        "</Ar",
					expectedTags: nil,
				},
				{
					input: "tifact>end",
					expectedTags: []*TagStreamData{
						{
							Type:    TagStreamTypeEnd,
							TagName: "Artifact",
							Attrs:   []TagAttr{{Name: "id", Value: "1"}},
							Content: "local a=1",
						},
						{Type: TagStreamTypeText, Text: "end"},
					},
				},
			},
		}
		testParserTest(t, testData, "Artifact")
	})

	t.Run("multi tag normal", func(t *testing.T) {
		items := parserTest{
			items: []parserTestItem{
				{
					input: `hello <Artifact "id"="123">worl`,
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "hello "},
						{Type: TagStreamTypeStart, TagName: "Artifact", Attrs: []TagAttr{{Name: "id", Value: "123"}}},
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "worl"},
					},
				},
				{
					input: `d</Artifact>!`,
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "d"},
						{Type: TagStreamTypeEnd, TagName: "Artifact", Attrs: []TagAttr{{Name: "id", Value: "123"}}, Content: "world"},
						{Type: TagStreamTypeText, Text: "!"},
					},
				},
				{
					input: "your <Think>>>!",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "your "},
						{Type: TagStreamTypeStart, TagName: "Think"},
						{Type: TagStreamTypeContent, TagName: "Think", Content: ">>!"},
					},
				},
				{
					input: "&<</Think>>",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeContent, TagName: "Think", Content: "&<"},
						{Type: TagStreamTypeEnd, TagName: "Think", Content: ">>!&<"},
						{Type: TagStreamTypeText, Text: ">"},
					},
				},
			},
		}
		testParserTest(t, items, "Artifact", "Think")
	})

	t.Run("tag with incorrect end", func(t *testing.T) {
		items := parserTest{
			items: []parserTestItem{
				{
					input: `333 <Artifact "name"="wu">444</Artif`,
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "333 "},
						{Type: TagStreamTypeStart, TagName: "Artifact", Attrs: []TagAttr{{Name: "name", Value: "wu"}}},
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "444"},
					},
				},
				{
					input: `>555`,
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "</Artif>555"},
					},
				},
				{
					input: "321</Artifact>123",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "321"},
						{Type: TagStreamTypeEnd, TagName: "Artifact", Attrs: []TagAttr{{Name: "name", Value: "wu"}}, Content: "444</Artif>555321"},
						{Type: TagStreamTypeText, Text: "123"},
					},
				},
			},
		}
		testParserTest(t, items, "Artifact")
	})

	t.Run("full str", func(t *testing.T) {
		items := parserTest{
			items: []parserTestItem{
				{
					input: `hello <Artifact "id"="1" name="a">world</Artifact> !`,
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "hello "},
						{Type: TagStreamTypeStart, TagName: "Artifact", Attrs: []TagAttr{{Name: "id", Value: "1"}, {Name: "name", Value: "a"}}},
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "world"},
						{
							Type:    TagStreamTypeEnd,
							TagName: "Artifact",
							Attrs:   []TagAttr{{Name: "id", Value: "1"}, {Name: "name", Value: "a"}},
							Content: "world",
						},
						{Type: TagStreamTypeText, Text: " !"},
					},
				},
			},
		}
		testParserTest(t, items, "Artifact", "Think")
	})

	t.Run("tag with incorrect start", func(t *testing.T) {
		items := parserTest{
			items: []parserTestItem{
				{
					input: "<Artifact>hello",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeStart, TagName: "Artifact"},
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "hello"},
					},
				},
				{
					input: "<Arti>hello</Arti>",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "<Arti>hello</Arti>"},
					},
				},
			},
			doneDatas: []*TagStreamData{
				{
					Type:    TagStreamTypeEnd,
					TagName: "Artifact",
					Content: "hello<Arti>hello</Arti>",
				},
			},
		}
		testParserTest(t, items, "Artifact")
	})

	t.Run("tag not match", func(t *testing.T) {
		testData := parserTest{
			items: []parserTestItem{
				{
					input: "hello <Artif",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "hello "},
					},
				},
				{
					input: ">end",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "<Artif>end"},
					},
				},
				{
					input: "</Artif>",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "</Artif>"},
					},
				},
			},
		}
		testParserTest(t, testData, "Artifact")
	})

	t.Run("tag error", func(t *testing.T) {
		items := parserTest{
			items: []parserTestItem{
				{
					input: "<Artifact>hello",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeStart, TagName: "Artifact"},
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "hello"},
					},
				},
				{
					input: "<Artifact>hello</Artifact>",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeContent, TagName: "Artifact", Content: "<Artifact>hello"},
						{Type: TagStreamTypeEnd, TagName: "Artifact", Content: "hello<Artifact>hello"},
					},
				},
				{
					input: "<Think>123",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "<Think>123"},
					},
				},
				{
					input: "456</Think>",
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: "456</Think>"},
					},
				},
			},
		}
		testParserTest(t, items, "Artifact")
	})

	t.Run("attr too long", func(t *testing.T) {
		items := parserTest{
			items: []parserTestItem{
				{
					input:        `111 <<Artifac`,
					expectedTags: []*TagStreamData{{Type: TagStreamTypeText, Text: "111 <"}},
				},
				{
					input: `t "id"="1" "name"="a" `,
				},
				{
					input: longText,
					expectedTags: []*TagStreamData{
						{Type: TagStreamTypeText, Text: fmt.Sprintf(`<Artifact "id"="1" "name"="a" %s`, longText)},
					},
				},
			},
		}
		testParserTest(t, items, "Artifact")
	})

	t.Run("parse Done", func(t *testing.T) {
		tests := []parserTest{
			{
				items: []parserTestItem{
					{
						input:        "222 <Artifact",
						expectedTags: []*TagStreamData{{Type: TagStreamTypeText, Text: "222 "}},
					},
				},
				doneDatas: []*TagStreamData{
					{Type: TagStreamTypeText, Text: "<Artifact"},
				},
			},
			{
				items: []parserTestItem{
					{
						input:        "222 <Artifact",
						expectedTags: []*TagStreamData{{Type: TagStreamTypeText, Text: "222 "}},
					},
					{
						input: ` "id"="1" name="a"`,
					},
				},
				doneDatas: []*TagStreamData{
					{Type: TagStreamTypeText, Text: `<Artifact "id"="1" name="a"`},
				},
			},
			{
				items: []parserTestItem{
					{
						input:        "222 <Artifact",
						expectedTags: []*TagStreamData{{Type: TagStreamTypeText, Text: "222 "}},
					},
					{
						input: ` "id"="1" name="a" 1 >`,
						expectedTags: []*TagStreamData{
							{Type: TagStreamTypeStart, TagName: "Artifact", Attrs: []TagAttr{{Name: "id", Value: "1"}, {Name: "name", Value: "a"}}},
						},
					},
					{
						input: "</Ar",
					},
				},
				doneDatas: []*TagStreamData{
					{
						Type:    TagStreamTypeEnd,
						TagName: "Artifact",
						Attrs:   []TagAttr{{Name: "id", Value: "1"}, {Name: "name", Value: "a"}},
						Content: "</Ar",
					},
				},
			},
		}
		for _, testData := range tests {
			testParserTest(t, testData, "Artifact")
		}
	})

	t.Run("special case", func(t *testing.T) {
		tests := []parserTest{
			{
				items: []parserTestItem{
					{
						input:        "<>",
						expectedTags: []*TagStreamData{{Type: TagStreamTypeText, Text: "<>"}},
					},
				},
			},
			{
				items: []parserTestItem{
					{
						input: "hello <<<Artifact>>>>123",
						expectedTags: []*TagStreamData{
							{Type: TagStreamTypeText, Text: "hello <<"},
							{Type: TagStreamTypeStart, TagName: "Artifact"},
							{Type: TagStreamTypeContent, TagName: "Artifact", Content: ">>>123"},
						},
					},
					{
						input: "456<<</Artifact>>>>",
						expectedTags: []*TagStreamData{
							{Type: TagStreamTypeContent, TagName: "Artifact", Content: "456<<"},
							{
								Type:    TagStreamTypeEnd,
								TagName: "Artifact",
								Content: ">>>123456<<",
							},
							{Type: TagStreamTypeText, Text: ">>>"},
						},
					},
				},
			},
		}
		for _, testData := range tests {
			testParserTest(t, testData, "Artifact")
		}
	})
}

func testParserTest(t *testing.T, testData parserTest, tags ...string) {
	parser := NewTagParser(tags...)

	for _, item := range testData.items {
		tags := parser.Parse(item.input)
		if len(item.expectedTags) != len(tags) {
			t.Fatalf(" input: %s, expected tags length: %d, got: %d", item.input, len(item.expectedTags), len(tags))
		}
		for i, tag := range tags {
			tagEqual(t, item.expectedTags[i], tag)
		}
	}

	doneItems := parser.ParseDone()
	if len(testData.doneDatas) != len(doneItems) {
		t.Fatalf("expected done tags length: %d, got: %d", len(testData.doneDatas), len(doneItems))
	}
	for i, tag := range doneItems {
		tagEqual(t, testData.doneDatas[i], tag)
	}
}

func tagEqual(t *testing.T, expected, actual *TagStreamData) {
	if expected.Type != actual.Type {
		t.Fatalf("expected type: %s, got: %s", expected.Type, actual.Type)
	}
	if expected.Text != actual.Text {
		t.Fatalf("expected text: %s, got: %s", expected.Text, actual.Text)
	}
	if expected.TagName != actual.TagName {
		t.Fatalf("expected tag name: %s, got: %s", expected.TagName, actual.TagName)
	}
	if expected.Content != actual.Content {
		t.Fatalf("expected content: %s, got: %s", expected.Content, actual.Content)
	}
	if len(expected.Attrs) != len(actual.Attrs) {
		t.Fatalf("expected attrs length: %d, got: %d", len(expected.Attrs), len(actual.Attrs))
	}
	for i, attr := range actual.Attrs {
		if expected.Attrs[i].Name != attr.Name {
			t.Fatalf("expected attr name: %s, got: %s", expected.Attrs[i].Name, attr.Name)
		}
		if expected.Attrs[i].Value != attr.Value {
			t.Fatalf("expected attr value: %s, got: %s", expected.Attrs[i].Value, attr.Value)
		}
	}
}
