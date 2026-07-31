package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)

	lru := NewLRUCache()

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "CAP") {

			str := strings.TrimSpace(line[4:])
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println(err)
			}

			lru.setCapacity(num)

		} else if strings.HasPrefix(line, "PUT") {

			strs := strings.Split(line[4:], " ")

			num, err := strconv.Atoi(strs[1])
			if err != nil {
				fmt.Println(err)
			}

			ops := lru.put(strs[0], num)
			fmt.Printf("ops=%d\n", ops)

		} else if strings.HasPrefix(line, "GET") {

			str := strings.TrimSpace(line[4:])

			node, ops := lru.get(str)
			fmt.Printf("value=%d ops=%d", node.value, ops)
			fmt.Println()
		} else if strings.HasPrefix(line, "STATE") {
			lru.state()
		}
	}
}

type Node struct {
	prev, next *Node
	key        string
	value      int
}

func NewNode() *Node {
	return &Node{
		prev:  nil,
		next:  nil,
		key:   "",
		value: 0,
	}
}

type DLL struct {
	head, tail *Node
}

func NewDLL() *DLL {
	head := NewNode()
	tail := NewNode()
	head.next = tail
	tail.prev = head
	return &DLL{
		head: head,
		tail: tail,
	}
}

func (d *DLL) addFront(n *Node) {
	n.prev = d.head
	n.next = d.head.next
	d.head.next.prev = n
	d.head.next = n
}

func (d *DLL) removeKey(n *Node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (d *DLL) moveFront(n *Node) {
	d.removeKey(n)
	d.addFront(n)
}

func (d *DLL) findNodeByKey(key string) (*Node, error) {
	for n := d.head; n != nil; n = n.next {
		if n.key == key {
			return n, nil
		}
	}
	return nil, fmt.Errorf("key not found")
}

func (d *DLL) list() {
	for n := d.head.next; n != d.tail; n = n.next {
		fmt.Printf("%s=%d ", n.key, n.value)
	}
	fmt.Println()
}

type LRUCache struct {
	cap   int
	dll   *DLL
	cache map[string]*Node
}

func NewLRUCache() *LRUCache {
	return &LRUCache{
		cap:   0,
		dll:   NewDLL(),
		cache: make(map[string]*Node),
	}
}

func (c *LRUCache) setCapacity(cap int) {
	c.cap = cap
}

func (c *LRUCache) get(key string) (*Node, int) {
	ops := 0
	if node, ok := c.cache[key]; ok {
		ops++

		c.dll.removeKey(node)
		ops++

		c.dll.moveFront(node)
		ops++
		return node, ops
	}

	return nil, 0
}

func (c *LRUCache) put(key string, value int) int {
	ops := 0
	if c.cap == 0 {
		return 0
	}

	if node, ok := c.cache[key]; ok {
		ops++

		node.value = value
		c.dll.removeKey(node)
		ops++

		c.dll.addFront(node)
		ops++

		c.cache[key] = node
		ops += 2
		return ops
	}

	if len(c.cache) >= c.cap {
		node := c.dll.tail.prev
		c.dll.removeKey(c.dll.tail.prev)
		ops++

		delete(c.cache, node.key)
		ops++
	}

	newNode := NewNode()
	newNode.key = key
	newNode.value = value
	c.dll.addFront(newNode)
	ops++

	c.cache[key] = newNode
	ops++
	return ops
}

func (c *LRUCache) state() {
	c.dll.list()
}
