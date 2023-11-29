package main

import "testing"

func TestSum(t *testing.T){
	result := sum(2,3)
	if result != 5{
		t.Error("Result must be 5")
	}
}
func TestSub(t *testing.T){
	result := sub(5,3)
	if result != 2{
		t.Error("Result must be 2")
	}
}
func TestTimes(t *testing.T){
	result := times(4,3)
	if result != 12{
		t.Error("Result must be 12")
	}
}
func TestDivide(t *testing.T){
	result := divide(9,3)
	if result != 3{
		t.Error("Result must be 3")
	}
}
func TestSumX(t *testing.T){
	result := sumX(1,3,5)
	if result != 9{
		t.Error("Result must be 9")
	}
}