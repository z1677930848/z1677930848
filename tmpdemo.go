package main
import "fmt"

type Base struct{}
func (b *Base) Save(){fmt.Println("save")}

type Wrapper struct{ Base }

type AliasWrapper Base

type Derived Base

type Alias = Base

func main(){
    var w Wrapper
    w.Save()
    var aw AliasWrapper
    aw.Save()
    var d Derived
    _ = d
    var al Alias
    al.Save()
}
