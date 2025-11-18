#!/bin/bash

# マイグレーションスクリプト
# 使い方: ./scripts/migrate.sh

set -e

# .envファイルから環境変数を読み込む
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# DATABASE_URLが設定されているか確認
if [ -z "$DATABASE_URL" ]; then
    echo "エラー: DATABASE_URLが設定されていません"
    exit 1
fi

echo "マイグレーションを実行します..."
echo "接続先: ${DATABASE_URL%%\?*}"  # パスワードなどを隠すため?以降を削除して表示

# PostgreSQL/CockroachDB用のクライアントでマイグレーションを実行
# psqlがインストールされている場合
if command -v psql &> /dev/null; then
    for file in migrations/*.sql; do
        echo "実行中: $file"
        psql "$DATABASE_URL" -f "$file"
    done
    echo "マイグレーション完了"
else
    echo "エラー: psqlコマンドが見つかりません"
    echo "代替手段: CockroachDB Cloud Console から以下のSQLを実行してください:"
    echo ""
    cat migrations/*.sql
    exit 1
fi
