package main

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

func main() {
	input := []byte("function hello() { console.log('hello') }; function goodbye(){}")

	parser := sitter.NewParser()
	defer parser.Delete()
	parser.SetLanguage(javascript.GetLanguage())

	tree := parser.Parse(input)
	defer tree.Delete()

	n := tree.RootNode()

	fmt.Println("AST:", n)
	fmt.Println("Root type:", n.Type())
	fmt.Println("Root children:", n.ChildCount())

	fmt.Println("\nFunctions in input:")
	iter := sitter.NewIterator(n, sitter.DFSMode)
	var funcs []*sitter.Node
	iter.ForEach(func(n *sitter.Node) error {
		if n.Type() == "function" {
			fmt.Println("-", sitter.FuncName(input, n))
			funcs = append(funcs, n)
		}
		return nil
	})

	fmt.Println("\nEdit input")
	input = []byte("function hello() { console.log('hello') }; function goodbye(){ console.log('goodbye') }")
	// reuse tree
	tree.Edit(sitter.EditInput{
		StartIndex:  62,
		OldEndIndex: 63,
		NewEndIndex: 87,
		StartPosition: sitter.Position{
			Row:    0,
			Column: 62,
		},
		OldEndPosition: sitter.Position{
			Row:    0,
			Column: 63,
		},
		NewEndPosition: sitter.Position{
			Row:    0,
			Column: 87,
		},
	})

	for _, f := range funcs {
		var textChange string
		if f.HasChanges() {
			textChange = "has change"
		} else {
			textChange = "no changes"
		}
		fmt.Println("-", sitter.FuncName(input, f), ">", textChange)
	}

	newTree := parser.ParseWithTree(input, tree)
	n = newTree.RootNode()
	fmt.Println("\nNew AST:", n)
}
