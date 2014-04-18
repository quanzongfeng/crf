package crfpp

import (
    "log"
    "io"
    "bytes"
    "fmt"
)

type StringError  string

func (r StringError)Error() string {
    return string(r)
}


//read line less than 1024 bytes
func ReadLine(r io.ReadSeeker, bf []byte)(n int, e error) {
    newline := make([]byte, 1024)
    n, e = r.Read(newline)
    if  e != nil && e != io.EOF {
        return 0, e
    }

    lineTag := bytes.IndexByte(newline[:n], '\n')
    if lineTag == -1 {
        if e != io.EOF {
            return 0, StringError(fmt.Sprintf("line too long than %d bytes", n)) 
        }

        if n > len(bf) {
            return 0, StringError(fmt.Sprintf("too small bytes to hold line"))
        }
        copy(bf, newline[:n])
        return n, io.EOF
    }

    if lineTag > len(bf) {
        return 0, StringError(fmt.Sprintf("too small bytes to hold line"))
    }

    back := lineTag -n +1
    _, err:=r.Seek(int64(back), 1)
    if err != nil {
        log.Fatal("Read error")
    }

    copy(bf, newline[:lineTag])
    return lineTag, e
}




