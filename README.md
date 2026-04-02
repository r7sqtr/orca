  <h1 align="center">Orca</h1>

  <p align="center">
    <strong>Docker Compose TUI管理ツール</strong>
  </p>

  <p align="center">
    <a href="#インストール">インストール</a> •
    <a href="#使い方">使い方</a> •
    <a href="#キーバインド">キーバインド</a> •
    <a href="#設定">設定</a>
  </p>

  ---

  ## 特徴

  **サービス管理**
  - サービスの起動・停止・再起動をキーボード操作で実行
  - イメージのビルド
  - コンテナへのシェルアクセス
  - 環境変数の確認

  **モニタリング**
  - サービス一覧・状態をリアルタイム表示
  - 停止中・未作成のサービスも表示
  - コンテナログのリアルタイムフォロー・検索・エクスポート

  **ユーザビリティ**
  - プロジェクトレジストリによる自動検出（一度認識したプロジェクトを記憶）
  - サイドバーのプロジェクト折りたたみ
  - Vim ライクなキーバインド
  - ダーク / ライトテーマ対応
  - 設定ファイルによるカスタマイズ

  ## 要件

  - Docker デーモン（[Colima](https://github.com/abiosoft/colima)、[Finch](https://github.com/runfinch/finch) など）

  ## インストール

  ### Homebrew

  ```bash
  brew install r7sqtr/tap/orca-tui
  ```

  ### ソースからビルド

  > Go 1.25.5 以上が必要です。

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

  ### ナビゲーション

  | キー | 操作 |
  |------|------|
  | `k` / `↑` | 上へ移動 |
  | `j` / `↓` | 下へ移動 |
  | `Enter` | 選択 |
  | `Esc` | 戻る |
  | `Tab` | パネル切替 |
  | `Ctrl+h` | 左パネルフォーカス |
  | `Ctrl+l` | 右パネルフォーカス |
  | `h` | プロジェクト折りたたみ切替 |

  ### サービス操作

  | キー | 操作 |
  |------|------|
  | `u` | 起動（`docker compose up -d`） |
  | `d` | 停止（`docker compose stop`） |
  | `r` | 再起動（`docker compose restart`） |
  | `b` | ビルド（`docker compose build`） |
  | `e` | シェル接続（`docker compose exec`） |

  ### ログ・情報

  | キー | 操作 |
  |------|------|
  | `l` | ログ表示 |
  | `f` | ログフォロー切替 |
  | `/` | ログ検索 |
  | `y` | ログコピー |
  | `o` | ログエクスポート |
  | `i` | サービス情報 |
  | `v` | 環境変数 |

  ### その他

  | キー | 操作 |
  |------|------|
  | `?` | ヘルプ |
  | `q` | 終了 |

  ## 設定

  設定ファイルは `~/.config/orca/config.yml` に配置します。初回起動時にデフォルト設定ファイルが自動生成されます。

  ```yaml
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
  ```

  <details>
  <summary>キーバインドのカスタマイズ</summary>

  ```yaml
  keybindings:
    up: k          # 上へ移動
    down: j        # 下へ移動
    select: enter  # 選択
    back: esc      # 戻る
    quit: q        # 終了
    tab: tab       # パネル切替
    focus_left: ctrl+h   # 左パネルフォーカス
    focus_right: ctrl+l  # 右パネルフォーカス
    start: u       # サービス起動
    stop: d        # サービス停止
    restart: r     # サービス再起動
    search: /      # 検索
    follow: f      # ログフォロー
    logs: l        # ログ表示
    info: i        # サービス情報
    help: "?"      # ヘルプ
    exec: e        # シェル (exec)
    copy: y        # コピー
    export: o      # エクスポート
    env_vars: v    # 環境変数
    build: b       # イメージビルド
  ```

  </details>
