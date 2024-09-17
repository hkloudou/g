package gapi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator
var Validate *validator.Validate

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		Validate = v
	} else {
		panic("can't read `*validator.Validate` from gin engine")
	}
	// 注册自定义字段名称，提取结构体中的 label 标签
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		label := fld.Tag.Get("label")
		if label == "" {
			return fld.Name // 没有 label 则使用字段名
		}
		return label
	})

	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	// if !ok {
	// 	panic("not read translator zh")
	// }
	err := zh_translations.RegisterDefaultTranslations(Validate, trans)
	if err != nil {
		panic(err)
	}
}

func translateErrors(err error) map[string]string {
	errs := err.(validator.ValidationErrors)
	errors := map[string]string{}
	fmt.Println("tran", err.Error())
	for _, e := range errs {
		errors[e.Field()] = e.Translate(trans)
	}
	// validator.err
	return removeTopStruct(errors)
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		// 去掉结构体名前缀
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
