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

	var lru *NaiveLRU

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "CAP") {

			numStr := strings.TrimSpace(line[4:])
			num, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Println(err)
			}

			lru = NewNaiveLRU(num)
			fmt.Println("OK")

		} else if strings.HasPrefix(line, "PUT") {

			strs := strings.Split(line[4:], " ")
			num, err := strconv.Atoi(strs[1])
			if err != nil {
				fmt.Println(err)
			}

			lru.put(strs[0], num)
			fmt.Println("OK")

		} else if strings.HasPrefix(line, "GET") {

			str := strings.TrimSpace(line[4:])

			fmt.Println(lru.get(str))
		} else if strings.HasPrefix(line, "STATE") {
			lru.state()
		}
	}
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
