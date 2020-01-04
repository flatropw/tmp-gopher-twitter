package utils

var publicRoutes = []string{"/", "/api/v1/users/login", "/api/v1/users/register"}

func IsPublicRoute(requestPath string) bool {
	for _, v := range publicRoutes {
		if v == requestPath {
			return true
		}
	}
	return false
}
