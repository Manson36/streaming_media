package main

import (
	"fmt"
	"testing"
)

func testPrint(t *testing.T) {
	res := Print1to20()
	fmt.Println("hey")
	if res != 210 {
		t.Errorf("Wrong result of Print1to20")
	}
}

func testPrint2(t *testing.T) {
	res := Print1to20()
	res ++
	if res != 211 {
		t.Errorf("testPrint2 failed")
	}
}

func TestAll(t *testing.T) {
	t.Run("TestPrint", testPrint)
	t.Run("TestPrint2", testPrint2)
}

func TestMain(m *testing.M) {
	fmt.Println("Tests begin...")
	m.Run()
}

func BenchmarkAll(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Print1to20()
	}
}