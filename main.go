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

	var shardedLRU []*LRUCache
	now := 0
	evictions := 0

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "INIT") {
			now = 0
			evictions = 0

			strs := strings.Split(line[5:], " ")

			numShards, err := strconv.Atoi(strs[0])
			if err != nil {
				fmt.Println(err)
			}

			shardCap, err := strconv.Atoi(strs[1])
			if err != nil {
				fmt.Println(err)
			}

			shardedLRU = make([]*LRUCache, numShards)
			for i := 0; i < numShards; i++ {
				shardedLRU[i] = NewLRUCache()
				shardedLRU[i].cap = shardCap
			}

			fmt.Println("OK")

		} else if strings.HasPrefix(line, "NOW") {

			str := strings.TrimSpace(line[4:])

			time, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println(err)
			}

			now = time

		} else if strings.HasPrefix(line, "PUT") {

			strs := strings.Split(line[4:], " ")

			num, err := strconv.Atoi(strs[1])
			if err != nil {
				fmt.Println(err)
			}

			shardIndex := keyHash(strs[0]) % len(shardedLRU)
			shardedLRU[shardIndex].put(strs[0], num)

			if len(strs) > 2 {
				sec, err := strconv.Atoi(strs[2])
				if err != nil {
					fmt.Println(err)
				}
				shardedLRU[shardIndex].putTime(strs[0], sec)
			}

			fmt.Printf("OK shard=%d", shardIndex)
			fmt.Println()

		} else if strings.HasPrefix(line, "GET") {

			str := strings.TrimSpace(line[4:])

			shardIndex := keyHash(str) % len(shardedLRU)
			node, _ := shardedLRU[shardIndex].get(str)

			if node == nil {
				fmt.Println("<nil>")
				continue
			}

			if node.time > now || node.time == -1 {
				fmt.Printf("%d shard=%d", node.value, shardIndex)
				fmt.Println()
				continue
			}

			shardedLRU[shardIndex].expire(str)
			fmt.Println("<nil>")

		} else if strings.HasPrefix(line, "STATS") {
			hits := 0
			misses := 0

			for _, v := range shardedLRU {
				hits += v.hits
				misses += v.misses
			}

			var hitRate float64
			if hits == 0 && misses == 0 {
				hitRate = 0
			} else {
				hitRate = float64(hits) / float64(hits+misses)
			}

			fmt.Printf("hits=%d misses=%d hit_rate=%.2f", hits, misses, hitRate)
			fmt.Println()

		} else if strings.HasPrefix(line, "EVICTIONS") {

			for _, v := range shardedLRU {
				evictions += v.evictions
			}
			fmt.Println(evictions)
		}
	}
}

func keyHash(key string) int {
	var h uint32
	for i := 0; i < len(key); i++ {
		h = 31*h + uint32(key[i])
	}
	return int(h)
}

type Node struct {
	prev, next *Node
	key        string
	value      int
	time       int
}

func NewNode() *Node {
	return &Node{
		prev:  nil,
		next:  nil,
		key:   "",
		value: 0,
		time:  -1,
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
	cap       int
	hits      int
	misses    int
	evictions int
	dll       *DLL
	cache     map[string]*Node
}

func NewLRUCache() *LRUCache {
	return &LRUCache{
		cap:       0,
		hits:      0,
		misses:    0,
		evictions: 0,
		dll:       NewDLL(),
		cache:     make(map[string]*Node),
	}
}

func (c *LRUCache) setCapacity(cap int) {
	c.cap = cap
}

func (c *LRUCache) get(key string) (*Node, int) {
	ops := 0
	if node, ok := c.cache[key]; ok {
		c.hits++
		ops++

		c.dll.removeKey(node)
		ops++

		c.dll.moveFront(node)
		ops++
		return node, ops
	} else {
		c.misses++
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
		c.evictions++
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

func (c *LRUCache) putTime(key string, time int) {
	c.cache[key].time = time
}

func (c *LRUCache) expire(key string) {
	c.dll.removeKey(c.cache[key])
	delete(c.cache, key)
}

func (c *LRUCache) state() {
	c.dll.list()
}

func (c *LRUCache) stats() {
	var hitRate float64
	if c.hits == 0 && c.misses == 0 {
		hitRate = 0
	} else {
		hitRate = float64(c.hits) / float64(c.hits+c.misses)
	}

	fmt.Printf("hits=%d misses=%d hit_rate=%.2f", c.hits, c.misses, hitRate)
}
