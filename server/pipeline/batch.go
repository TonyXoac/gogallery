package pipeline

import (
	"runtime"
	"sync"
)

type BatchProcessing[T any] struct {
	work      func(T) error
	chunkSize int
}

func chunkSlice[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func (batch *BatchProcessing[T]) Run(items []T) {
	var wg sync.WaitGroup
	for _, chunk := range chunkSlice(items, batch.chunkSize) {
		wg.Add(1)
		go batch.processing(chunk, &wg)
	}
	wg.Wait()
}

func (poc *BatchProcessing[T]) processing(batch []T, wg *sync.WaitGroup) {
	for _, pic := range batch {
		poc.work(pic)
	}
	wg.Done()
}

func NewBatchProcessing[T any](processing func(T) error) *BatchProcessing[T] {
	proc := BatchProcessing[T]{}
	proc.work = processing
	proc.chunkSize = runtime.NumCPU()
	return &proc
}
