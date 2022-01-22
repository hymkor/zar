package stringstack

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

func (node *_Node) push(b []byte) *_Node {
	length := len(b)

	if node.top < length+2 {
		return newNode(node).push(b)
	}

	newTop := node.top - length
	copy(node.buffer[newTop:node.top], b)

	newTop -= 2
	node.buffer[newTop+0] = byte(length & 0xFF)
	node.buffer[newTop+1] = byte(length >> 8)

	node.top = newTop
	return node
}

func (node *_Node) pop() (*_Node, []byte) {
	length := int(node.buffer[node.top]) + (int(node.buffer[node.top+1]) << 8)
	node.top += 2
	b := node.buffer[node.top : node.top+length]
	node.top += length

	if node.top >= len(node.buffer) {
		return node.next, b
	}
	return node, b
}

type Stack struct {
	first *_Node
}

func (stack *Stack) Push(b []byte) {
	if stack.first == nil {
		stack.first = newNode(nil)
	}
	stack.first = stack.first.push(b)
}

func (stack *Stack) PushString(s string) {
	stack.Push([]byte(s))
}

func (stack *Stack) Pop() []byte {
	if stack.first == nil {
		return nil
	}
	var b []byte
	stack.first, b = stack.first.pop()
	return b
}

func (stack *Stack) PopString() (string, bool) {
	b := stack.Pop()
	if b == nil {
		return "", false
	}
	return string(b), true
}
