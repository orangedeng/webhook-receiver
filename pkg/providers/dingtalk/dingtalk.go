package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rancher/webhook-receiver/pkg/providers"
)

const (
	Name = "DINGTALK"

	webhookURLKey = "webhook_url"
	secretKey     = "secret"
	proxyURLKey   = "proxy_url"

	//PANDARIA: dingtalk alert message limit
	dingtalkMsgLimit = 19000
)

type sender struct {
	webhookURL string
	secret     string
	proxyURL   string

	client *http.Client
}

func New(opt map[string]string) (providers.Sender, error) {
	if err := validate(opt); err != nil {
		return nil, err
	}

	c := &http.Client{}

	return &sender{
		webhookURL: opt[webhookURLKey],
		secret:     opt[secretKey],
		proxyURL:   opt[proxyURLKey],
		client:     c,
	}, nil
}

func (s *sender) Send(msg string, receiver providers.Receiver) error {
	if len(msg) > dingtalkMsgLimit {
		msg = msg[:dingtalkMsgLimit]
	}

	payload, err := newPayload(msg)
	if err != nil {
		return err
	}

	webhook, err := getWebhook(s.webhookURL, s.secret)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	if s.proxyURL != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(s.proxyURL)
		}

		transport.Proxy = proxy
	}

	s.client.Transport = transport

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dtr := dingtalkResp{}
	if err := json.Unmarshal(respData, &dtr); err != nil {
		return err
	}
	if dtr.ErrCode != 0 {
		return fmt.Errorf("dingtalk response errcode: %d, errmsg: %s", dtr.ErrCode, dtr.ErrMsg)
	}

	return nil
}

type payload struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		IsAtAll bool `json:"isAtAll"`
	} `json:"at"`
}

type dingtalkResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func newPayload(msg string) ([]byte, error) {
	p := payload{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: msg,
		},
		At: struct {
			IsAtAll bool `json:"isAtAll"`
		}{
			IsAtAll: true,
		},
	}

	data, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}
	return data, nil
}

func validate(opt map[string]string) error {
	if _, exists := opt[webhookURLKey]; !exists {
		return fmt.Errorf("%s empty", webhookURLKey)
	}

	return nil
}

func getWebhook(webhook, secret string) (string, error) {
	timestamp := time.Now().UnixNano() / 1e6

	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)

	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	_, err := h.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}

	signData := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sign := url.QueryEscape(signData)

	webhook = fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, timestamp, sign)

	return webhook, nil
}
