package feishu

import "github.com/xujiahua/alertmanager-webhook-feishu/model"

type FakeBot struct {
	SecretKey string
}

func (f FakeBot) Send(*model.WebhookMessage) error {
	return nil
}

func (f FakeBot) GetSecret() string {
	return f.SecretKey
}
