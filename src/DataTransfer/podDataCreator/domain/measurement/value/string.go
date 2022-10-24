package value

import "fmt"

type String string
 
func (s *String) ToDisplayString() string {
	return fmt.Sprintf("%v",*s)
}
