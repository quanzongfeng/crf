package crfpp

import (
    "math"
    "heap"
)

const MINUS_LOG_ESPILON = 50

func logsumexp(x, y float64, flg bool) float64 {
    if flg {
        return y
    }

    vmin := math.Min(x, y)
    vmax := math.Max(x, y)

    if vmax > vmin + MINUS_LOG_ESPILON {
        return vmax
    } 
    return vmax + math.Log(math.Exp(vmin - vmax) + 1.0)
}

//define Node used in vitebi 
type Node struct {
    x           int
    y           TAG_TYPE    //present tag type
    alpha       float64     //forward probabilities
    beta        float64     //backward probabilities
    cost        float64     //p(xi|yi)
    bestcost    float64
    prev        *Node       //best prev node
    fvector     []ID_TYPE   //feature ids
    lpath       []*Path     //left transfer
    rpath       []*Path     //right transfer
}

//define for transfermition
type Path struct {
    rnode   *Node
    lnode   *Node
    fvector []ID_TYPE
    cost    float64
}


func (n *Node)CalcCost(expected []float64) {
    n.cost = 0.0
    for _,s := range n.fvector {
        n.cost += expected[s + ID_TYPE(n.y)]
    }
}

func (n *Node) CalcAlpha() {
    n.alpha = 0.0
    for i,s := range n.lpath {
        n.alpha = logsumexp(n.alpha, s.cost+s.lnode.alpha, i==0)
    }
    n.alpha += n.cost
}

func (n *Node) CalcBeta() {
    n.beta = 0.0
    for i, s:= range n.rpath {
        n.beta = logsumexp(n.beta, s.cost+s.rnode.beta, i==0)
    }
    n.beta += n.cost
}

func (n *Node) CalcExpectation(expected map[ID_TYPE]float64, z float64, size int) {
    c := math.Exp(n.alpha + n.beta - n.cost - z)
    for _,s := range n.fvector {
        expected[s + ID_TYPE(n.y)] += c
    }

    for _,lp:= range n.lpath {
        lp.CalcExpectation(expected, z, size)
    }
}

func (n *Node) Clear() {
    n.x = 0
    n.y = 0
    n.alpha, n.beta, n.cost = 0.0, 0.0, 0.0
    n.bestcost = 0.0
    n.prev = nil
    n.fvector = nil
    n.lpath = nil
    n.rpath = nil
}


func (p *Path) CalcCost(expected []float64, size int) {
    p.cost = 0.0
    for _,s:= range p.fvector {
        p.cost += expected[s + ID_TYPE(p.lnode.y) * ID_TYPE(size) + ID_TYPE(p.rnode.y)]
    }
}

func (p *Path) CalcExpectation(expected map[ID_TYPE]float64,  z float64, size int) {
    c := math.Exp(p.lnode.alpha + p.cost + p.rnode.beta - z)
    for _,s := range p.fvector {
        expected[s + ID_TYPE(p.lnode.y) * ID_TYPE(size) + ID_TYPE(p.rnode.y)] += c
    }
}

func (p *Path) Add(l *Node, r *Node) {
    p.rnode = r
    p.lnode = l
    l.rpath = append(l.rpath, p)
    r.lpath = append(r.lpath, p)
}

func (p *Path) Clear() {
    p.rnode = nil
    p.lnode = nil
    p.fvector = nil
    p.cost = 0.0
}

//*Node realize heap.ICompare interface, so *Node can make a heap.HeapFactory as elements
func (n *Node)Value() float64 {
    return n.bestcost
}

func (n *Node)Compare(b heap.Value)(bool, error) {
    return n.Value() < b.Value(), nil
}
    
