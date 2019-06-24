package main

import "fmt"

func main() {
	t := &tree{}
	t.insert("/aa/bb/cc", func() {
		fmt.Println("aabbcc")
	})
	t.find("/aa/bb/cc")()

}

type tree struct {
	children map[string]*tree
	handler  func()
}

func (t *tree) insert(path string, handler func()) {
	low, high := 1, 1
	for high < len(path) {
		for high < len(path) && path[high] != '/' {
			high++
		}

		if t.children == nil {
			t.children = map[string]*tree{}
		}
		if t.children[path[low:high]] == nil {
			t.children[path[low:high]] = &tree{}
		}

		t = t.children[path[low:high]]
		high++
		low = high
	}

	t.handler = handler
}

func (t *tree) find(path string) func() {
	low, high := 1, 1
	for high < len(path) {
		for high < len(path) && path[high] != '/' {
			high++
		}

		if t.children == nil || t.children[path[low:high]] == nil {
			return nil
		}

		t = t.children[path[low:high]]
		high++
		low = high
	}

	return t.handler
}
