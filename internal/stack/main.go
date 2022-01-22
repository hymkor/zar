package stringstack

import (
	"io"
	"strings"
)

type _Node struct {
	next   *_Node
	top    int
	buffer [4096]byte
}

func newNode(next *_Node) *_Node {
	return &_Node{
		top:  len(next.buffer),
		next: next,
	}
}

func (node *_Node) push(s string) *_Node {
	length := len(s)

	if node.top < length+2 {
		return newNode(node).push(s)
	}

	newTop := node.top - length
	for i := 0; i < length; i++ {
		node.buffer[newTop+i] = s[i]
	}

	newTop -= 2
	node.buffer[newTop+0] = byte(length & 0xFF)
	node.buffer[newTop+1] = byte(length >> 8)

	node.top = newTop
	return node
}

func (node *_Node) pop(buffer io.Writer) *_Node {
	length := int(node.buffer[node.top]) + (int(node.buffer[node.top+1]) << 8)
	node.top += 2
	buffer.Write(node.buffer[node.top : node.top+length])
	node.top += length

	if node.top >= len(node.buffer) {
		return node.next
	}
	return node
}

type Stack struct {
	first *_Node
}

func (stack *Stack) Push(s string) {
	if stack.first == nil {
		stack.first = newNode(nil)
	}
	stack.first = stack.first.push(s)
}

func (stack *Stack) PopTo(buffer io.Writer) bool {
	if stack.first == nil {
		return false
	}
	stack.first = stack.first.pop(buffer)
	return true
}

func (stack *Stack) Pop() (string, bool) {
	var buffer strings.Builder
	ok := stack.PopTo(&buffer)
	if !ok {
		return "", false
	}
	return buffer.String(), true
}
