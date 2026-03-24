package validation

import "regexp"

var codePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
var permissionCodePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func IsValidCode(value string) bool {
	if value == "" {
		return false
	}
	return codePattern.MatchString(value)
}

func IsValidPermissionCode(value string) bool {
	if value == "" {
		return false
	}
	return permissionCodePattern.MatchString(value)
}
