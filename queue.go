package main

import "errors"

type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Enqueue(value T) {
	q.elements = append(q.elements, value)
}

func (q *Queue[T]) Dequeue() (T, error) {

	if q.IsEmpty() {
		var zero T
		return zero, errors.New("queue is empty")
	}

	value := q.elements[0]

	var zero T
	q.elements[0] = zero
	q.elements = q.elements[1:]

	return value, nil
}

func (q *Queue[T]) Peek() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, errors.New("queue is empty")
	}
	return q.elements[0], nil
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.elements) == 0
}
