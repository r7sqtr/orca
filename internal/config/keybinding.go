package config

// DefaultKeyBindings はデフォルトのキーバインドを返す
func DefaultKeyBindings() map[string]string {
	return map[string]string{
		"up":           "k",
		"down":         "j",
		"select":       "enter",
		"back":         "esc",
		"quit":         "q",
		"tab":          "tab",
		"service_up":   "u",
		"service_down": "d",
		"restart":      "r",
		"search":       "/",
		"follow":       "f",
		"logs":         "l",
		"info":         "i",
		"help":         "?",
	}
}
