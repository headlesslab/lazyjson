package lazyjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/headlesslab/lazyjson"
)

func ExampleJSON() {
	obj := lazyjson.NewFrom(`{"a": {"b": [1, 2]}}`)

	fmt.Println(obj.Get("a.b.0").Int())

	obj.Set("a.b.1", "ok").Set("c", 2)
	obj.Del("c")
	fmt.Println(">", obj.JSON("> ", "  "))

	// Output:
	// 1
	// > {
	// >   "a": {
	// >     "b": [
	// >       1,
	// >       "ok"
	// >     ]
	// >   }
	// > }
}

func Test(t *testing.T) {
	eq := genEq(t)

	eq(lazyjson.NewFrom("true").Bool(), true)
	eq(lazyjson.New([]byte("10")).Int(), 10)
	eq(lazyjson.New(10).Int(), 10)
	eq(lazyjson.New(lazyjson.New(10)).Int(), 10)
	eq(lazyjson.JSON{}.Int(), 0)
	eq(lazyjson.JSON{}.JSON("", ""), "null")
	eq(lazyjson.New(nil).Num(), 0.0)
	eq(lazyjson.New(nil).Bool(), false)

	buf := bytes.NewBufferString("10")
	fromBuf := lazyjson.New(buf)
	eq(fromBuf.Int(), 10)
	eq(fromBuf.Int(), 10)

	b, _ := n(`"ok"`).MarshalJSON()
	eq(string(b), `"ok"`)

	eq(lazyjson.JSON{}.Raw(), nil)
	eq(lazyjson.New([]byte(`ok"`)).Raw(), []byte(`ok"`))
	eq(n(`"ok"`).Str(), "ok")
	eq(n(`1`).Str(), "1")
	eq(n(`1.2`).Num(), 1.2)
	eq(n(`"ok"`).Num(), 0.0)
	eq(n(`1`).Int(), 1)
	eq(n(`"ok"`).Int(), 0)
	eq(n(`true`).Bool(), true)
	eq(n(`1`).Bool(), false)
	eq(n(`null`).Nil(), true)

	eq(n(`1`).Join(""), "")

	j := n(`{
		"a": {
			"b": 1
		},
		"c": ["x", "y", "z"]
	}`)

	eq(j.Get("a.b").Int(), 1)
	eq(j.Get("c.1").Str(), "y")

	v, _ := j.Gets("c", lazyjson.Query(func(i interface{}) (val interface{}, has bool) {
		return i.([]interface{})[1], true
	}))
	eq(v.Str(), "y")

	eq(j.Get("c").Arr()[1].Str(), "y")
	eq(lazyjson.New([]int{1, 2}).Arr()[1].Int(), 2)
	eq(j.Get("c").Join(" "), "x y z")
	eq(j.Get("a").Map()["b"].Int(), 1)
	eq(lazyjson.New(map[string]int{"a": 1}).Map()["a"].Int(), 1)
	eq(len(j.Get("c").Map()), 0)

	eq(lazyjson.New([]lazyjson.JSON{
		lazyjson.New(1),
		lazyjson.New(map[string]int{"a": 2}),
	}).Get("1.a").Int(), 2)

	v, _ = lazyjson.New(map[float64]int{2: 3}).Gets(2.0)
	eq(v.Int(), 3)

	_, has := j.Gets(true)
	eq(has, false)

	eq(j.Has("a.b"), true)
	eq(j.Has("a.x"), false)
	eq(j.Has("c.10"), false)

	onNil := lazyjson.JSON{}
	eq(onNil.Set("a.b", 10).Get("a.b").Int(), 10)

	self := lazyjson.New(nil)
	self.Set("1", "a")
	self.Set("a.b.1", 1)
	self.Sets("ok")
	eq(self.Str(), "ok")
	self.Sets(map[string]int{"a": 1})
	eq(self.Get("a").Int(), 1)
	self.Sets([]int{1})
	eq(self.Get("0").Int(), 1)

	j.Sets(2.0, "a", "b")
	eq(j.Get("a.b").Int(), 2)

	j.Set("c.1", 2)
	eq(j.Get("c.1").Int(), 2)

	j.Sets(3, "a", "b")
	eq(j.Get("a.b").Int(), 3)

	eq(j.Get("s.10.b").Nil(), true)
	eq(j.Get("c.10").Nil(), true)

	j.Set("s.1.a", 10)
	j.Set("c.5", "ok")
	eq(fmt.Sprint(j), `map[a:map[b:3] c:[x 2 z <nil> <nil> ok] s:[<nil> map[a:10]]]`)

	eq(j.Dels("s", 1, "a"), true)
	eq(fmt.Sprint(j), `map[a:map[b:3] c:[x 2 z <nil> <nil> ok] s:[<nil> map[]]]`)
	eq(j.Dels("s", 1), true)
	eq(fmt.Sprint(j), `map[a:map[b:3] c:[x 2 z <nil> <nil> ok] s:[<nil>]]`)
	j.Del("c.1")
	eq(fmt.Sprint(j), `map[a:map[b:3] c:[x z <nil> <nil> ok] s:[<nil>]]`)
	eq(j.Dels("c", 10), false)
	eq(j.Dels("c", "1"), false)
	eq(j.Dels("xxx", "1"), false)
	eq(j.Dels(1), false)
	d := lazyjson.New(1)
	d.Dels()
	eq(d.Val(), nil)
}

func TestUnmarshal(t *testing.T) {
	eq := genEq(t)

	v := struct {
		A int `json:"a"`
	}{}

	g := lazyjson.New([]byte(`{"a":1}`))

	err := g.Unmarshal(&v)
	if err != nil {
		t.Fatal(err)
	}

	if v.A != 1 {
		t.Fatal("parse error")
	}

	g.Get("a")
	eq(g.Unmarshal(&v).Error(), "lazyjson: value has been parsed")

	eq(lazyjson.JSON{}.Unmarshal(&v).Error(), "lazyjson: no value to unmarshal")
}

func TestConvertors(t *testing.T) {
	eq := genEq(t)

	n := 1.2
	i := 1
	s := "ok"
	b := true

	eq(lazyjson.Num(n), &n)
	eq(lazyjson.Int(i), &i)
	eq(lazyjson.Str(s), &s)
	eq(lazyjson.Bool(b), &b)
}

func TestLab(t *testing.T) {
}

func n(s string) (j lazyjson.JSON) {
	_ = json.Unmarshal([]byte(s), &j)
	return
}

func genEq(t *testing.T) func(a, b interface{}) {
	return func(a, b interface{}) {
		t.Helper()
		if !reflect.DeepEqual(a, b) {
			t.Log(a, "!=", b)
			t.Fail()
		}
	}
}
