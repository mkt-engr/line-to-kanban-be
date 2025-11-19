-- messagesテーブルにuser_idカラムを追加
ALTER TABLE messages ADD COLUMN user_id TEXT NOT NULL DEFAULT '';

-- user_idにインデックスを追加（ユーザー別の検索を高速化）
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
