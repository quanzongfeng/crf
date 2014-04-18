package main

import (
    "heap"
    "fmt"
)

func main() {
    a := []float64{1.0,5.0,4.0,3.0,8.0,6.0,7.0}
    h := heap.GetHeap(a)
    fmt.Println(h.Items)
    fmt.Println(h.Top())
    h.Push(3.3)
    fmt.Println(h.Top(),h.Items, h.Size)
    h.Push(9.0)
    fmt.Println(h.Top(),h.Items, h.Size)
    h.Pop()
    h.Push(10.0)
    h.Push(2.5)
    h.Push(7.1)
    h.Push(8.7)
    fmt.Println(h.Top(),h.Items, h.Size)
    h.Sort()
    fmt.Println(h.Items)
    
}
