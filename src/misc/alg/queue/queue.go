package queue

type Queue struct {
	capacity int32
	size	int32
	front	int32
	rear	int32
	elements []interface{}
}

func New(max int32) *Queue {
	queue := &Queue {capacity:max, size:0, front:0, rear:-1}
	queue.elements = make([]interface{}, max)
	return queue
}

//----------------------------------------------- Enqueue
func (q *Queue) Enqueue(elem interface{}) bool {
	if q.size < q.capacity {
		q.size++
		q.rear++

		if (q.rear == q.capacity) {
			q.rear = 0
		}

		q.elements[q.rear] = elem

		return true
	}

	return false
}

//----------------------------------------------- Dequeue
func (q *Queue) Dequeue()(ret interface{}, result bool) {
	if q.size > 0 {
		ret = q.elements[q.front]
		result = true

		q.size--
		q.front++

		if q.front == q.capacity {
			q.front = 0
		}
		return
	}

	return nil, false
}

//----------------------------------------------- return queue
func (q *Queue) All()(all []interface{}) {
	all = make([]interface{}, q.size)

	count := q.size
	idx := q.front

	for count > 0  {
		all[idx] = q.elements[idx]

		if idx < q.capacity {
			idx++
		} else {
			idx = 0
		}

		count--
	}

	return
}
