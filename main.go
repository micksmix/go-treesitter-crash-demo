package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

func main() {
	path := flag.String("path", ".", "The path to the directory to scan")
	flag.Parse()

	err := filepath.Walk(*path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			processFile(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through files:", err)
	}
}

func processFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %s, error: %v\n", path, err)
		return
	}

	parser := sitter.NewParser()
	defer parser.Close()

	lang := cpp.GetLanguage()
	parser.SetLanguage(lang)

	ctx := context.Background()
	tree, err := parser.ParseCtx(ctx, nil, content)
	if err != nil {
		fmt.Printf("Error parsing file: %s, error: %v\n", path, err)
		return
	}
	defer tree.Close()

	queryString := `(declaration 
		declarator: (init_declarator
			declarator: (identifier) @key
			value: (string_literal) @val
		)
	)`

	query, err := sitter.NewQuery([]byte(queryString), lang)
	if err != nil {
		fmt.Printf("Error creating query: %v\n", err)
		return
	}
	defer query.Close()

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()

	cursor.Exec(query, tree.RootNode())

	fmt.Printf("Results for %s:\n", path)
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			fmt.Printf("Captured text: %s\n", content[capture.Node.StartByte():capture.Node.EndByte()])
		}
	}
}
