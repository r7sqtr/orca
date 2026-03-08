# Orca

Docker Compose プロジェクトをターミナルから管理できる TUIツール

## 特徴
- Docker Compose プロジェクトのサービス一覧・状態をリアルタイム表示
- サービスの起動・停止・再起動をキーボード操作で実行
- コンテナログのリアルタイムフォロー・検索・エクスポート
- コンテナへのシェルアクセス
- 環境変数の確認
- イメージのビルド
- Vim ライクなキーバインド
- ダーク / ライトテーマ対応
- 設定ファイルによるカスタマイズ

## 要件
- [Docker](https://docs.docker.com/get-docker/) および [Docker Compose](https://docs.docker.com/compose/)
- Go 1.25.5 以上

## インストール
### Homebrew

```bash
brew install r7sqtr/tap/orca
```

### ソースからビルド
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
Docker Compose プロジェクトのディレクトリで実行します。

```bash
cd your-compose-project
orca
```

## キーバインド

| キー | 操作 |
|------|------|
| `k` / `↑` | 上へ移動 |
| `j` / `↓` | 下へ移動 |
| `Enter` | 選択 |
| `Esc` | 戻る |
| `Tab` | パネル切替 |
| `u` | サービス起動 |
| `d` | サービス停止 |
| `r` | サービス再起動 |
| `l` | ログ表示 |
| `f` | ログフォロー |
| `/` | 検索 |
| `i` | サービス情報 |
| `e` | シェル（exec） |
| `v` | 環境変数 |
| `b` | イメージビルド |
| `y` | コピー |
| `o` | エクスポート |
| `?` | ヘルプ |
| `q` | 終了 |

## 設定

設定ファイルは `~/.config/orca/config.yml` に配置します。

```yaml
# 言語 (デフォルト: "ja")
language: ja

# テーマ: "dark", "light", "auto" (デフォルト: "dark")
theme: dark

# ログバッファサイズ (デフォルト: 10000)
log_buffer_size: 10000

# Docker ホスト (未設定時は環境変数 DOCKER_HOST を使用)
docker_host: ""

# サイドバー幅のパーセント (0で自動) (デフォルト: 0)
sidebar_width: 0

# 破壊的操作の確認ダイアログ (デフォルト: true)
confirm_actions: true

# キーバインドのカスタマイズ
keybindings:
  up: k
  down: j
```

