package crfpp

import (
//    "fmt"
)

type ID_TYPE    int64
type NUM_TYPE   int64

type feature struct {
    id  ID_TYPE
    num NUM_TYPE
}

type FeatureSet struct {
    dict    map[string]*feature
    maxid   ID_TYPE
    tagnum  int

}

func (f *FeatureSet)Size() ID_TYPE {
    return f.maxid
}


func GetFeatureSet(t int)(*FeatureSet) {
    fs := new(FeatureSet)
    fs.dict = make(map[string]*feature)
    fs.tagnum = t
    return fs
}


func (f *FeatureSet)FeatureID(key string, learnflag bool) (id ID_TYPE) {

    v, ok := f.dict[key]
    if learnflag {          //learn stage
        if !ok {
            f.dict[key] = &feature{f.maxid,NUM_TYPE(1)}
            id  = f.maxid
            if key[0] == 'U' {
                f.maxid += ID_TYPE(f.tagnum)
            }else {
                f.maxid += ID_TYPE(f.tagnum*f.tagnum)
            }
            return id
        } 
    
        v.num +=1
        return v.id
    }

    if !ok { //not found
        return ID_TYPE(-1)
    }
    return v.id
}

func (f *FeatureSet)Shrink(freq NUM_TYPE) (old2new map[ID_TYPE]ID_TYPE) {
//    fmt.Println("set is ", len(f.dict), f.maxid)
    var id ID_TYPE = 0
    old2new = make(map[ID_TYPE]ID_TYPE)

    for k, v:= range f.dict {
//        fmt.Println(k, v)
        if v.num < freq {
            delete(f.dict, k)
        }else {
            old2new[v.id] = id
            v.id = id
            if k[0] == 'U' {
                id += ID_TYPE(f.tagnum)
            }else {
                id += ID_TYPE(f.tagnum *f.tagnum)
            }
        }
    }
    f.maxid = id
    return
}

func (f *FeatureSet)formatFeature(in *featureMsgIn, learnflag bool) (out *featureMsgOut) {
    out = new(featureMsgOut)
    out.id = in.id
    out.col = in.col
    out.ty = in.ty
    keys := in.keys
    for _,t := range keys{
        id := f.FeatureID(t, learnflag)
        if id == ID_TYPE(-1) {
            continue
        }
        out.f = append(out.f, id)
    }
    return
}


func (f *FeatureSet)GenerateFeatureGo(in chan *featureMsgIn, out chan *featureMsgOut, learnflag bool) {
    var fmi *featureMsgIn
    var ok bool
    for {
        if fmi, ok = <- in; !ok {
            break
        }

        fmo := f.formatFeature(fmi, learnflag)
        out <- fmo
    }
    close(out)
}



