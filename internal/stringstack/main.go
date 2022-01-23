package stringstack

import (
	"io"
	"strings"
)

type _Node struct {
	next   *_Node
	top    int
	buffer [4000]byte
}

func newNode(next *_Node) *_Node {
	return &_Node{
		top:  len(next.buffer),
		next: next,
	}
}

type Stack struct {
	first *_Node
}

func (stack *Stack) Push(s string) {
	if stack.first == nil {
		stack.first = newNode(nil)
	}
	length := len(s)

	if stack.first.top < length+2 {
		stack.first = newNode(stack.first)
	}
	node := stack.first

	newTop := node.top - length
	for i := 0; i < length; i++ {
		node.buffer[newTop+i] = s[i]
	}

	newTop -= 2
	node.buffer[newTop+0] = byte(length & 0xFF)
	node.buffer[newTop+1] = byte(length >> 8)

	node.top = newTop
}

func (stack *Stack) PopTo(buffer io.Writer) bool {
	if stack.first == nil {
		return false
	}
	node := stack.first

	length := int(node.buffer[node.top]) + (int(node.buffer[node.top+1]) << 8)
	node.top += 2
	buffer.Write(node.buffer[node.top : node.top+length])
	node.top += length

	if node.top >= len(node.buffer) {
		stack.first = node.next
	} else {
		stack.first = node
	}
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
