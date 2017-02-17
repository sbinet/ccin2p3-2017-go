// package badfmt shows a badly formatted code.
// Its content will be automatically reformatted by 'gofmt'.
//
// Run like so:
//  $> gofmt -d ./badfmt
//
// to display the diff.
// 
// Or, like so:
//  $> gofmt -w ./badfmt
//
// to rewrite in place.
//
// Also, one can use the -s switch to simplify the code:
//  $> gofmt -d -s ./badfmt
package badfmt

import ("fmt"; "io")
func Bad(x,y,z float64)   float64{
return x+y   +
z;}

func myPrintf(format string,args ...interface{}){fmt.Printf(format,args...)}
func myWrite(w io.Writer,  data[]byte)(int,error){
	return w.Write(data,
)}

type Point struct{ X,Y float64}

func makePointMap() map[string]Point{
return map[string]Point{
"Paris": Point{X:0,Y:0},
"La Rochelle": Point{X:42,Y:42},
}
}
