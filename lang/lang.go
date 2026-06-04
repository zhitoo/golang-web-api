package lang

import (
	"fmt"

	"github.com/zhitoo/golang-web-api/config"
)

var translations = map[string]map[string]map[string]string{
	"fa": fa,
	"en": en,
}

var Lang = "en"

func init() {
	cfg := config.GetConfig()
	Lang = cfg.App.Lang
}

func Trans(scope, key string, values ...any) string {

	if len(values) > 0 {
		fieldName := values[0].(string)

		values[0] = Trans("feild", fieldName)
	}

	if scopes, ok := translations[Lang]; ok {
		if messages, ok := scopes[scope]; ok {
			if msg, ok := messages[key]; ok {
				return fmt.Sprintf(msg, values...)
			}
		}
	}

	return key
}
