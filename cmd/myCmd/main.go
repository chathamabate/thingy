package main

import "fmt"

type I interface {
    Mutate()
}

type X struct {
    Val int
}

func (x *X) Mutate() {
    x.Val++
}


func main() {
    a := X{Val:10}

    var i I

    i = &a
    j := i

    i.Mutate()
    j.Mutate()

    fmt.Println(a)
}
