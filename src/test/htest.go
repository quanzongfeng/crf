package main

import (
    "heap"
    "fmt"
)

func main() {
    a, b:=0.0, 1.0

    fmt.Println(heap.MinFloat64(a).Greater(heap.MinFloat64(b)))
    f :=  heap.MinFloat64(a)
    g := heap.MinFloat64(b)
    fmt.Println(f.Greater(g))
    fmt.Println(f.Value())


    t:=make([]heap.IGreater, 10,10)

    for i,_:=range t {
        t[i] = heap.MinFloat64(i*1.0)
    }

    fheap := new(heap.HeapFactory)
    fheap.MakeHeap(t)
    fmt.Println(t)
    fmt.Println(fheap)
    fheap.Push(heap.MinFloat64(2.2))
    fmt.Println(fheap)

    num := fheap.Size
    for i:=0;i< num;i++ {
        fmt.Println(fheap.Pop())
    }

}
    
