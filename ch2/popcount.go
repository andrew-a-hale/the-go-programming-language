package main

import (
	"fmt"
	"testing"
	"time"
)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountLoop(x uint64) int {
	var count int
	for i := 0; i < 8; i++ {
		count += int(pc[byte(x>>(i*8))])
	}
	return count
}

func PopCountShift(x uint64) int {
	var count int
	for i := 0; i < 64; i++ {
		if x&1 > 0 {
			count++
		}
		x >>= 1
	}
	return count
}

func PopCountClear(x uint64) int {
	var count int
	for x != 0 {
		x &= x - 1
		count++
	}
	return count
}

func BenchmarkPopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount(uint64(0x73fa55623a30abed))
	}
}

func BenchmarkPopCountLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountLoop(uint64(0x73fa55623a30abed))
	}
}

func BenchmarkPopCountShift(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountShift(uint64(0x73fa55623a30abed))
	}
}

func BenchmarkPopCountClear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountClear(uint64(0x73fa55623a30abed))
	}
}

func runBenchmark(f func(b *testing.B)) {
	res := testing.Benchmark(f)
	fmt.Println(res.N, res.T, time.Duration(int64(res.T)/int64(res.N)))
}

func main() {
	fmt.Println(PopCount(uint64(0x73fa55623a30abed)))
	fmt.Println(PopCountLoop(uint64(0x73fa55623a30abed)))
	fmt.Println(PopCountShift(uint64(0x73fa55623a30abed)))
	fmt.Println(PopCountClear(uint64(0x73fa55623a30abed)))
	runBenchmark(BenchmarkPopCount)
	runBenchmark(BenchmarkPopCountLoop)
	runBenchmark(BenchmarkPopCountShift)
	runBenchmark(BenchmarkPopCountClear)
}
