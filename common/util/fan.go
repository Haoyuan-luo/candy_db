package util

import "sync"

// 场景函数
type sceneType[T, R any] func(chan R, T, func(T) R)

// 工作函数
type fanType[T, R any] func(T) R

func Producer[T any](done <-chan struct{}, stream []T) <-chan T {
	fan := make(chan T)
	go func() {
		defer close(fan)
		for _, k := range stream {
			select {
			case fan <- k:
			case <-done:
				return
			}
		}
	}()
	return fan
}

func Processor[T, R any](done <-chan struct{}, stream <-chan T, fn fanType[T, R], scene sceneType[T, R]) <-chan R {
	fan := make(chan R)
	go func() {
		defer close(fan)
		for v := range stream {
			select {
			case <-done:
				return
			default:
				scene(fan, v, fn)
			}
		}
	}()
	return fan
}

func Consumer[R any](done <-chan struct{}, streams ...<-chan R) <-chan R {
	fan := make(chan R)
	var wg sync.WaitGroup
	wg.Add(len(streams))
	for _, stream := range streams {
		go func(stream <-chan R) {
			defer wg.Done()
			for v := range stream {
				select {
				case fan <- v:
				case <-done:
					return
				}
			}
		}(stream)
	}
	go func() {
		wg.Wait()
		close(fan)
	}()
	return fan
}
