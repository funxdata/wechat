package pay

import (
	"encoding/xml"
	"errors"

	"github.com/funxdata/wechat/context"
	"github.com/funxdata/wechat/util"
)

var payBankURI = "https://api.mch.weixin.qq.com/mmpaysptrans/pay_bank"

// Bank struct extends context
type Bank struct {
	*context.Context
}

// OrderData 提款所需传入参数
type OrderData struct {
	EncBankNo      string
	EncTrueName    string
	BankCode       string
	Amount         int
	Desc           string
	PartnerTradeNo string
}

// BankRequest 接口请求参数
type BankRequest struct {
	// MchID 商户号
	MchID string `xml:"mch_id"`
	// PartnerTradeNo 商户企业付款单号
	PartnerTradeNo string `xml:"partner_trade_no"`
	// NonceStr 随机字符串
	NonceStr string `xml:"nonce_str"`
	// Sign sign
	Sign string `xml:"sign"`
	// EncBankNo 收款方银行卡号
	EncBankNo string `xml:"enc_bank_no"`
	// EncTrueName 收款方用户名
	EncTrueName string `xml:"enc_true_name"`
	// BankCode 收款方开户行
	BankCode string `xml:"bank_code"`
	// Amount 付款金额
	Amount int `xml:"amount"`
	// Desc 付款说明
	Desc string `xml:"desc"`
}

// BankResponse 微信提款接口的返回
type BankResponse struct {
	// ReturnCode 通信标识
	ReturnCode string `xml:"return_code"`
	// ReturnMsg 错误信息
	ReturnMsg string `xml:"return_msg"`
	// ResultCode 业务结果
	ResultCode string `xml:"result_code,omitempty"`
	// ErrCode 错误代码
	ErrCode string `xml:"err_code,omitempty"`
	// ErrCodeDes 错误代码描述
	ErrCodeDes string `xml:"err_code_des,omitempty"`
	// MchID 商户ID
	MchID string `xml:"mch_id,omitempty"`
	// PartnerTradeNo 企业商户付款单号
	PartnerTradeNo string `xml:"partner_trade_no,omitempty"`
	// Amount 付款金额
	Amount int `xml:"amount"`
	// NonceStr 随机字符串
	NonceStr string `xml:"nonce_str,omitempty"`
	// Sign 签名
	Sign string `xml:"sign,omitempty"`
	// PaymentNo 微信企业付款单号
	PaymentNo string `xml:"payment_no"`
	// CmmsAmt 手续费金额
	CmmsAmt int `xml:"cmms_amt"`
}

// NewWithdrawBank return an instance of Bank package
func NewWithdrawBank(ctx *context.Context) *Bank {
	pay := &Bank{Context: ctx}
	return pay
}

// PreBankOrder return data for invoke wechat payment
func (pcf *Bank) PreBankOrder(p *OrderData) (payOrder BankResponse, err error) {
	nonceStr := util.RandomStr(32)
	param := make(map[string]interface{})
	param["appid"] = pcf.AppID
	param["mch_id"] = pcf.PayMchID
	param["nonce_str"] = nonceStr
	param["notify_url"] = pcf.PayNotifyURL

	bizKey := "&key=" + pcf.PayKey
	str := orderParam(param, bizKey)
	sign := util.MD5Sum(str)
	request := BankRequest{
		MchID:          pcf.PayMchID,
		PartnerTradeNo: p.PartnerTradeNo,
		NonceStr:       nonceStr,
		Sign:           sign,
		EncBankNo:      p.EncBankNo,
		EncTrueName:    p.EncTrueName,
		BankCode:       p.BankCode,
		Amount:         p.Amount,
		Desc:           p.Desc,
	}
	rawRet, err := util.PostXML(payGateway, request)
	if err != nil {
		return
	}
	err = xml.Unmarshal(rawRet, &payOrder)
	if err != nil {
		return
	}
	if payOrder.ReturnCode == "SUCCESS" {
		//pay success
		if payOrder.ResultCode == "SUCCESS" {
			err = nil
			return
		}
		err = errors.New(payOrder.ErrCode + payOrder.ErrCodeDes)
		return
	}
	err = errors.New("[msg : xmlUnmarshalError] [rawReturn : " + string(rawRet) + "] [params : " + str + "] [sign : " + sign + "]")
	return
}
