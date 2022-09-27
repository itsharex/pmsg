package message

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/lenye/pmsg/pkg/flags"
	"github.com/lenye/pmsg/pkg/http/client"
	"github.com/lenye/pmsg/pkg/weixin"
	"github.com/lenye/pmsg/pkg/weixin/token"
)

type CmdMpSendTemplateSubscribeParams struct {
	UserAgent   string
	AccessToken string
	AppID       string
	AppSecret   string
	ToUser      string
	TemplateID  string
	Url         string
	Mini        map[string]string
	Scene       string
	Title       string
	Data        string
}

func (t *CmdMpSendTemplateSubscribeParams) Validate() error {
	if t.AccessToken == "" && t.AppID == "" {
		return flags.ErrWeixinAccessToken
	}

	if len(t.Title) > 15 {
		return fmt.Errorf("flag %q maximum length is within 15", flags.Title)
	}

	// 跳小程序
	if t.Mini != nil {
		if miniAppID, ok := t.Mini[flags.MiniAppID]; !ok {
			return fmt.Errorf("mini flag %q not set", flags.MiniAppID)
		} else {
			if miniAppID == "" {
				return fmt.Errorf("mini flag %q not set", flags.MiniAppID)
			}
		}

		if miniPagePath, ok := t.Mini[flags.MiniPagePath]; !ok {
			return fmt.Errorf("mini flag %q not set", flags.MiniPagePath)
		} else {
			if miniPagePath == "" {
				return fmt.Errorf("mini flag %q not set", flags.MiniPagePath)
			}
		}
	}

	return nil
}

// CmdMpSendTemplateSubscribe 发送微信公众号一次性订阅消息
func CmdMpSendTemplateSubscribe(arg *CmdMpSendTemplateSubscribeParams) error {

	if err := arg.Validate(); err != nil {
		return err
	}

	var dataItem map[string]TemplateDataItem
	buf := bytes.NewBufferString("")
	buf.WriteString(arg.Data)
	if buf.String() != "" {
		if err := json.Unmarshal(buf.Bytes(), &dataItem); err != nil {
			return fmt.Errorf("invalid json format, %v", err)
		}
		for k, v := range dataItem {
			if v.Value == "" {
				return fmt.Errorf("data %v.value not set", k)
			}
			if len(v.Value) > 200 {
				return fmt.Errorf("data %v.value maximum length is within 200", k)
			}
		}
	}

	msg := TemplateSubscribeMessage{
		ToUser:     arg.ToUser,
		TemplateID: arg.TemplateID,
		Data:       dataItem,
		URL:        arg.Url,
		Scene:      arg.Scene,
		Title:      arg.Title,
	}

	// 跳小程序
	if arg.Mini != nil {
		miniAppID, _ := arg.Mini[flags.MiniAppID]
		miniPagePath, _ := arg.Mini[flags.MiniPagePath]
		msg.MiniProgram = &MiniProgramMeta{
			AppID:    miniAppID,
			PagePath: miniPagePath,
		}
	}

	client.UserAgent = arg.UserAgent

	if arg.AccessToken == "" {
		accessTokenResp, err := token.GetAccessToken(arg.AppID, arg.AppSecret)
		if err != nil {
			return err
		}
		arg.AccessToken = accessTokenResp.AccessToken
	}

	if err := SendTemplateSubscribe(arg.AccessToken, &msg); err != nil {
		return err
	}
	fmt.Println(weixin.MessageOK)

	return nil
}
