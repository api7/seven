package conf

var BaseUrl = "http://172.16.20.90:30116/apisix/admin"
var UrlGroup = make(map[string]string)

func AddGroup(group string){
	UrlGroup[group] = "http://" + group + "/apisix/admin"
}

func FindUrl(group string) string {
	if UrlGroup[group] != ""{
		return UrlGroup[group]
	} else {
		return BaseUrl
	}
}
