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
	"status.not_created": "未作成",
	"status.unknown":     "不明",

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
	"action.stop":    "停止",
	"action.restart": "再起動",
	"action.build":   "ビルド",
	"action.logs":    "ログ表示",
	"action.down":    "コンテナ削除",

	// 確認ダイアログ
	"confirm.title":   "確認",
	"confirm.up":      "%s を起動しますか？",
	"confirm.stop":    "%s を停止しますか？",
	"confirm.restart": "%s を再起動しますか？",
	"confirm.build":   "%s をビルドしますか？",
	"confirm.exec":    "%s のシェルに接続しますか？",
	"confirm.down":    "%s のコンテナを削除しますか？",
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
	"help.toggle":  "[h]折りたたみ",

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
	"exec.not_running":      "サービスが起動していません",
	"action.not_created":    "未作成のサービスにはこの操作を実行できません",

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
	"help.desc.toggle":         "プロジェクト折りたたみ切替",
	"help.desc.quit":           "終了",

	// イメージ管理
	"detail.tab.images":       "イメージ",
	"detail.tab.images.short": "Img",
	"images.title":            "Docker イメージ",
	"images.no_images":        "イメージがありません",
	"images.repo":             "リポジトリ:タグ",
	"images.size":             "サイズ",
	"images.created":          "作成日時",
	"images.status":           "状態",
	"images.used":             "使用中",
	"images.unused":           "未使用",
	"images.dangling":         "Dangling",
	"images.removing":         "イメージを削除中...",
	"images.removed":          "イメージを削除しました: %s",
	"images.remove_failed":    "イメージの削除に失敗しました: %s",
	"images.pruning":          "未使用イメージを削除中...",
	"images.pruned":           "イメージを削除しました: %s 解放",
	"images.prune_failed":     "イメージの一括削除に失敗しました: %s",
	"images.in_use":           "使用中のイメージは削除できません",

	// ボリューム管理
	"detail.tab.volumes":       "ボリューム",
	"detail.tab.volumes.short": "Vol",
	"volumes.title":            "Docker ボリューム",
	"volumes.no_volumes":       "ボリュームがありません",
	"volumes.name":             "名前",
	"volumes.driver":           "ドライバ",
	"volumes.mountpoint":       "マウントポイント",
	"volumes.status":           "状態",
	"volumes.used":             "使用中",
	"volumes.unused":           "未使用",
	"volumes.removing":         "ボリュームを削除中...",
	"volumes.removed":          "ボリュームを削除しました: %s",
	"volumes.remove_failed":    "ボリュームの削除に失敗しました: %s",
	"volumes.pruning":          "未使用ボリュームを削除中...",
	"volumes.pruned":           "ボリュームを削除しました: %s 解放",
	"volumes.prune_failed":     "ボリュームの一括削除に失敗しました: %s",
	"volumes.in_use":           "使用中のボリュームは削除できません",

	// 確認ダイアログ（イメージ・ボリューム）
	"confirm.remove_image":  "イメージ %s を削除しますか？",
	"confirm.remove_volume": "ボリューム %s を削除しますか？",
	"confirm.prune_images":  "未使用イメージを全て削除しますか？",
	"confirm.prune_volumes": "未使用ボリュームを全て削除しますか？",

	// ヘルプ（イメージ・ボリューム）
	"help.images":  "[I]イメージ",
	"help.volumes": "[V]ボリューム",
	"help.delete":  "[x]削除",
	"help.prune":   "[p]一括削除",
	"help.desc.images":        "イメージ一覧表示",
	"help.desc.volumes":       "ボリューム一覧表示",
	"help.desc.delete":        "選択項目を削除",
	"help.desc.prune":         "未使用項目を一括削除",
	"help.desc.tab_switch_all": "情報/ログ/環境変数/イメージ/ボリューム タブ切替",

	// パス解決
	"resolve.paths_updated": "パスを自動解決しました: %s",

	// エラー
	"error.docker_connect": "Docker への接続に失敗しました: %s",
	"error.compose_exec":   "docker compose コマンドの実行に失敗しました: %s",
	"error.log_stream":     "ログストリームのエラー: %s",

	// エラー診断 - 共通
	"diag.cause": "考えられる原因: %s",
	"diag.hints": "確認事項:",

	// Docker接続エラー診断
	"diag.conn.cause.not_running":      "Docker デーモンが起動していません",
	"diag.conn.cause.permission":       "Docker ソケットへのアクセス権限がありません",
	"diag.conn.cause.no_socket":        "Docker ソケットが見つかりません",
	"diag.conn.cause.timeout":          "Docker デーモンが応答しません",
	"diag.conn.cause.version_mismatch": "Docker クライアントとデーモンのバージョンが一致しません",
	"diag.conn.cause.unknown":          "原因を特定できません",

	"diag.conn.hint.start_docker":   "ご利用の Docker 環境を起動してください (Docker Desktop / colima / finch 等)",
	"diag.conn.hint.add_group":      "sudo usermod -aG docker $USER を実行してください",
	"diag.conn.hint.relogin":        "実行後、再ログインしてください",
	"diag.conn.hint.check_installed": "Docker がインストールされているか確認してください",
	"diag.conn.hint.check_load":     "Docker デーモンが過負荷でないか確認してください",
	"diag.conn.hint.check_resources": "システムリソース (CPU/メモリ) の状況を確認してください",
	"diag.conn.hint.update_docker":  "Docker を最新版にアップデートしてください",
	"diag.conn.hint.run_docker_info": "docker info を実行して状態を確認してください",

	// Compose実行エラー診断
	"diag.compose.cause.port_conflict":    "ポート競合",
	"diag.compose.cause.image_not_found":  "イメージ取得失敗",
	"diag.compose.cause.file_not_found":   "ファイルが見つかりません",
	"diag.compose.cause.network_not_found": "ネットワークが見つかりません",
	"diag.compose.cause.timeout":          "タイムアウト",
	"diag.compose.cause.no_config":        "compose 設定ファイルが見つかりません",

	"diag.compose.hint.check_port":         "該当ポートを使用中のプロセスを確認・停止してください",
	"diag.compose.hint.check_image":        "イメージ名を確認するか docker compose build を実行してください",
	"diag.compose.hint.docker_login":       "プライベートレジストリの場合は docker login を実行してください",
	"diag.compose.hint.check_path":         "Dockerfile やビルドコンテキストのパスを確認してください",
	"diag.compose.hint.create_network":     "docker network create で必要なネットワークを作成してください",
	"diag.compose.hint.check_network":      "ネットワーク接続を確認してください",
	"diag.compose.hint.check_daemon":       "Docker デーモンの状態を確認してください",
	"diag.compose.hint.check_compose_file": "作業ディレクトリに compose.yml が存在するか確認してください",
}
