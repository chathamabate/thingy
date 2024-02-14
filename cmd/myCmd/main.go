package main

import "fmt"


type X struct {
    Val int
}

func (x *X) Mutate() {
    x.Val++
}

type Y struct {
    X
}
 
func (y *Y) Mutate() {
    y.Val += 2
}


func main() {
    a := Y{
        X: X{
            Val: 10,
        },
    }

    a.X.Mutate()

    fmt.Println(a)
}
