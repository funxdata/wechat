package qr

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/funxdata/wechat/context"
	"github.com/funxdata/wechat/util"
)

const (
	tmpQRCreateURL = "https://api.weixin.qq.com/cgi-bin/qrcode/create"
	qrImgURL       = "https://mp.weixin.qq.com/cgi-bin/showqrcode"
)

//QR 模板消息
type QR struct {
	*context.Context
}

//NewQR 实例化
func NewQR(context *context.Context) *QR {
	q := new(QR)
	q.Context = context
	return q
}

const (
	actionId  = "QR_SCENE"
	actionStr = "QR_STR_SCENE"
)

type TmpQR struct {
	ExpireSeconds int    `json:"expire_seconds"`
	ActionName    string `json:"action_name"`
	ActionInfo    struct {
		Scene struct {
			SceneStr string `json:"scene_str,omitempty"`
			SceneId  int    `json:"scene_id,omitempty"`
		} `json:"scene"`
	} `json:"action_info"`
}

type Ticket struct {
	util.CommonError `json:",inline"`
	Ticket           string `json:"ticket"`
	ExpireSeconds    int    `json:"expire_seconds"`
	url              string `json:"url"`
}

func newStrTmpQrRequest(exp time.Duration, str string) *TmpQR {
	tq := &TmpQR{
		ExpireSeconds: int(exp.Seconds()),
		ActionName:    actionStr,
	}
	tq.ActionInfo.Scene.SceneStr = str

	return tq
}

//Ticket 获取QR Ticket
func (q *QR) GetStrQRTicket(exp time.Duration, str string) (t *Ticket, err error) {
	tq := newStrTmpQrRequest(exp, str)

	var accessToken string
	accessToken, err = q.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", tmpQRCreateURL, accessToken)
	response, err := util.PostJSON(uri, tq)
	if err != nil {
		err = fmt.Errorf("get qr ticket failed, %s", err)
		return
	}

	t = new(Ticket)
	err = json.Unmarshal(response, &t)
	if err != nil {
		return
	}
	if t.CommonError.ErrCode > 0 {
		err = fmt.Errorf("[%v] %s", t.CommonError.ErrCode, t.CommonError.ErrMsg)
	}
	return
}

func ShowQRCode(tk *Ticket) string {
	u, _ := url.Parse(qrImgURL)
	q := u.Query()
	q.Set("ticket", tk.Ticket)
	u.RawQuery = q.Encode()
	return u.String()
}
