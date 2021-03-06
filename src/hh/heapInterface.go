package hh

import (
    //    "fmt"
    "stringError"
)

type HeapFactory struct {
    Size    int
    Items   []interface{}
}

type IGreate interface {
    Greater(b interface{})(bool, error)
}

type Int int
type Float float64

func (f Float)Greater(b interface{})(bool, error) {
    fb,ok := b.(float64) 
    if !ok {
        return false, stringError.StringError("parameters not float64")
    }
    return float64(f) > float64(fb), nil
}

