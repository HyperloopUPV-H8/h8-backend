```go

type Vector2 struct {
    x int
    y int
}

func Print(vector Vector2) {
    fmt.Println(x, y)
}


type Vector3 struct {
    x int
    y int
    z int
}

type Stringer interface {
    String() string
}

func (vector Vector2) String() string {
    return fmt.Sprintf("%d, %d", vector.x, vector.y)
}

func (vector Vector3) String() string {
    return fmt.Sprintf("%d, %d", vector.x, vector.y, vector.z)
}

func Print(s Stringer) {
    fmt.Println(s.String())
}


```
