package conf

var BaseUrl = "http://172.16.20.90:30116/apisix/admin"

func conf(baseUrl string){
	BaseUrl = baseUrl
}
