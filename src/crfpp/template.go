package crfpp

import (
    "os"
    "io"
//    "fmt"
//    "strconv"
)

type Template struct {
    unigram     []string
    bigram      []string
}

const (
    U_TYPE = 0
    B_TYPE = 1
)

const KMaxContextSize = 8
var    BOS []string =[]string{"_B-1", "_B-2", "_B-3", "_B-4", "_B-5","_B-6", "_B-7", "_B-8"}
var    EOS []string =[]string{"_B+1", "_B+2", "_B+3", "_B+4", "_B+5","_B+6", "_B+7", "_B+8"}


func ReadTemplateFromFile(name string) (*Template,error) {
    f, e := os.Open(name)
    if e!= nil {
        return nil, e
    }
    defer f.Close()

    tm := new(Template)

    e =tm.ReadTemplate(f)
    return tm,e
}


func (t *Template)ReadTemplate(r io.ReadSeeker) (err error) {
    bf := make([]byte, 1024)
    for {
        n, e := ReadLine(r, bf)
        if e != nil && e!= io.EOF {
            return e
        }

        if n == 0 {  //empty line
            if e == io.EOF{ //file end
                break
            }else {
                continue
            }
        }

        if bf[0] == ' ' || bf[0] == '\t' || bf[0] == '\n' || bf[0] == '#'{
            continue
        }

        switch bf[0] {
        case 'U':
            err = t.processUnigram(bf[0:n])
        case 'B':
            err = t.processBigram(bf[0:n])
        default:
            return StringError("bad format")
        }
        if err != nil {
            return err
        }
        if e==io.EOF {
            break
        }
    }
    return nil
}

func (t *Template)processUnigram(line []byte) (error) {
    if len(line) ==0 {
        return StringError("empty line")
    }
    if line[0] != 'U'{
        return StringError("Not Unigram")
    }

    if t.unigram == nil {
        t.unigram = make([]string, 0, 5)
    }
    t.unigram = append(t.unigram, string(line))
    return nil
}

func (t *Template)processBigram(line []byte) (error) {
    if len(line) ==0 {
        return StringError("empty line")
    }
    if line[0] != 'B'{
        return StringError("Not Bigram")
    }
    
    if t.bigram == nil {
        t.bigram = make([]string, 0, 1)
    }
    t.bigram = append(t.bigram, string(line))
    return nil
}

func (t *Template)Size() (int, int) {
    return  len(t.unigram), len(t.bigram)
}


func (t *Template)ApplyRuleUnigram(index int, c *Clique)(keys []string, err error) {
//    defer func() {
//        fmt.Println(keys)
//    }()
    for _, s:= range t.unigram {
        k, e :=applyRule(s, index, c)
        if e == nil{
            keys = append(keys, k)
        }else {
            return keys, e
        }
    }
    return keys, nil
}

func (t *Template)ApplyRuleBigram(index int, c*Clique)(keys[]string, err error) {
    for _, s:= range t.bigram {
        k, e :=applyRule(s, index, c)
        if e == nil {
            keys = append(keys, k)
        }else {
            return keys, e
        }
    }
    return keys, nil
}


func applyRule(t string, index int, c *Clique) (k string, err error) {
    lt := len(t)
    for i:=0; i<lt; i++ {
        switch(t[i]) {
        default:
            k += string(t[i])
        case '%':
            switch t[i+1] {
            case 'x':
                i += 2
                re,l,e := getString(t[i:], index, c)
                if e != nil {
                    panic(e.Error())
                    return k, e
                }
                k += re
                i += l-1

            default:
                return k, StringError("parse template failed")
            }
        }
    }
    return k, nil
}

func getString(t string, index int, c *Clique)(re string, l int, e error) {
    lt := len(t)
    if(t[0] != '[') {
        return re, l, StringError("template format error 1")
    }
    var i = 1
    var neg int = 1
    if t[i] == '-' {
        neg = -1
        i += 1
    }

    row := 0
    col := 0
//    fmt.Println(t[i:])
    LABEL1:
    for ; i<lt; i++ {
        switch(t[i]) {
        case '1','2','3','4','5','6','7','8','9','0':
            row = row * 10 + int(t[i]-'0')
        case ',':
            i++
            break LABEL1
        default:
            return re, i, StringError("template format error 2")
        }
    }

    LABEL2:
    for ;i<lt;i++ {
        switch(t[i]) {
        case '1','2','3','4','5','6','7','8','9','0':
            col = col *10 + int(t[i]-'0')
        case ']':
            i++
            break LABEL2
        default:
            return re, i, StringError("template format error 3")
        }
    }

    row = row * neg

//    fmt.Println(row, KMaxContextSize)
    if row < (-1*KMaxContextSize) || row > KMaxContextSize {
        return re, i, StringError("template column too large")
    }

    if col > c.columns {
        return re, i, StringError("template columna more than train files")
    }

    idx := index + row
    if (idx < 0) {
        return BOS[-idx-1], i, nil
    }
    if (idx >= c.Size()) {
        return EOS[idx - c.Size()], i, nil
    }

    return c.lines[idx][col], i, nil
}

