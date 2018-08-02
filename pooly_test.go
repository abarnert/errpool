package main

import (
	"testing"
)

func BenchmarkSpam(b *testing.B) {
	for n := 0; n < b.N; n++ {
		doitSpam(n*5)
	}
}


func BenchmarkSemaphore(b *testing.B) {
	for n := 0; n < b.N; n++ {
		doitSemaphore(n*5)
	}
}

func BenchmarkSemaphoreBlock(b *testing.B) {
	for n := 0; n < b.N; n++ {
		doitSemaphoreBlock(n*5)
	}
}

func BenchmarkPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		doitPool(n*5)
	}
}

func BenchmarkErrPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		doitErrPool(n*5)
	}
}

func TestSpam(t *testing.T) {
	doitSpam(10)
}

func TestSemaphore(t *testing.T) {
	doitSemaphore(10)
}

func TestSemaphoreBlock(t *testing.T) {
	doitSemaphoreBlock(10)
}

func TestPool(t *testing.T) {
	doitPool(10)
}

func TestErrPool(t *testing.T) {
	doitErrPool(10)
}
