package lang

import (
	"fmt"

	"github.com/zhitoo/golang-web-api/config"
)

var trans = map[string]string{
	"en.validation.default":  "This field is not valid",
	"en.validation.required": "This field is required",
	"en.validation.irmobile": "This field must be a valid Iranian mobile number",
	"en.validation.min":      "This field must have a minimum value of %s",

	"fa.validation.default":  "مقدار وارد شده معتبر نیست",
	"fa.validation.required": "پر کردن این فیلد الزامی است",
	"fa.validation.irmobile": "شماره موبایل وارد شده معتبر نیست",
	"fa.validation.min":      "حداقل مقدار باید %s باشد",
	"fa.validation.numeric":  "این فیلد باید عددی باشد",
	"fa.validation.alpha":    "این فیلد باید متنی باشد",
}

func Trans(scope string, key string, values ...any) string {
	lang := config.GetConfig().App.Lang

	if scope != "" {
		key = lang + "." + scope + "." + key
	}

	if val, ok := trans[key]; ok {
		return fmt.Sprintf(val, values...)
	}

	return key
}
