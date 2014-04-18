package crfpp

import (
    "io"
    "bytes"
    "os"
    "fmt"
)

const MAX_LINE = 1024

type TAG_TYPE   int

type Tag struct {
    tags    map[string]TAG_TYPE
    id      TAG_TYPE
}


var DefaultTags = new(Tag)

func ascIsSpace(t rune)bool {
    if t == ' ' || t == '\t' {
        return true
    }
    return false
}

func (t *Tag)ReadTags(r io.ReadSeeker)(e error) {
    if t.tags == nil {
        t.tags = make(map[string]TAG_TYPE)
    }

    bf := make([]byte, 1024)
    size := 0
    for {
        n, e:= ReadLine(r, bf)
        if e!= nil && e!= io.EOF {
            return e
        }

        if n==0 {
            if e == io.EOF {
                break
            }else {
                continue
            }
        }

        sz := bytes.FieldsFunc(bf[:n], ascIsSpace)
        if len(sz) < 2 {
            return StringError("file format error, only 1 column")
        }

        if size==0 {
            size = len(sz)
        }else if size != len(sz) {
            return StringError(fmt.Sprintf("file format error, columns not same %d:%d:%s:%s",size, len(sz),sz, bf[0:n]))
        }

        tag := string(sz[len(sz) -1])
        
        val, isexist := t.tags[tag]
        if isexist {
            continue
        }

        val = t.id
        t.id += 1
        t.tags[tag] = val

        if e == io.EOF{
            break
        }
    }
    return nil
}
        
func (t *Tag) Find(tag string)(v TAG_TYPE, ok bool) {
    v, ok = t.tags[tag]
    return
}

//AddTag, cannot be used in multi-Threads
func (t *Tag) AddTag(tag string)(val TAG_TYPE, ok bool) {
    val, ok= t.tags[tag]
    if !ok {
        val = t.id
        t.id += 1 
        t.tags[tag] = val
        ok = true
    }
    return 
}

func (t *Tag)GetTagNum()(int) {
    return len(t.tags)
}

func (t *Tag)GetTags() map[string]TAG_TYPE {
    return t.tags
}

func ReadTagFromIO(r io.ReadSeeker)(e error) {
    return DefaultTags.ReadTags(r)
}

func AddTag(tag string)(val TAG_TYPE, ok bool) {
    return DefaultTags.AddTag(tag)
}
func FindTag(tag string)(TAG_TYPE, bool) {
    return DefaultTags.Find(tag)
}

func GetTagNum() (int) {
    return DefaultTags.GetTagNum()
}

func GetTags() map[string]TAG_TYPE {
    return DefaultTags.GetTags()
}

func ReadTagFromFile(name string)(*Tag ,error) {
    f, e:= os.Open(name)
    if e != nil {
        return nil,e
    }
    defer f.Close()
    t := new(Tag)
    e = t.ReadTags(f)
    return t, e
}



    
