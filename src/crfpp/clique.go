package crfpp

import (
    "io"
    "bytes"
    "fmt"
    "heap"
)

//present every nature lines

type Clique struct{
    id          int
    featureid   ID_TYPE
    size        int    //vertex num
    columns     int    //column num, text columns -1
    g           *Group
    lines       [][]string
    result      []TAG_TYPE
    answer      []TAG_TYPE
    nodes       [][]*Node
    z           float64
    s           float64
    cost        float64
    expect      map[ID_TYPE]float64
    learnflag   bool
    h           *heap.HeapFactory
}

type HeapElement struct {
    n   *Node
    gx  float64     //store sum cost of current path from end to n, used to calc fx
    fx  float64     //store -1*bestcost of current path
    next *HeapElement
}

func (h *HeapElement)Compare(b heap.Value )(bool, error) { //used for minheap
    return h.Value() < b.Value(), nil
}

func (h *HeapElement)Value() float64 {
    return h.fx
}


//type Line struct {
//    col     []string //column numbers
//    tag     uint
//}
//
//func getLine() *Line {
//    return new(Line)
//}

func (c *Clique)SetLearnFlag(b bool) {
    c.learnflag = b
}

func (c *Clique)SetGroup(g *Group){
    c.g = g
}

func (c *Clique)Size() int {
    return len(c.lines)
}

func (c *Clique)SetID(id int) {
    c.id = id
}
func (c *Clique)SetFeatureID(id ID_TYPE) {
    c.featureid = id
}

func ReadClique(r io.ReadSeeker, tags *Tag)(c *Clique, err error) {
    bf := make([]byte, 1024)
    lines := make([][]string,0, 8)
    answers:=make([]TAG_TYPE,0, 8)
    for {
        n, e := ReadLine(r, bf)
        if e != nil && e!= io.EOF {
            return nil,e
        }
        if n == 0 {
            err = e
            break
        }
        if bf[0] == ' ' || bf[0] == '\t' || bf[0] == '\n' {
            break
        }

        tokens := bytes.Fields(bf[:n])

        t := make([]string,0, 1)
        for _, s:= range tokens {
            t = append(t, string(s))
        }

        lencolum := len(t)-1
        m := make([]string, lencolum+1)
        copy(m, t)
        v, ok := tags.Find(t[lencolum])

        if c.learnflag { //when learn, need to judge tags
            if !ok {
                return nil, StringError(fmt.Sprintf("Unexpected tag: %s", string(t[lencolum])))
            }
        }

        if ok {
        answers = append(answers, v)    //when !ok, v = TAG_TYPE(0), just present in test
        } else {
            answers = append(answers, TAG_TYPE(-1))
        }

        lines = append(lines, m)
    }

    ll := len(lines)
    if ll == 0 {
        return nil, err
    }

    c = new(Clique)
    c.lines = make([][]string,  ll)
    c.answer = make([]TAG_TYPE, ll)
    c.result = make([]TAG_TYPE, ll)
    copy(c.lines, lines)
    copy(c.answer, answers)
    c.size = ll
    c.columns = len(c.lines[0]) - 1     //columns 不包括最后一列
    return c, err
}

func (c *Clique)BuildLattice() {
    if len(c.lines) == 0 {
        return 
    }

    if c.nodes == nil  {
        tagtypes := c.g.tags.GetTagNum()
        c.nodes = make([][]*Node,0, tagtypes)
        for i, _:= range c.lines {
            lnd := make([]*Node, 0, c.columns)
            for j := 0;j<tagtypes;j++ {
                n := new(Node)
                n.x = i
                n.y = TAG_TYPE(j)
                n.fvector = c.g.GetFvector(c.featureid, i,c.size, 0 )
                lnd = append(lnd, n)
            }
            c.nodes = append(c.nodes, lnd)
        }

        for i, _:= range c.lines {
            if i==0 {
                continue
            }
            for j :=0; j<tagtypes; j++ {
                for k:=0; k<tagtypes; k++ {
                    p:=new(Path)
                    p.Add(c.nodes[i-1][j], c.nodes[i][k])
                    p.fvector = c.g.GetFvector(c.featureid, i, c.size, 1)
                }
            }
        }
    }
    //    fmt.Println(len(c.nodes),c.nodes)
}

func (c *Clique)CalcCost(alpha[] float64) {
    for _, ln := range c.nodes {
        for _, n := range ln {
            n.CalcCost(alpha)
            for _,p := range n.lpath {
                p.CalcCost(alpha, c.g.tags.GetTagNum() )
            }
        }
    }

//    fmt.Printf("cost of clique:%d are :\n", c.id)
//    for _,ln:= range c.nodes {
//        for _,n:= range ln {
//            fmt.Printf("\t%f", n.cost)
//        }
//        fmt.Printf("\n")
//    }
//
//    n := c.nodes[0][0]
//    for _, t:= range n.rpath {
//        fmt.Printf("\tB\t%f", t.cost)
//    }
//    fmt.Println("")
}

func (c *Clique)ForwardBackward() {
    if len(c.nodes)==0 {
        return
    }

    for _,ln := range c.nodes {
        for _,n := range ln {
            n.CalcAlpha()
        }
    }

    l := len(c.nodes)
    //    fmt.Println(l)
    for j:=l-1; j>=0;j-- {
        ln := c.nodes[j]
        for _, n := range ln {
            n.CalcBeta()
        }
    }

    c.z = 0.0
    for j, n := range c.nodes[0] {
        c.z = logsumexp(c.z, n.beta, j==0)
    }
}


func (c *Clique)Vitebi() {
    for _,ln := range c.nodes {
        for _, n:= range ln {
            bestc := -1e37
            var best *Node = nil
            for _, p:= range n.lpath {
                cost := p.lnode.bestcost + p.cost + n.cost
                if cost > bestc {
                    bestc = cost
                    best = p.lnode
                }
            }
            n.prev = best
            if best != nil {
                n.bestcost = bestc
            }else {
                n.bestcost = n.cost
            }
        }
    }

    bestc := -1e37
    var best *Node = nil
    for _, n := range c.nodes[len(c.nodes)-1] {
        if bestc < n.bestcost {
            bestc = n.bestcost
            best = n
        }
    }

    for n:= best;n!= nil ;n= n.prev {
        c.result[n.x] = n.y
    }

    c.cost = -bestc
}

func (c *Clique)initBest() {
    if c.h == nil {
        c.h = heap.NewHeap()
    }

    for !c.h.Empty() {
        c.h.Pop()
    }

    k := len(c.nodes)-1
    j := len(c.nodes[k])
    for i:=0;i<j;i++ {
        h := new(HeapElement)
        h.n = c.nodes[k][i]
        h.gx = -1*h.n.cost
        h.fx = -1*h.n.bestcost
        h.next = nil
        c.h.Push(h)
    }
}

func (c *Clique)next() bool {
    if c.h == nil {
        return false
    }
    for !c.h.Empty() {
        mtop := c.h.Pop()
        top, ok := mtop.(*HeapElement)
        if !ok {
            panic("error heap type")
        }

        
        n := top.n

        if n.x == 0 {   //find a way
            for t:=top;t != nil ;t=t.next {
                c.result[t.n.x] = t.n.y
            }
            c.cost= top.gx  //different path, different cost
            return true
        }

        for _, p:= range n.lpath {
            h := new(HeapElement)
            h.n = p.lnode
            h.gx = -1*h.n.cost - p.cost + top.gx
            h.fx = -1*h.n.bestcost - p.cost + top.gx
            h.next = top
            c.h.Push(h)
        }
    }

    return false
}

func (c *Clique)GetResult() []*Node{
    ll := make([]*Node, 0,1)
    for i, t := range c.result {
        ll = append(ll, c.nodes[i][int(t)])
    }
    return ll
}

func (c *Clique)GetPath(num int) ([][]*Node,[]float64) {
    re := make([][]*Node, 0, 1)
    cl := make([]float64,0,1)
    if num == 1 {
        re = append(re, c.GetResult())
        cl = append(cl, c.cost)
        return re, cl
    }
    for i:=0; i<num; i++ {
        if !c.next() {
            break
        }
        re = append(re, c.GetResult())
        cl = append(cl, c.cost)
    }
    return re, cl
}

func (c *Clique)CalcBest(alpha []float64) {
    c.ClearResult()
    c.BuildLattice()
    c.CalcCost(alpha)
    c.ForwardBackward()
    c.Vitebi()
    c.initBest()
}

func (c *Clique) CalcCRFGo(alpha[]float64, ch chan int) {
    _,e := c.Gradient(alpha)
    if e!= nil {
        fmt.Println("Clique ",c.id, "failed: ", e.Error())
        panic("failed")
    }
    ch <- c.id
}

func (c *Clique)CalcBestGo(alpha []float64, ch chan int) {
    c.CalcBest(alpha)
    ch <- c.id
}


func (c *Clique) Gradient(alpha[]float64) (float64, error) {
    if len(c.lines) == 0 {
        return 0.0, nil
    }

    c.ClearResult()

    c.BuildLattice()
    c.CalcCost(alpha)
    c.ForwardBackward()
    tagtypes := c.g.tags.GetTagNum()

    for _, ln := range c.nodes {
        for _, n := range ln {
            n.CalcExpectation(c.expect, c.z, tagtypes)
        }
    }


    for i:=0; i<c.size; i++ {
        n := c.nodes[i][c.answer[i]]
        for _, f := range n.fvector {
            d := ID_TYPE(f) + ID_TYPE(c.answer[i]) 
            v, ok := c.expect[d]
            if !ok {
                return 0.0, StringError("unexpected error")
            }
            v--
            c.expect[d] = v

        }

        c.s += n.cost
        for _,p := range n.lpath {
            if (p.lnode.y == c.answer[p.lnode.x]) {
                for _, f:= range p.fvector {
                    d := ID_TYPE(f) + ID_TYPE(p.lnode.y) * ID_TYPE(tagtypes)+ ID_TYPE(p.rnode.y)
                    v, ok := c.expect[d]
                    if !ok {
                        return 0.0, StringError("unexpected error")
                    }
                    v--
                    c.expect[d] = v

                }

                c.s += p.cost
                break
            }
        }
    }

    c.Vitebi()
//    fmt.Println("clique message: ", c.id, c.size, c.z, c.s, c.z-c.s)
    return c.z - c.s, nil
}

func (c *Clique) Eval() int {
    var i, d int = 0, 0
    for i=0; i<c.size;i++ {
        if c.answer[i] != c.result[i] {
            d += 1
        }
    }
    return d
}

func (c *Clique) ClearNodes() {
    c.nodes = nil
}
func (c *Clique)ClearResult() {
    c.s = 0.0
    c.z = 0.0
    c.result = make([]TAG_TYPE, c.size)
    c.expect = make(map[ID_TYPE]float64)
}

func (c *Clique) BuildFeatureGo(out chan *featureMsgIn, t *Template) {
    un, bn := t.Size()
    all := c.size * un +(c.size-1)*bn
    sout := make([]*featureMsgIn,0, all)
    for i, _ := range c.lines {
        fout := new(featureMsgIn)
        fout.id = c.id
        fout.col = i
        fout.ty = U_TYPE
        lk, e:=t.ApplyRuleUnigram(i, c)
        if e!= nil {
            fmt.Println(e.Error())
            panic("here")
        }
        fout.keys = lk
        sout = append(sout, fout)
    }
    for i, _ := range c.lines {
        if i == 0{
            continue
        }
        fout := new(featureMsgIn)
        fout.id = c.id
        fout.col = i
        fout.ty = B_TYPE
        lk, e := t.ApplyRuleBigram(i, c)
        if e!= nil {
            fmt.Println(e.Error())
        }
        fout.keys = lk
        sout = append(sout, fout)
    }
    for _, t:= range sout {
        out <- t
    }

    fend := new(featureMsgIn)
    fend.id = c.id
    fend.col = -1
    out<-fend
}

