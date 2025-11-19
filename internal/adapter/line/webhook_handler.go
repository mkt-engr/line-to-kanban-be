package line

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"line-to-kanban-be/internal/adapter/repository/db"
)

type WebhookHandler struct {
	client  *Client
	queries db.Querier
}

func NewWebhookHandler(client *Client, queries db.Querier) *WebhookHandler {
	return &WebhookHandler{
		client:  client,
		queries: queries,
	}
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, req *http.Request) {
	// 1. Webhookリクエストのパースと署名検証
	events, err := h.client.GetBot().ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Printf("エラー: 不正な署名 (Status 400)")
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Printf("エラー: Webhookパース中にサーバーエラー (Status 500): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// 2. 受信した各イベントの処理
	for _, event := range events {
		userID := event.Source.UserID
		log.Printf("受信イベントの処理を開始: ユーザーID: %s, イベントタイプ: %s", userID, event.Type)

		// メッセージイベントのみを対象とする
		if event.Type == linebot.EventTypeMessage {
			switch lineMessage := event.Message.(type) {
			case *linebot.TextMessage:
				// テキストメッセージの内容をログに出力
				log.Printf("【メッセージ検出】ユーザー: %s, 内容: %s", userID, lineMessage.Text)

				ctx := context.Background()

				// 「一覧」コマンドの処理
				if lineMessage.Text == "一覧" {
					// ユーザーの過去のメッセージを取得
					messages, err := h.queries.ListMessagesByUser(ctx, userID)
					if err != nil {
						log.Printf("メッセージ一覧取得エラー: %v", err)
						replyMessage := "メッセージ一覧の取得に失敗しました"
						if _, err := h.client.GetBot().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							log.Printf("返信エラー: %v", err)
						}
						return
					}

					// メッセージがない場合
					if len(messages) == 0 {
						replyMessage := "まだメッセージが登録されていません"
						if _, err := h.client.GetBot().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							log.Printf("返信エラー: %v", err)
						}
						return
					}

					// 箇条書きでメッセージ一覧を作成
					var replyText string
					replyText = "【あなたのメッセージ一覧】\n\n"
					for i, msg := range messages {
						// 「一覧」コマンド自体は除外
						if msg.Content == "一覧" {
							continue
						}
						replyText += fmt.Sprintf("%d. %s\n", i+1, msg.Content)
					}

					if _, err := h.client.GetBot().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do(); err != nil {
						log.Printf("返信エラー: %v", err)
					}
					return
				}

				// 通常のメッセージ処理（「一覧」以外）
				// メッセージをデータベースに保存 (sqlc使用)
				savedMsg, err := h.queries.CreateMessage(ctx, db.CreateMessageParams{
					Content: lineMessage.Text,
					Status:  db.MessageStatusTodo,
					UserID:  userID,
				})
				if err != nil {
					log.Printf("メッセージ保存エラー: %v", err)
					// 保存に失敗してもユーザーには返信を続ける
				} else {
					log.Printf("メッセージを保存しました: ID=%v", savedMsg.ID)
				}

				// ユーザープロフィール取得
				profile, err := h.client.GetBot().GetProfile(userID).Do()
				if err != nil {
					log.Printf("プロフィール取得エラー: %v", err)
					// エラー時はユーザー名なしで返信
					replyMessage := "こんにちは!" + lineMessage.Text
					if _, err := h.client.GetBot().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Printf("返信エラー: %v", err)
					}
					return
				}

				// ユーザーに返信
				replyMessage := "こんにちは！" + profile.DisplayName + "さん," + lineMessage.Text
				if _, err := h.client.GetBot().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Printf("返信エラー: %v", err)
				}

			default:
				// テキスト以外のメッセージ（スタンプ、画像など）は無視
				log.Printf("未対応のメッセージタイプ: %T", lineMessage)
			}
		}
	}

	// LINEへの応答（Webhookハンドラーは常に200 OKを返す必要がある）
	w.WriteHeader(http.StatusOK)
}
