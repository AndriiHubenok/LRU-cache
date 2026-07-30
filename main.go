package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO (cache-basics): implement per the lesson description.

type LRUCache struct {
	cap  int
	keys []string
	hits int
}

func (c *LRUCache) Access(key string) {
	if c.cap == 0 {
		return
	}

	for i, k := range c.keys {
		if k == key {
			c.hits++

			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			c.keys = append(c.keys, key)
			return
		}
	}

	if len(c.keys) == c.cap {
		c.keys = c.keys[1:]
	}
	c.keys = append(c.keys, key)
}

type FIFOCache struct {
	cap  int
	keys []string
	hits int
}

func (c *FIFOCache) Access(key string) {
	if c.cap == 0 {
		return
	}

	for _, k := range c.keys {
		if k == key {
			c.hits++
			return
		}
	}

	if len(c.keys) == c.cap {
		c.keys = c.keys[1:]
	}
	c.keys = append(c.keys, key)
}

type LFUCache struct {
	cap  int
	keys []string
	freq map[string]int
	hits int
}

func (c *LFUCache) Access(key string) {
	if c.cap == 0 {
		return
	}

	for _, k := range c.keys {
		if k == key {
			c.hits++
			c.freq[key]++
			return
		}
	}

	if len(c.keys) == c.cap {
		minFreq := -1
		evictIdx := -1

		for i, k := range c.keys {
			if minFreq == -1 || c.freq[k] < minFreq {
				minFreq = c.freq[k]
				evictIdx = i
			}
		}

		evictKey := c.keys[evictIdx]
		delete(c.freq, evictKey)
		c.keys = append(c.keys[:evictIdx], c.keys[evictIdx+1:]...)
	}

	c.keys = append(c.keys, key)
	c.freq[key] = 1
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)

	var lru LRUCache
	var fifo FIFOCache
	var lfu LFUCache

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "CAP") {
			numStr := strings.TrimSpace(line[4:])
			size, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			lru = LRUCache{cap: size, keys: make([]string, 0, size)}
			fifo = FIFOCache{cap: size, keys: make([]string, 0, size)}
			lfu = LFUCache{
				cap:  size,
				keys: make([]string, 0, size),
				freq: make(map[string]int),
			}
		}

		if strings.HasPrefix(line, "ACCESS") {
			value := strings.TrimSpace(line[7:])

			lru.Access(value)
			fifo.Access(value)
			lfu.Access(value)
		}

		if strings.HasPrefix(line, "STATS") {
			fmt.Printf("lru_hits=%d fifo_hits=%d lfu_hits=%d\n", lru.hits, fifo.hits, lfu.hits)
		}
	}
}
