package utils

var unauthorizedRoutes = []string{"/", "/api/v1/login", "/api/v1/register"}
func IsAuthorizedRoute(requestPath string) bool {
	for _, v := range unauthorizedRoutes {
		if v == requestPath {
			return false
		}
	}
	return true
}
