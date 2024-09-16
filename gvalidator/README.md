## Install
``` sh
go get github.com/hkloudou/g/gvalidator
```


## work with gin
``` go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hkloudou/g/gvalidator"
)

var validate *validator.Validate

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
	} else {
		panic("can't read `*validator.Validate` from gin engine")
	}
	for k, v := range gvalidator.GetRegexp() {
		re := v()
		validate.RegisterValidation(k, func(fl validator.FieldLevel) bool {
			return re.MatchString(fl.Field().String())
		})
	}
}
```