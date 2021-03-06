package sitter_test

import (
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/stretchr/testify/assert"
)

func TestRootNode(t *testing.T) {
	assert := assert.New(t)

	n, close := sitter.Parse([]byte("let a = 1"), javascript.GetLanguage())
	defer close()

	assert.Equal(uint32(0), n.StartByte())
	assert.Equal(uint32(9), n.EndByte())
	assert.Equal("(program (lexical_declaration (variable_declarator (identifier) (number))))", n.String())
	assert.Equal("program", n.Type())

	assert.Equal(true, n.IsNamed())
	assert.Equal(false, n.IsMissing())
	assert.Equal(false, n.HasChanges())
	assert.Equal(false, n.HasError())

	assert.Equal(uint32(1), n.ChildCount())
	assert.Equal(uint32(1), n.NamedChildCount())

	assert.Nil(n.Parent())
	assert.Nil(n.NextSibling())
	assert.Nil(n.NextNamedSibling())
	assert.Nil(n.PrevSibling())
	assert.Nil(n.PrevNamedSibling())

	assert.NotNil(n.Child(0))
	assert.NotNil(n.NamedChild(0))
}

func TestTree(t *testing.T) {
	assert := assert.New(t)

	parser := sitter.NewParser()
	defer parser.Delete()

	parser.Debug()
	parser.SetLanguage(javascript.GetLanguage())
	tree := parser.Parse([]byte("let a = 1"))
	defer tree.Delete()
	n := tree.RootNode()

	assert.Equal(uint32(0), n.StartByte())
	assert.Equal(uint32(9), n.EndByte())
	assert.Equal("program", n.Type())
	assert.Equal("(program (lexical_declaration (variable_declarator (identifier) (number))))", n.String())

	tree2 := parser.Parse([]byte("let a = 'a'"))
	defer tree2.Delete()
	n = tree2.RootNode()
	assert.Equal("(program (lexical_declaration (variable_declarator (identifier) (string))))", n.String())

	// change 'a' -> true
	newText := []byte("let a = true")
	tree2.Edit(sitter.EditInput{
		StartIndex:  8,
		OldEndIndex: 11,
		NewEndIndex: 12,
		StartPosition: sitter.Position{
			Row:    0,
			Column: 8,
		},
		OldEndPosition: sitter.Position{
			Row:    0,
			Column: 11,
		},
		NewEndPosition: sitter.Position{
			Row:    0,
			Column: 12,
		},
	})
	// check that it changed tree
	assert.True(n.HasChanges())
	assert.True(n.Child(0).HasChanges())
	assert.False(n.Child(0).Child(0).HasChanges()) // left side of tree didn't change
	assert.True(n.Child(0).Child(1).HasChanges())

	tree3 := parser.ParseWithTree(newText, tree2)
	n = tree3.RootNode()
	assert.Equal("(program (lexical_declaration (variable_declarator (identifier) (true))))", n.String())
}
