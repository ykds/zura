package queue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	queue := NewTimeoutQueue(10, 10)
	queue.Push("test1", func() {
		fmt.Println("1")
	})
	queue.Push("test2", func() {
		fmt.Println("2")
	})
	queue.Push("test3", func() {
		fmt.Println("3")
	})
	go func() {
		time.Sleep(3 * time.Second)
		queue.Finish("test2")
	}()
	go func() {
		time.Sleep(time.Second * 30)
		queue.Push("test4", func() {
			fmt.Println("4")
		})
	}()
	queue.Run(context.TODO())
}
