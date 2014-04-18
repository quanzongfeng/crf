package main

import(
    "crfpp"
    "fmt"
)

func main() {
    gr := new(crfpp.Group)
    e := gr.ReadFile("template", "train.txt")
    if e!= nil {
        fmt.Println(e)
        return
    }else {
        fmt.Println("gr is ready")
    }
    e = gr.ProcessData("model.txt", 1, 0.0001, 4, 10000 )
    if e!= nil {
        fmt.Println(e.Error())
    }
}


