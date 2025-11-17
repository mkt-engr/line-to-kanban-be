package line

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type WebhookHandler struct {
	client *Client
}

func NewWebhookHandler(client *Client) *WebhookHandler {
	return &WebhookHandler{
		client: client,
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
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// テキストメッセージの内容をログに出力
				log.Printf("【メッセージ検出】ユーザー: %s, 内容: %s", userID, message.Text)

				// ユーザープロフィール取得
				profile, err := h.client.GetBot().GetProfile(userID).Do()
				if err != nil {
					log.Printf("プロフィール取得エラー: %v", err)
					// エラー時はユーザー名なしで返信
					replyMessage := "こんにちは!" + message.Text
					if _, err := h.client.GetBot().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Printf("返信エラー: %v", err)
					}
					return
				}

				// ユーザーに返信
				replyMessage := "こんにちは！" + profile.DisplayName + "さん," + message.Text
				if _, err := h.client.GetBot().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Printf("返信エラー: %v", err)
				}

			default:
				// テキスト以外のメッセージ（スタンプ、画像など）は無視
				log.Printf("未対応のメッセージタイプ: %T", message)
			}
		}
	}

	// LINEへの応答（Webhookハンドラーは常に200 OKを返す必要がある）
	w.WriteHeader(http.StatusOK)
}
