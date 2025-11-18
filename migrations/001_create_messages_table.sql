-- ステータスENUM型の作成
CREATE TYPE message_status AS ENUM ('todo', 'in_progress', 'done');

-- messagesテーブルの作成
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    status message_status NOT NULL DEFAULT 'todo',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- インデックス作成（created_atでの検索を高速化）
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at DESC);

-- インデックス作成（statusでの検索を高速化）
CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status);
