package main

import (
    "crfpp"
    "fmt"
    "os"
    "io"
)

func main() {
    f , err:= os.Open("temp")
    if err != nil {
        fmt.Println(err.Error())
        return 
    }

    var s  [1024]byte
    bf := s[:]
    fmt.Println(len(bf), cap(bf), len(s), cap(s))
    fmt.Printf("%p \n%p\n", bf, &s)
    for i:=0;;i++{
        n, e:=crfpp.ReadLine(f, bf)
        if e!= nil && e != io.EOF{
            fmt.Println(e.Error())
            break
        }
        fmt.Printf("n=%p, e=%p\n", &n, e)
        fmt.Println(bf[:n], e==nil)
        fmt.Println(string(bf[:n]))
        fmt.Println(i)
        if e == io.EOF {
            break
        }
    }
//    for i:=0;i<2048;i++ {
//        bf = append(bf, 1)
//        fmt.Printf("%p\n", bf)
//    }
    fmt.Println(len(bf), cap(bf), len(s), cap(s))

}


