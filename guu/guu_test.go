package guu

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router{
	r := newRouter()
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name", nil)
	r.addRouter("GET", "/hello/b/c", nil)
	r.addRouter("GET", "/hi/:name", nil)
	r.addRouter("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T){
	//设计函数的测试用例
	ok := reflect.DeepEqual(parsePattern("/p/:name"),[]string{"p",":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"),[]string{"p","*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*path/:name"),[]string{"p","*path"})
	if ok == false{
		t.Fatal("test parsePattern failed")
	}
	fmt.Println("parsePattern pass")
}

func TestGetRoute(t *testing.T){
	r := newTestRouter()
	//if r == nil{
	//	t.Fatal("wrong")
	//}
	n, ps := r.getRoute("GET","/hello/guu")
	if n == nil{
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name"{
		t.Fatal("should match /hello/:name")
	}

	v,ok := ps["name"]
	if ok == false || v != "guu"{
		t.Fatal("name should be equal to 'guu'")
	}
	fmt.Printf("matched path：%s, params['name']: %s \n", n.pattern, ps["name"])
}