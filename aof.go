package main

import (
	"bufio"
	"os"
	"sync"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}
