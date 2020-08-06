package encode

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)
func setup() {
	fmt.Println("Before all tests")
}

func teardown() {
	fmt.Println("After all tests")
}

func Mul(a,b int)int  {
	return a*b
}
//  calc_test.go
func TestMul(t *testing.T) {
	cases := []struct {
		Name           string
		A, B, Expected int
	}{
		{"pos", 2, 3, 6},
		{"neg", 2, -3, -6},
		{"zero", 2, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := Mul(c.A, c.B); ans != c.Expected {
				t.Fatalf("%d * %d expected %d, but %d got",
					c.A, c.B, c.Expected, ans)
			}
		})
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

//mock
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用


}

func TestSpec(t *testing.T) {
	assert.Equal(t,1,1,"err")

	// Only pass t into top-level Convey calls
	Convey("Given some integer with a starting value", t, func() {
		x := 1
		So(1,ShouldEqual,1)
		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
				So(true,ShouldBeTrue)
			})
		})
	})
}