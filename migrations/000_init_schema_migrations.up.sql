-- マイグレーション管理用テーブル
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INT PRIMARY KEY,
    dirty BOOLEAN NOT NULL DEFAULT FALSE
);
