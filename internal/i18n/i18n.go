package i18n

import "fmt"

var currentLang = "ja"

var translations = map[string]map[string]string{
	"ja": jaTranslations,
	"en": enTranslations,
}

// 表示言語を設定
func SetLanguage(lang string) {
	if _, ok := translations[lang]; ok {
		currentLang = lang
	}
}

// キーに対応する翻訳文字列を返す
func T(key string) string {
	if msgs, ok := translations[currentLang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// フォールバック: キーをそのまま返す
	return key
}

// フォーマット付き翻訳を返す
func TF(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}
