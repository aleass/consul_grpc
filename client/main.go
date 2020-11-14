package main

import (
	"fmt"
	"strings"
)

func init() {
	str := "add_fa_qs"
	f := str[0 : 1]
	var StrType string
	ReStrType := "D"
	if strings.Index(str,"_") > -1 {
		StrType = "L"
	}else if strings.ToUpper(f) == f {
		StrType = "D"
	}else {
		StrType = "X"
	}
	var strs string
	i := 0


	for _,v := range []byte(str){
		if StrType == "L"  {
			if v == 95 {
				i = 0
			}else {
				if i == 0{
					if ReStrType != "L"{
						strs += strings.ToUpper(string(v))
						i = 1
					}else  {

					}
				}else {
					strs += string(v)
				}
			}
		}else if StrType == "U" {
			fmt.Println("a")
		} else {
			fmt.Println("a")
		}
	}
	fmt.Println(strs)
}



func main() {
	//route.RouteInit()
	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
}