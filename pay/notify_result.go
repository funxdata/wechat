package pay

import (
	"encoding/xml"
	"io"
	"sort"
	"strconv"

	"github.com/funxdata/wechat/util"
)

// NotifyResult 下单回调
type NotifyResult map[string]string

// NotifyResp 消息通知返回
type NotifyResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m *NotifyResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = NotifyResult{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func (n NotifyResult) ReturnCode() string    { return n["return_code"] }
func (n NotifyResult) ReturnMsg() string     { return n["return_msg"] }
func (n NotifyResult) ResultCode() string    { return n["result_code"] }
func (n NotifyResult) OpenID() string        { return n["openid"] }
func (n NotifyResult) IsSubscribe() string   { return n["is_subscribe"] }
func (n NotifyResult) TradeType() string     { return n["trade_type"] }
func (n NotifyResult) BankType() string      { return n["bank_type"] }
func (n NotifyResult) TotalFee() int         { return mustParseInt(n["total_fee"]) }
func (n NotifyResult) FeeType() string       { return n["fee_type"] }
func (n NotifyResult) CashFee() int          { return mustParseInt(n["cash_fee"]) }
func (n NotifyResult) CashFeeType() string   { return n["cash_fee_type"] }
func (n NotifyResult) TransactionID() string { return n["transaction_id"] }
func (n NotifyResult) OutTradeNo() string    { return n["out_trade_no"] }
func (n NotifyResult) Attach() string        { return n["attach"] }
func (n NotifyResult) TimeEnd() string       { return n["time_end"] }

// IsSucc 返回是否为成功
func (n NotifyResult) IsSucc() bool { return n.ResultCode() == "SUCCESS" && n.ReturnCode() == "SUCCESS" }

func mustParseInt(val string) int {
	n, _ := strconv.Atoi(val)
	return n
}

// NewNotifyResp 商户处理后同步返回给微信的参数
func NewNotifyResp(isSucc bool, msg ...string) *NotifyResp {
	if isSucc {
		return &NotifyResp{
			ReturnCode: "SUCCESS",
			ReturnMsg:  "OK",
		}
	}
	rmsg := "FAILED"
	if len(msg) > 0 && msg[0] != "" {
		rmsg = msg[0]
	}
	return &NotifyResp{
		ReturnCode: "FAILED",
		ReturnMsg:  rmsg,
	}
}

// VerifySign 验签
func (pcf *Pay) VerifySign(notifyRes NotifyResult) bool {
	// 支付key
	sortedKeys := make([]string, 0, len(notifyRes))
	for k := range notifyRes {
		if k == "sign" {
			continue
		}
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	// STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sortedKeys {
		value := notifyRes[k]
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	// STEP3, 在键值对的最后加上key=API_KEY
	signStrings = signStrings + "key=" + pcf.PayKey
	// STEP4, 进行MD5签名并且将所有字符转为大写.
	sign := util.MD5Sum(signStrings)
	if sign != notifyRes["sign"] {
		return false
	}
	return true
}
