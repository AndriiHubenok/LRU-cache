package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO (cache-basics): implement per the lesson description.

//type LRUCache struct {
//	cap  int
//	keys []string
//	hits int
//}
//
//func (c *LRUCache) Access(key string) {
//	if c.cap == 0 {
//		return
//	}
//
//	for i, k := range c.keys {
//		if k == key {
//			c.hits++
//
//			c.keys = append(c.keys[:i], c.keys[i+1:]...)
//			c.keys = append(c.keys, key)
//			return
//		}
//	}
//
//	if len(c.keys) == c.cap {
//		c.keys = c.keys[1:]
//	}
//	c.keys = append(c.keys, key)
//}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)

	dll := NewDLL()

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "ADD-FRONT") {

			strs := strings.Split(line[10:], " ")
			num, err := strconv.Atoi(strs[1])
			if err != nil {
				fmt.Println(err)
			}

			node := NewNode()
			node.key = strs[0]
			node.value = num
			dll.addFront(node)
			fmt.Println("OK")

		} else if strings.HasPrefix(line, "REMOVE-KEY") {

			str := strings.TrimSpace(line[11:])

			node, err := dll.findNodeByKey(str)
			if err != nil {
				fmt.Println(err)
			}

			dll.removeKey(node)
			fmt.Println("OK")

		} else if strings.HasPrefix(line, "MOVE-FRONT") {

			str := strings.TrimSpace(line[11:])

			node, err := dll.findNodeByKey(str)
			if err != nil {
				fmt.Println(err)
			}

			dll.moveFront(node)
			fmt.Println("OK")
		} else if strings.HasPrefix(line, "LIST") {
			dll.list()
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

type NaiveLRU struct {
	capacity int
	keys     []string
	values   map[string]int
}

func NewNaiveLRU(capacity int) *NaiveLRU {
	return &NaiveLRU{
		capacity: capacity,
		keys:     make([]string, 0, capacity),
		values:   make(map[string]int),
	}
}

func (c *NaiveLRU) put(key string, value int) {

	if c.values[key] != 0 {
		c.values[key] = value
		c.updateKeys(key)
		return
	}

	if len(c.keys) >= c.capacity {
		delete(c.values, c.keys[len(c.keys)-1])
		c.values[key] = value
		c.updateKeys(key)
		return
	}

	c.values[key] = value
	c.updateKeys(key)
}

func (c *NaiveLRU) get(key string) int {
	c.updateKeys(key)
	return c.values[key]
}

func (c *NaiveLRU) state() {
	for _, k := range c.keys {
		fmt.Printf("%s=%d ", k, c.values[k])
	}
}

func (c *NaiveLRU) updateKeys(key string) {
	var newKeys []string
	newKeys = append(newKeys, key)

	for _, k := range c.keys {
		if len(newKeys) >= c.capacity {
			break
		}
		if k == key {
			continue
		}
		newKeys = append(newKeys, k)
	}
	c.keys = newKeys
}
