package mysql

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/Mictrlan/Miuer/smsservice/model/mysql"
)

// SendSmsReply -
type SendSmsReply struct {
	Message   string `json:"Message,omitempty"`
	RequestID string `json:"RequestID,omitempty"`
	BizID     string `json:"BizID,omitempty"`
	Code      string `json:"Code,omitempty"`
}

// SMS -
type SMS struct {
	Mobile string
	Date   int64
	Code   string
	Sign   string
}

// newSMS return a new *SMS
func newSMS() *SMS {
	sms := &SMS{}
	return sms
}

// Config -
type Config struct {
	Host           string
	Appcode        string
	Digits         int
	ResendInterval int
	OnCheck        SMSVerify
	DB             *sql.DB
}

// SMSVerify -
type SMSVerify interface {
	OnVerifySucceed(targetID, mobile string)
	OnVerifyFailed(targetID, mobile string)
}

// Send -
func Send(db *sql.DB, mobile, sign string, conf *Config) error {
	sms := newSMS()
	sms.prepare(mobile, sign, conf.Digits)

	if err := sms.checkvalid(db, conf); err != nil {
		return err
	}

	if err := sms.save(db); err != nil {
		return err
	}

	// err !!!
	if err := sms.sendmsg(conf); err != nil {
		return err
	}

	return nil
}

// prepare prepare data to send
func (sms *SMS) prepare(mobile, sign string, digits int) {
	sms.Mobile = mobile
	sms.Date = time.Now().Unix()
	sms.Code = Code(digits)
	sms.Sign = sign
}

// Code generate size bits of captcha
func Code(size int) string {
	data := make([]byte, size)
	out := make([]byte, size)
	buffer := len(numbers)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}

	for ix, key := range data {
		temp := byte(int(key) % buffer) // temp < buffer
		out[ix] = numbers[temp]
	}
	return string(out)
}

var numbers = []byte("012345678998765431234567890987654321")

// checkvalid get date(unixtime) and verify the validity of the phone number
func (sms *SMS) checkvalid(db *sql.DB, conf *Config) error {
	unixtime, _ := mysql.GetDateBySign(db, sms.Sign)

	if unixtime > 0 && sms.Date-unixtime < int64(conf.ResendInterval) {
		return errors.New("短时间内不允许发送两次")
	}
	if err := VailMobile(sms.Mobile); err != nil {
		return errors.New("手机号不符合规则")
	}

	return nil
}

// VailMobile  verify the validity of the phone number
func VailMobile(mobile string) error {

	if len(mobile) != 11 {
		return errors.New("[mobile]参数不对")
	}
	reg, err := regexp.Compile("^1[3-8][0-9]{9}$")
	if err != nil {
		panic("regexp error")
	}
	if !reg.MatchString(mobile) {
		return errors.New("手机号码[mobile]格式不正确")
	}
	return nil
}

// save store data in database
func (sms *SMS) save(db *sql.DB) error {
	err := mysql.AddSmsMessage(db, sms.Mobile, sms.Date, sms.Code, sms.Sign)
	return err
}

// Sendmsg -
func (sms *SMS) sendmsg(conf *Config) error {
	host := conf.Host

	url := host + "?" + "code=" + sms.Code + "&phone=" + sms.Mobile + "&skin=1"

	client := &http.Client{}

	// Configure the HTTP request message
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "APPCODE "+conf.Appcode)

	// get http response
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// parse response
	ssr := &SendSmsReply{}
	if err := json.Unmarshal(body, ssr); err != nil {
		return err
	}

	if ssr.Code != "OK" {
		return err
	}

	return nil
}

// Check delete messgae
func Check(code, sign string, conf *Config, db *sql.DB) error {
	sms := newSMS()
	sms.Date = time.Now().Unix()
	sms.Code = code
	sms.Sign = sign

	smsCode, _ := mysql.GetCodeBySign(db, sign)

	if sms.Code == smsCode {
		sms.delete(db, sign)
		return nil
	}

	return errors.New("Unknown error")
}

func (sms *SMS) delete(db *sql.DB, sign string) {
	mysql.DeleteSmsMessage(db, sign)
}
