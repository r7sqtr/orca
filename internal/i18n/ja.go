package i18n

var jaTranslations = map[string]string{
	// アプリケーション
	"app.title":       "orca - Docker Compose マネージャー",
	"app.connecting":  "Docker に接続中...",
	"app.connected":   "Docker に接続しました",
	"app.no_docker":   "Docker に接続できません",
	"app.colima_hint": "colima start を実行してください",

	// サイドバー
	"sidebar.title":       "サービス一覧",
	"sidebar.no_projects": "Compose プロジェクトが見つかりません",

	// サービス状態
	"status.running":    "実行中",
	"status.exited":     "停止",
	"status.paused":     "一時停止",
	"status.restarting": "再起動中",
	"status.created":    "作成済み",
	"status.removing":   "削除中",
	"status.dead":       "異常停止",
	"status.unknown":    "不明",

	// 詳細パネル
	"detail.tab.info": "情報",
	"detail.tab.logs": "ログ",
	"detail.tab.env":  "環境変数",
	"detail.service":  "サービス",
	"detail.image":    "イメージ",
	"detail.state":    "状態",
	"detail.ports":    "ポート",
	"detail.health":   "ヘルス",
	"detail.id":       "コンテナID",
	"detail.created":  "作成日時",
	"detail.started":  "起動日時",

	// ログ
	"log.title":         "ログ",
	"log.follow":        "フォロー中",
	"log.paused":        "停止中",
	"log.no_logs":       "ログがありません",
	"log.searching":     "検索: %s",
	"log.matches":       "%d 件一致",
	"log.lines":         "%d 行",
	"log.stream.stdout": "stdout",
	"log.stream.stderr": "stderr",

	// 操作
	"action.up":      "起動",
	"action.down":    "停止",
	"action.restart": "再起動",
	"action.build":   "ビルド",
	"action.logs":    "ログ表示",

	// 確認ダイアログ
	"confirm.title":   "確認",
	"confirm.up":      "%s を起動しますか？",
	"confirm.down":    "%s を停止しますか？",
	"confirm.restart": "%s を再起動しますか？",
	"confirm.build":   "%s をビルドしますか？",
	"confirm.exec":    "%s のシェルに接続しますか？",
	"confirm.yes":     "はい",
	"confirm.no":      "いいえ",

	// ヘルプ
	"help.move":    "[j/k]移動",
	"help.up":      "[u]起動",
	"help.down":    "[d]停止",
	"help.restart": "[r]再起動",
	"help.search":  "[/]検索",
	"help.follow":  "[f]フォロー",
	"help.tab":       "[Tab]タブ切替",
	"help.focus":     "[C-h/l]パネル移動",
	"help.quit":    "[q]終了",
	"help.help":    "[?]ヘルプ",
	"help.enter":   "[Enter]選択",
	"help.esc":     "[Esc]戻る",
	"help.logs":    "[l]ログ",
	"help.info":    "[i]情報",
	"help.env":     "[v]環境変数",
	"help.exec":    "[e]シェル",
	"help.copy":    "[y]コピー",
	"help.export":  "[o]エクスポート",
	"help.build":   "[b]ビルド",

	// 検索
	"search.placeholder": "検索...",
	"search.no_results":  "一致する結果がありません",

	// ログ操作
	"log.copied":          "ログをクリップボードにコピーしました",
	"log.copy_failed":     "クリップボードへのコピーに失敗しました: %s",
	"log.exported":        "ログを保存しました: %s",
	"log.export_failed":   "ログの保存に失敗しました: %s",

	// 環境変数
	"env.title":    "環境変数",
	"env.no_env":   "環境変数がありません",
	"env.no_container": "コンテナが起動していません",

	// シェル接続
	"exec.not_running": "サービスが起動していません",

	// ヘルプオーバーレイ
	"help.overlay.title":       "キーバインド一覧",
	"help.overlay.sidebar":     "サイドバー",
	"help.overlay.detail":      "Detail パネル",
	"help.overlay.global":      "共通",
	"help.overlay.close":       "Esc/? で閉じる",
	"help.desc.move":           "上下移動",
	"help.desc.up":             "サービス起動",
	"help.desc.down":           "サービス停止",
	"help.desc.restart":        "サービス再起動",
	"help.desc.build":          "イメージビルド",
	"help.desc.exec":           "シェル接続",
	"help.desc.tab_switch":     "情報/ログ/環境変数タブ切替",
	"help.desc.panel_switch":   "パネル切替",
	"help.desc.follow":         "ログフォロー切替",
	"help.desc.search":         "ログ検索",
	"help.desc.copy":           "ログをクリップボードにコピー",
	"help.desc.export":         "ログをファイルに保存",
	"help.desc.back":           "サイドバーに戻る",
	"help.desc.help":           "このヘルプを表示",
	"help.desc.quit":           "終了",

	// エラー
	"error.docker_connect": "Docker への接続に失敗しました: %s",
	"error.compose_exec":   "docker compose コマンドの実行に失敗しました: %s",
	"error.log_stream":     "ログストリームのエラー: %s",
}
