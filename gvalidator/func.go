package gvalidator

import (
	"regexp"
	"sync"
)

var datas = make(map[string]func() *regexp.Regexp)

func GetRegexp() map[string]func() *regexp.Regexp {
	return datas
}

func registe(key string, pattern string) {
	datas[key] = lazyRegexCompile(pattern)
}

func lazyRegexCompile(str string) func() *regexp.Regexp {
	var regex *regexp.Regexp
	var once sync.Once
	return func() *regexp.Regexp {
		once.Do(func() {
			regex = regexp.MustCompile(str)
		})
		return regex
	}
}
