package crfpp

import (
    "io"
    "os"
    "fmt"
    "strings"
    "strconv"
    "bytes"
)

type Model struct {
    version     string
    size        int64
    alpha       []float64
    set         *FeatureSet
    tags        *Tag
    template    *Template
}

func SaveModel(gc *Group, modelFile string)error {
    ml := new(Model)
    alpha := gc.alpha
    ml.size = int64(gc.Set.Size())
    if ml.size != int64(len(alpha)) {
        fmt.Println("size not equal")
        return StringError("alpha and dict size not euqal")
    }

    ml.alpha = alpha
    ml.set = gc.Set
    ml.template = gc.template
    ml.tags = gc.tags

    pa, e := os.Getwd()
    bi := strings.Index(pa,"crf")
    if bi == -1 {
        ml.version = "1.0"
    }else {
        be := strings.Index(pa[bi:], "/")
        if be == -1 {
            if len(pa) == bi+3 {
                ml.version = "1.0"
            }  else {
                ml.version = string(pa[bi+3:])
            }
        }else {
            fmt.Println(bi, be)
            bv := pa[bi+3:bi+be]
            if bv == ""{
                ml.version = "1.0"
            }else {
                ml.version = string(bv[3:])
            }
        }
    }

    f,er := os.OpenFile(modelFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if er != nil {
        return er
    }

    //write version
    f.WriteString(ml.version + "\n")
    //write template num: unigram num and bigram num
    f.WriteString(strconv.Itoa(len(ml.template.unigram)) + "\t" +strconv.Itoa(len(ml.template.bigram))+"\n")

    //write unigram
    for _,lu := range ml.template.unigram {
        f.WriteString(lu + "\n")
    }
    //write bigram
    for _,lb := range ml.template.bigram {
        f.WriteString(lb + "\n")
    }


    //tag type numbers
    lt := ml.tags.GetTagNum()
    f.WriteString(strconv.Itoa(lt)+"\n")

    //write tags and id 
    for k, v := range ml.tags.GetTags() {
        f.WriteString(k+"\t"+strconv.Itoa(int(v))+"\n")
    }

    //write alpha len
    la := len(alpha)
    f.WriteString(strconv.Itoa(la)+ "\n")

    for i, a:= range alpha{
        f.WriteString(strconv.FormatFloat(a, 'g', 10,64)+"\t")
        if (i+1)%10 == 0 {
            f.WriteString("\n")
        }
    }
    f.WriteString("\n")

    lf  := len(ml.set.dict)
    f.WriteString(strconv.Itoa(lf)+"\n")
    for k, v := range ml.set.dict {
        f.WriteString(k+"\t"+strconv.FormatInt(int64(v.id), 10)+"\n")
    }

    f.Close()

    return e
}

func  LoadModel(modelfile string) (*Model, error) {
    f, e:= os.Open(modelfile)
    if e!= nil {
        return nil, e
    }

    ml := new(Model)
    bf := make([]byte, 1024)
    n, err:= ReadLine(f, bf)
    if err != nil && err != io.EOF {
        return nil, err
    }
    ml.version = string(bf[:n])

    ml.template = new(Template)
    n, err =ReadLine(f, bf)
    if err != nil && err != io.EOF{
        return nil, err
    }
    listb :=bytes.FieldsFunc(bf[:n], ascIsSpace)
    if len(listb) < 2 {
        return nil, StringError("template nums error in modle file")
    }

    fmt.Println(string(listb[0]), string(listb[1]))

    un ,e0:= strconv.Atoi(string(listb[0]))
    bn ,e1:= strconv.Atoi(string(listb[1]))
    if e0 != nil{
        return nil, e0 
    }
    if e1 != nil {
        return nil ,e1
    }
    fmt.Println(listb, len(listb[0]), len(listb[1]))

    for i:= 0; i< un;i++ {
        n, err =ReadLine(f, bf)
        if err != nil && err != io.EOF{
            return nil, err
        }
        ml.template.unigram = append(ml.template.unigram, string(bf[:n]))
    }

    for i:=0; i<bn;i++ {
        n, err =ReadLine(f, bf)
        if err != nil && err != io.EOF{
            return nil, err
        }
        ml.template.bigram = append(ml.template.bigram, string(bf[:n]))
    }


    n, err =ReadLine(f, bf)
    if err != nil && err != io.EOF{
        return nil, err
    }
    tagnum, e2:= strconv.Atoi(string(bf[:n]))
    if e2!= nil {
        return nil,e2
    }
    ml.tags = new(Tag)
    ml.tags.tags = make(map[string]TAG_TYPE)
    for i:=0; i<tagnum; i++ {
        n, err =ReadLine(f, bf)
        if err != nil && err != io.EOF{
            return nil, err
        }
        listtags :=bytes.FieldsFunc(bf[:n], ascIsSpace)
        if len(listtags) < 2{
            return nil, StringError("read tag types failed")
        }
        tid, e:= strconv.Atoi(string(listtags[1]))
        if e != nil {
            return nil, e
        }
        ml.tags.tags[string(listtags[0])] = TAG_TYPE(tid)
    }


    n, err =ReadLine(f, bf)
    if err != nil && err != io.EOF{
        return nil, err
    }
    la, e3:= strconv.Atoi(string(bf[:n]))
    if e3!= nil {
        return nil,e3
    }
    fmt.Println(la)
    ml.alpha = make([]float64, 0, la)    

    i := 0
    for i=0; i<la;{
        n, err =ReadLine(f, bf)
        if err != nil && err != io.EOF{
            return nil, err
        }
        listalpha:=bytes.FieldsFunc(bf[:n], ascIsSpace)
        lla := len(listalpha)
        fmt.Println(listalpha, lla)
        for j:=0; j<lla; j++ {
            sf, e := strconv.ParseFloat(string(listalpha[j]),64)
            if e != nil {
                return nil, e
            }
            ml.alpha = append(ml.alpha, sf)
        }
        
        i += lla
    }

    if i != la {
        return nil, StringError("alpha size error")
    }

    n, err =ReadLine(f, bf)
    if err != nil && err != io.EOF{
        return nil, err
    }
    ls, e4:= strconv.Atoi(string(bf[:n]))
    if e4!= nil {
        return nil,e3
    }

    ml.set = GetFeatureSet(ml.tags.GetTagNum())
    for i=0; i< ls; i++ {
        n, err =ReadLine(f, bf)
        if err != nil && err != io.EOF{
            return nil, err
        }
        listalpha:=bytes.FieldsFunc(bf[:n], ascIsSpace)
        if len(listalpha) != 2 {
            return nil, StringError("model format error")
        }
        sid, e:= strconv.ParseInt(string(listalpha[1]), 10, 64)
        if e!= nil {
            return nil, e
        }
        ml.set.dict[string(listalpha[0])] = &feature{id:ID_TYPE(sid)}
    }

    f.Close()
    return ml, nil  
}
