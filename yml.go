package main

import (
	"sync"
)

func main() {
	var locks = new(sync.Map)

	lock, _ := locks.LoadOrStore("a", new(sync.Mutex))

	lock.(*sync.Mutex).Lock()

	println("===========")
	lock.(*sync.Mutex).Lock()
	println("===========")
	lock.(*sync.Mutex).Unlock()

}
