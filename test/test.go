package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	m := make(map[int]int)
	m[1] = 1
	m[2] = 2
	Check(m)
}

func Check(m map[int]int) {
	var once sync.Once
	for {
		// once.Do只允许执行一次
		once.Do(func() {
			for index, item := range m {
				go func(timeDuration int, o int) {
					t1 := time.NewTicker(time.Duration(timeDuration) * time.Second)
					for {
						select {
						case <-t1.C:
							if o == 1 {
								fmt.Println("11111111111")
							} else if o == 2 {
								fmt.Println("22222222222")
							} else {
								fmt.Println("3333333333")
							}
						}
					}
				}(index, item)
			}
		})
	}
}
