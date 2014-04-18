package main

import (
    "fmt"
    "lbfgs"
)

func gradient(x float64) float64{
    return 4.0 * x * x *x
}

func f(x float64)float64 {
    return x * x* x*x
}

func main () {
    var x, gx []float64
    x = []float64{10.0}
    for {
        gx = []float64{gradient(x[0])}
        fx := f(x[0])

        fmt.Println(x, gx, fx)
        t, e := lbfgs.Optimize(1, x, fx, gx, false, 1)
        if t <=0 {
            fmt.Println(x, e, t)
            break
        }
        
    }
}


