package parsetime_test

import (
	"fmt"
	"github.com/richardliao/parsetime"
)

func Example() {
	fmt.Println(parsetime.Parse("2006-01-02T15:04:05.999999999+08:00"))
	fmt.Println(parsetime.Parse("2006-01-02T15:04:05.999999+08:00"))
	fmt.Println(parsetime.Parse("2006-01-02T15:04:05.999+08:00"))
	fmt.Println(parsetime.Parse("2006-01-02T15:04:05+08:00"))
	fmt.Println(parsetime.Parse("2006-01-02T15:04:05+0800"))
	fmt.Println(parsetime.Parse("2006-01-02T15:04:05Z"))
	fmt.Println(parsetime.Parse("2006-01-02T15:04:05"))
	fmt.Println(parsetime.Parse("2006-01-02 15:04:05"))
	fmt.Println(parsetime.Parse("2006-01-02"))

	// Output:
	// 2006-01-02 15:04:05.999999999 +0800 CST <nil>
	// 2006-01-02 15:04:05.999999 +0800 CST <nil>
	// 2006-01-02 15:04:05.999 +0800 CST <nil>
	// 2006-01-02 15:04:05 +0800 CST <nil>
	// 2006-01-02 15:04:05 +0800 CST <nil>
	// 2006-01-02 23:04:05 +0800 CST <nil>
	// 2006-01-02 23:04:05 +0800 CST <nil>
	// 2006-01-02 23:04:05 +0800 CST <nil>
	// 2006-01-02 08:00:00 +0800 CST <nil>
}
