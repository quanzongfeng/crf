package heap

import (
    "fmt"
    "stringError"
//    "reflect"
)

type HeapFactory struct {
    Size    int
    Items   []IComparer
}

type IComparer interface {
    Compare(b Value) (bool, error)
    Value() float64
}

type Value interface {
    Value() float64
}

//对应最大堆
type Float64 float64
func (f Float64)Value() float64 {
    return float64(f)
}

func(f Float64)Compare(b Value)(bool, error) {
    return f.Value() > b.Value(), nil
}

//对应最小堆
type MinFloat64 float64
func (f MinFloat64)Value() float64 {
    return float64(f)
}
func (f MinFloat64)Compare(b Value)(bool, error) {
    return f.Value() < b.Value(), nil
}
       
func NewHeap() *HeapFactory {
    h := new(HeapFactory)
    h.Items = make([]IComparer,1)
    return h
}

//type IHeap interface {
//    func Top() interface{}
//    func Pop()
//    func Push(interface{})
//    func Empty() bool 
//}

func (hf *HeapFactory)Adjust(s, ln int) (error) {
    if len(hf.Items) < ln {
        return stringError.StringError(fmt.Sprintf("$3:%d is longer then len of $1", ln))
    }
    nchild := 0
    for ; 2*s + 1 < ln ;s = nchild {
        nchild = 2*s +1
        if nchild < ln -1 {
            b, e:=hf.Items[nchild+1].Compare(hf.Items[nchild])
            if e!= nil {
                return e
            }
            if b{
                nchild++
            }
        }
        b, e:= hf.Items[nchild].Compare(hf.Items[s])
        if e!= nil {
            return e
        }
        if b {
            hf.Items[nchild], hf.Items[s] = hf.Items[s], hf.Items[nchild]
        }else {
            break
        }
    }
    return nil
}

func (hf *HeapFactory)AdjustBack(s, ln int)(error) {
    if len(hf.Items) < ln {
        return stringError.StringError(fmt.Sprintf("$3:%d is longer then len of $1", ln))
    }

    np := (s-1)/2 
    for ;np >=0 && s>0; np = (s-1)/2 {
        nchild := 2*np +1
        if nchild < ln -1 {
            b, e:=hf.Items[nchild+1].Compare(hf.Items[nchild])
            if e!= nil {
                return e
            }
            if b{
                nchild++
            }
        }
        b ,e := hf.Items[nchild].Compare(hf.Items[np])
        if e!= nil {
            return e
        }
        if b {
            hf.Items[nchild], hf.Items[np] = hf.Items[np], hf.Items[nchild]
        }else {
            break
        }
    }
    return nil
}

func GetHeap(ig []IComparer)*HeapFactory {
    hf := new(HeapFactory)
    hf.MakeHeap(ig)
    return hf
}


func (hf *HeapFactory)MakeHeap(ig []IComparer) {
    hf.Size = len(ig)
    hf.Items = ig
    hf.makeHeap()
}

func (hf *HeapFactory)makeHeap() {
    ln := hf.Size
    for i:=(ln-2)/2; i>=0; i-- {
        hf.Adjust(i, ln)
    }
}

func (hf *HeapFactory)Sort() {
    ln := hf.Size
    for i:=ln-1; i>0; i-- {
        hf.Items[0], hf.Items[i] = hf.Items[i], hf.Items[0]
        hf.Adjust(0, i)
    }
}

func (hf *HeapFactory)Top() IComparer {
    return hf.Items[0]
}

func (hf *HeapFactory)Pop() IComparer {
    ln := hf.Size
    hf.Items[0], hf.Items[ln-1] = hf.Items[ln-1], hf.Items[0]
    hf.Adjust(0, ln-1)
    hf.Size--
    return hf.Items[ln-1]
}

func (hf *HeapFactory)Push(a IComparer) {
    ln := len(hf.Items)
    if ln > hf.Size {
        hf.Items[hf.Size] = a
    }else {
        hf.Items = append(hf.Items, a)
    }
    hf.Size++
    hf.AdjustBack(hf.Size-1, hf.Size)
}

func (hf *HeapFactory)Empty() bool {
    return hf.Size == 0
}



