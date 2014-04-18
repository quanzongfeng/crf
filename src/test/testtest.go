package main

import (
    "crfpp"
    "fmt"
)

func main() {
    ml, e := crfpp.LoadModel("model.txt")
    if e!= nil {
        fmt.Println(e.Error())
        return
    }
    if ml == nil {
        fmt.Println("load model file failed")
    }
    return
}
