package main

type Queue struct {
	queue []int

	addItem chan int
	getItem chan int
	getSize chan int
}

func CreateQueue() *Queue {
	return &Queue{
		queue:   make([]int, 0),
		addItem: make(chan int),
		getItem: make(chan int),
		getSize: make(chan int),
	}
}

func (q *Queue) Start() {
	go func() {
		for {
			if len(q.queue) > 0 {
				select {
				case item := <-q.addItem:
					q.queue = append(q.queue, item)

				case q.getSize <- len(q.queue):

				case q.getItem <- q.queue[0]:
					q.queue = q.queue[1:]
				}
			} else {
				select {
				case item := <-q.addItem:
					q.queue = append(q.queue, item)

				case q.getSize <- len(q.queue):
				}
			}
		}
	}()
}
