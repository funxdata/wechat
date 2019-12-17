package wechat

import (
	"net/http"
	"sync"

	"github.com/funxdata/wechat/cache"
	"github.com/funxdata/wechat/context"
	"github.com/funxdata/wechat/device"
	"github.com/funxdata/wechat/js"
	"github.com/funxdata/wechat/material"
	"github.com/funxdata/wechat/menu"
	"github.com/funxdata/wechat/message"
	"github.com/funxdata/wechat/miniprogram"
	"github.com/funxdata/wechat/oauth"
	"github.com/funxdata/wechat/pay"
	"github.com/funxdata/wechat/qr"
	"github.com/funxdata/wechat/server"
	"github.com/funxdata/wechat/tcb"
	"github.com/funxdata/wechat/user"
)

// Wechat struct
type Wechat struct {
	Context *context.Context
}

// Config for user
type Config struct {
	AppID          string `yaml:"appID"`
	AppSecret      string `yaml:"appSecret"`
	Token          string `yaml:"token"`
	EncodingAESKey string `yaml:"encodingAESKey"`
	PayMchID       string `yaml:"payMchID"`     //支付 - 商户 ID
	PayNotifyURL   string `yaml:"payNotifyURL"` //支付 - 接受微信支付结果通知的接口地址
	PayKey         string `yaml:"payKey"`       //支付 - 商户后台设置的支付 key
	Cache          cache.Cache
}

// NewWechat init
func NewWechat(cfg *Config) *Wechat {
	context := new(context.Context)
	copyConfigToContext(cfg, context)
	return &Wechat{context}
}

func copyConfigToContext(cfg *Config, context *context.Context) {
	context.AppID = cfg.AppID
	context.AppSecret = cfg.AppSecret
	context.Token = cfg.Token
	context.EncodingAESKey = cfg.EncodingAESKey
	context.PayMchID = cfg.PayMchID
	context.PayKey = cfg.PayKey
	context.PayNotifyURL = cfg.PayNotifyURL
	context.Cache = cfg.Cache
	context.SetAccessTokenLock(new(sync.RWMutex))
	context.SetJsAPITicketLock(new(sync.RWMutex))
}

// GetServer 消息管理
func (wc *Wechat) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	wc.Context.Request = req
	wc.Context.Writer = writer
	return server.NewServer(wc.Context)
}

//GetAccessToken 获取access_token
func (wc *Wechat) GetAccessToken() (string, error) {
	return wc.Context.GetAccessToken()
}

// GetOauth oauth2网页授权
func (wc *Wechat) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(wc.Context)
}

// GetMaterial 素材管理
func (wc *Wechat) GetMaterial() *material.Material {
	return material.NewMaterial(wc.Context)
}

// GetJs js-sdk配置
func (wc *Wechat) GetJs() *js.Js {
	return js.NewJs(wc.Context)
}

// GetMenu 菜单管理接口
func (wc *Wechat) GetMenu() *menu.Menu {
	return menu.NewMenu(wc.Context)
}

// GetUser 用户管理接口
func (wc *Wechat) GetUser() *user.User {
	return user.NewUser(wc.Context)
}

// GetTemplate 模板消息接口
func (wc *Wechat) GetTemplate() *message.Template {
	return message.NewTemplate(wc.Context)
}

// GetPay 返回支付消息的实例
func (wc *Wechat) GetPay() *pay.Pay {
	return pay.NewPay(wc.Context)
}

// GetBankPay 返回银行卡提现的实例
func (wc *Wechat) GetBankPay() *pay.Bank {
	return pay.NewWithdrawBank(wc.Context)
}

// GetQR 返回二维码的实例
func (wc *Wechat) GetQR() *qr.QR {
	return qr.NewQR(wc.Context)
}

// GetMiniProgram 获取小程序的实例
func (wc *Wechat) GetMiniProgram() *miniprogram.MiniProgram {
	return miniprogram.NewMiniProgram(wc.Context)
}

// GetDevice 获取智能设备的实例
func (wc *Wechat) GetDevice() *device.Device {
	return device.NewDevice(wc.Context)
}

// GetTcb 获取小程序-云开发的实例
func (wc *Wechat) GetTcb() *tcb.Tcb {
	return tcb.NewTcb(wc.Context)
}
