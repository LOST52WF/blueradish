package tools

/*
func  isdemo(){
	var config = parseConfig.New("config.json")
	// 此为interface{}格式数据
	var name = config.Get("name")
	// 断言
	var nameString = name.(string)

	// 取数组
	var urls = config.Get("urls").([]interface{})
	var urlsString []string
	for _,v := range urls {
		urlsString = append(urlsString, v.(string))
	}

	// 取嵌套map内数据
	var qq = config.Get("info > qq").("string")
	var weixin = config.Get("info > weixin").("string")
}
*/