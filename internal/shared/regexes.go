package shared

import "regexp"

type RegexesModel struct {
	VersionRegex *regexp.Regexp
	NameRegex    *regexp.Regexp
	UserRegex    *regexp.Regexp
}

var TrakRegexes RegexesModel = RegexesModel{
	VersionRegex: regexp.MustCompile(`^v?[0-9]+(\.[0-9]+)*.*$`),
	NameRegex:    regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`),
	UserRegex:    regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-]*$`),
}
