package i18n

import "fmt"

var currentLang = "ja"

var translations = map[string]map[string]string{
	"ja": jaTranslations,
	"en": enTranslations,
}

// SetLanguage は表示言語を設定する
func SetLanguage(lang string) {
	if _, ok := translations[lang]; ok {
		currentLang = lang
	}
}

// T はキーに対応する翻訳文字列を返す
func T(key string) string {
	if msgs, ok := translations[currentLang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// フォールバック: キーをそのまま返す
	return key
}

// TF はフォーマット付き翻訳を返す
func TF(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}
