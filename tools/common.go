package tools

import (
	"strconv"
	"strings"
)

func StringToIntArray(str string)[]int{
	strings_id := strings.Split(str,",")
	id := []int{}
	for _,val := range strings_id {
		id_for_int,err := strconv.Atoi(val)
		if err == nil{  //若果转换出错，将忽略
			id = append(id,id_for_int)
		}
	}
	return id
}

func ReturnAllError(args ...error) string {
	err_string := make([]string,len(args))
	return_str := ""
	for i,errval := range args {
		if errval != nil{
			err_string[i] = errval.Error()
		}else{
			err_string[i] = "nil"
		}
	}
	return_str  = strings.Join(err_string,";")
	return return_str
}