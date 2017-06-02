package main

import "fmt"

//Pcall2 acts as a protected call, returning wether the call went through successfully, and its return value.
func Pcall2(f func([]interface{}) []interface{}, params []interface{}) (success bool, result []interface{}) {
	defer func() {
		if r := recover(); r != nil {
			success = false
			result = make([]interface{}, 1)
			result[0] = r
			fmt.Printf("An error occoured: %v", r)
		}
	}()
	return true, f(params)
}

func main2() {
	success, response := Pcall2(func(arg1 []interface{}) []interface{} {
		var i int
		for i = 0; i < 4; i++ {
			arg1[i] = i
		}
		return arg1
	}, make([]interface{}, 3))
	fmt.Printf("Success: %t, Response: %v", success, response)
}
