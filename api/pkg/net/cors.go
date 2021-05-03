package net

var Headers = map[string]map[string]string{
	"dev": {
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
	},
	"prod": {
		"Access-Control-Allow-Origin":  "https://bscfees.com",
		"Access-Control-Allow-Methods": "GET OPTIONS",
		"Access-Control-Allow-Headers": "Accept-Content",
	},
}
// SetupCORSHeaders specifies CORS related headers based on the Origin of the request. There should only 2 domains
// allowed: https://bscfees.com and https://www.bscfees.com
// It does what the infrastructure is supposed to do. This could be bypassed by a simple http redirection.
func SetupCORSHeaders(env string, originHeader string) map[string]string {
	if env == "dev" {
		return Headers["dev"]
	}

	prodHeaders := Headers["prod"]
	if originHeader != "https://www.bscfees.com" && originHeader != "https://bscfees.com" {
		return prodHeaders
	}

	prodHeaders["Access-Control-Allow-Origin"] = originHeader
	return prodHeaders
}
