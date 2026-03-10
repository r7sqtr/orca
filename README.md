# Orca
Docker Compose プロジェクトをターミナルから管理できる TUIツール

## 特徴
- Docker Compose プロジェクトのサービス一覧・状態をリアルタイム表示
- サービスの起動・停止・再起動をキーボード操作で実行
- コンテナログのリアルタイムフォロー・検索・エクスポート
- コンテナへのシェルアクセス
- 環境変数の確認
- イメージのビルド
- プロジェクトレジストリによる自動検出（一度認識したプロジェクトを記憶）
- サイドバーのプロジェクト折りたたみ
- 停止中・未作成のサービスも表示
- Vim ライクなキーバインド
- ダーク / ライトテーマ対応
- 設定ファイルによるカスタマイズ

## 要件
- Dockerデーモン（[Colima](https://github.com/abiosoft/colima)、[Finch](https://github.com/runfinch/finch) など）

## インストール
### Homebrew

```bash
brew install r7sqtr/tap/orca-tui
```

### ソースからビルド
Go 1.25.5 以上が必要です。

```bash
git clone https://github.com/r7sqtr/orca.git
cd orca
make install
```

デフォルトでは `/usr/local/bin` にインストールされます。インストール先を変更するには `PREFIX` を指定してください。

```bash
make install PREFIX=$HOME/.local
```

## 使い方

```bash
orca
```

任意のディレクトリから実行できます。
一度認識した Docker Compose プロジェクトはレジストリ（`~/.config/orca/registry.json`）に記憶され、次回以降も自動的に表示されます。

## キーバインド

| キー | 操作 |
|------|------|
| `k` / `↑` | 上へ移動 |
| `j` / `↓` | 下へ移動 |
| `Enter` | 選択 |
| `Esc` | 戻る |
| `Tab` | パネル切替 |
| `Ctrl+h` | 左パネルフォーカス |
| `Ctrl+l` | 右パネルフォーカス |
| `u` | サービス起動（`docker compose up -d`） |
| `d` | サービス停止（`docker compose stop`） |
| `r` | サービス再起動（`docker compose restart`） |
| `b` | イメージビルド（`docker compose build`） |
| `e` | シェル接続（`docker compose exec`） |
| `l` | ログ表示 |
| `f` | ログフォロー切替 |
| `/` | ログ検索 |
| `i` | サービス情報 |
| `v` | 環境変数 |
| `h` | プロジェクト折りたたみ切替 |
| `y` | ログコピー |
| `o` | ログエクスポート |
| `?` | ヘルプ |
| `q` | 終了 |

## 設定

設定ファイルは `~/.config/orca/config.yml` に配置します。初回起動時にデフォルト設定ファイルが自動生成されます。

```yaml
# Orca 設定ファイル

# 言語設定: "ja" (日本語), "en" (English)
language: ja

# テーマ: "dark", "light", "auto"
theme: dark

# ログバッファサイズ (1〜100000)
log_buffer_size: 10000

# Docker ホスト (未設定時は環境変数 DOCKER_HOST または自動検出)
# docker_host: ""

# サイドバー幅 (パーセント, 0で自動)
sidebar_width: 0

# 操作の確認ダイアログ
confirm_actions:
  exec: true
  up: true
  stop: true
  restart: true
  build: true

# キーバインドのカスタマイズ
# 各キーにはデフォルト値が設定されています
# keybindings:
#   up: k          # 上へ移動
#   down: j        # 下へ移動
#   select: enter  # 選択
#   back: esc      # 戻る
#   quit: q        # 終了
#   tab: tab       # パネル切替
#   focus_left: ctrl+h   # 左パネルフォーカス
#   focus_right: ctrl+l  # 右パネルフォーカス
#   start: u       # サービス起動
#   stop: d        # サービス停止
#   restart: r     # サービス再起動
#   search: /      # 検索
#   follow: f      # ログフォロー
#   logs: l        # ログ表示
#   info: i        # サービス情報
#   help: "?"      # ヘルプ
#   exec: e        # シェル (exec)
#   copy: y        # コピー
#   export: o      # エクスポート
#   env_vars: v    # 環境変数
#   build: b       # イメージビルド
```

