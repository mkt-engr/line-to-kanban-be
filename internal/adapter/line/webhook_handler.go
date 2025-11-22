package line

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"line-to-kanban-be/internal/app/usecase"
)

var deleteCommandPattern = regexp.MustCompile(`^削除\s*(\d+)$`)

type WebhookHandler struct {
	client  *Client
	usecase *usecase.TaskUsecase
}

func NewWebhookHandler(client *Client, taskUsecase *usecase.TaskUsecase) *WebhookHandler {
	return &WebhookHandler{
		client:  client,
		usecase: taskUsecase,
	}
}

func isCommand(text string) bool {
	return text == "一覧" || deleteCommandPattern.MatchString(text)
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

				// コマンドかどうか判定
				if isCommand(lineMessage.Text) {
					// コマンド処理（DBに保存しない）
					h.handleCommand(ctx, lineMessage.Text, userID, event.ReplyToken)
					return
				}

				// 通常のタスクとして保存（usecase経由）
				_, err := h.usecase.CreateTask(ctx, &usecase.CreateTaskRequest{
					UserID:  userID,
					Content: lineMessage.Text,
				})
				if err != nil {
					log.Printf("メッセージ保存エラー: %v", err)
					// 保存に失敗してもユーザーには返信を続ける
				} else {
					log.Printf("メッセージを保存しました")
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
				// テキスト以外のメッセージ（スタンプ、画像など）
				log.Printf("未対応のメッセージタイプ: %T", lineMessage)

				// ユーザーに対応外である旨を通知
				replyMessage := "申し訳ございません。テキストメッセージのみ対応しています。\n\n以下のコマンドが使用できます:\n・一覧: タスク一覧を表示\n・削除 [番号]: タスクを削除\n・その他のテキスト: 新しいタスクとして登録"
				if _, err := h.client.GetBot().ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Printf("返信エラー: %v", err)
				}
			}
		}
	}

	// LINEへの応答（Webhookハンドラーは常に200 OKを返す必要がある）
	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) handleCommand(ctx context.Context, text string, userID string, replyToken string) {
	if text == "一覧" {
		h.handleListCommand(ctx, userID, replyToken)
		return
	}

	// 削除コマンド
	if matches := deleteCommandPattern.FindStringSubmatch(text); matches != nil {
		num, _ := strconv.Atoi(matches[1])
		h.handleDeleteCommand(ctx, num, userID, replyToken)
		return
	}
}

func (h *WebhookHandler) handleListCommand(ctx context.Context, userID string, replyToken string) {
	// ユーザーのタスク一覧を取得（usecase経由）
	tasks, err := h.usecase.ListTasksByUser(ctx, userID)
	if err != nil {
		log.Printf("タスク一覧取得エラー: %v", err)
		h.replyError(replyToken, "タスク一覧の取得に失敗しました")
		return
	}

	// タスクがない場合
	if len(tasks) == 0 {
		h.replyError(replyToken, "まだタスクが登録されていません")
		return
	}

	// 箇条書きでタスク一覧を作成
	var replyText string
	replyText = "【あなたのタスク一覧】\n\n"
	for i, task := range tasks {
		replyText += fmt.Sprintf("%d. %s\n", i+1, task.Content)
	}

	if _, err := h.client.GetBot().ReplyMessage(replyToken, linebot.NewTextMessage(replyText)).Do(); err != nil {
		log.Printf("返信エラー: %v", err)
	}
}

func (h *WebhookHandler) handleDeleteCommand(ctx context.Context, num int, userID string, replyToken string) {
	// バリデーション
	if num <= 0 {
		h.replyError(replyToken, "正しい番号を指定してください")
		return
	}

	// ユーザーのタスク一覧を取得（usecase経由）
	tasks, err := h.usecase.ListTasksByUser(ctx, userID)
	if err != nil {
		log.Printf("タスク一覧取得エラー: %v", err)
		h.replyError(replyToken, "タスクの取得に失敗しました")
		return
	}

	// 範囲チェック
	index := num - 1
	if index < 0 || index >= len(tasks) {
		h.replyError(replyToken, "指定されたタスクが見つかりません")
		return
	}

	// 削除対象のタスク
	targetTask := tasks[index]

	// 削除実行（usecase経由）
	err = h.usecase.DeleteTask(ctx, targetTask.ID, userID)
	if err != nil {
		log.Printf("削除エラー: %v", err)
		h.replyError(replyToken, "削除に失敗しました")
		return
	}

	// 成功レスポンス
	replyMessage := fmt.Sprintf("タスク『%s』を削除しました", targetTask.Content)
	if _, err := h.client.GetBot().ReplyMessage(replyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Printf("返信エラー: %v", err)
	}
	log.Printf("タスク削除成功: ID=%v, Content=%s", targetTask.ID, targetTask.Content)
}

func (h *WebhookHandler) replyError(replyToken string, message string) {
	if _, err := h.client.GetBot().ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		log.Printf("返信エラー: %v", err)
	}
}
