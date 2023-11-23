package ch1

import (
	"errors"
	"fmt"
	"testing"
)

func checkShrink(cap int, len int) int {
	if (cap % len) >= 2 {
		return len
	} else {
		return cap
	}
}
func Remove[T any](s []T, index int) (res []T, err error) {
	if index > len(s)-1 || index < 0 {
		return s, errors.New("Index > Slice len or index < 0")
	}
	s1Cap := checkShrink(cap(s), len(s))
	s1 := make([]T, len(s)-1, s1Cap)
	copy(s1[:index], s[:index])
	copy(s1[index:len(s1)], s[index+1:len(s)])
	return s1, nil
}

func StringSlice_remove() {
	s1 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
	s2, _ := Remove(s1, 3)
	fmt.Println(s2, cap(s2))
	s3, _ := Remove(s2, 0)
	fmt.Println(s3, cap(s3))
	s4, _ := Remove(s3, 3)
	fmt.Println(s4, cap(s4))
	s5, _ := Remove(s4, -1)
	fmt.Println(s5, cap(s5))
	s6, _ := Remove(s5, 4)
	fmt.Println(s6, cap(s6))
	s7, _ := Remove(s6, len(s6)-1)
	fmt.Println(s7, cap(s7))
	s8, _ := Remove(s7, 3)
	fmt.Println(s8, cap(s8))
}
func BenchmarkSlice(b *testing.B) {
	b.StartTimer()
	StringSlice_remove() //BenchmarkSlice-24       1000000000               0.0005196 ns/op
	b.StopTimer()
}
