// support {} pattern like
// - /api/user/{userid}/info
// - /api/user/user-{userid}/info
// not support
// - /api/user/{userid}-user/info

package main

import (
	"fmt"
)

const (
	patternStart = '{'
	patternEnd   = '}'
	separator = '/'
)

func main() {
	t := &tree{}
	t.Insert("/api/user", func() {
		fmt.Println("certain 1")
	})
	t.Insert("/api/user/123/info", func() {
		fmt.Println("certain 2")
	})
	t.Insert("/api/user-{id}/info", func() {
		fmt.Println("prefix pattern 1")
	})
	t.Insert("/api/user/user-{id}/info", func() {
		fmt.Println("prefix pattern 2")
	})
	t.Insert("/api/user/{userid}/info", func() {
		fmt.Println("pattern")
	})
	t.Insert("/api/user/456/info", func() {
		fmt.Println("certain 3")
	})
	t.Insert("/api/unit", func(){
		fmt.Println("certain 4")
	})
	print(t)
	fmt.Println()

	fmt.Println(t.Find("/api/unit/") == nil)
	fmt.Println()

	t.Find("/api/user")()
	t.Find("/api/user/123/info")()
	t.Find("/api/user-123/info")()
	t.Find("/api/user/user-1234/info")()
	t.Find("/api/user/1234/info")()
	t.Find("/api/user/456/info")()
	t.Find("/api/unit")()
}

type tree struct {
	children  []*tree
	char      byte
	isPattern bool
	handler   func()
}

func hasPattern(path string, index int) int {
	if path[index] != patternStart {
		return -1
	}
	for i := index + 1; i < len(path); i++ {
		if path[i] == patternEnd {
			return i
		}
	}
	return -1
}

func (t *tree) Insert(path string, handler func()) {
	index := 0
	for index < len(path) {
		var nextTree *tree

		if patternIndex := hasPattern(path, index); patternIndex != -1 {
			for _, child := range t.children {
				if child.isPattern {
					nextTree = child
				}
				break
			}
			if nextTree == nil {
				nextTree = &tree{isPattern: true}
				t.children = append(t.children, nextTree)
			}
			index = patternIndex + 1
		} else {
			char := path[index]
			for _, child := range t.children {
				if child.char == char {
					nextTree = child
					break
				}
			}
			if nextTree == nil {
				nextTree = &tree{char: char}
				t.children = append(t.children, nextTree)
			}
			index++
		}

		t = nextTree
	}

	t.handler = handler
}

func (t *tree) Find(path string) func() {
	return t.find(path, -1)
}

func (t *tree) find(path string, index int) func() {
	if t.isPattern {
		for index+1 < len(path) && path[index+1] != separator {
			index++
		}
	} else {
		if index != -1 && index < len(path) {
			if path[index] != t.char {
				return nil
			}
		}
	}

	if index == len(path)-1 {
		return t.handler
	}

	for _, child := range t.children {
		if handler := child.find(path, index+1); handler != nil {
			return handler
		}
	}

	return nil
}

func print(t *tree) {
	helper(t, "")
}

func helper(t *tree, path string) {
	if t == nil {
		return
	}
	if t.isPattern {
		path += string([]byte{patternStart, patternEnd})
	} else {
		path += string(t.char)
	}

	if t.handler != nil {
		fmt.Println(path)
	}

	for _, child := range t.children {
		helper(child, path)
	}
}
