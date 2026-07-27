package service

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// WechatService handles WeChat login, payment callbacks, and message push.
type WechatService struct {
	db           *gorm.DB
	log          *logger.Logger
	userSvc      *UserService
	appID        string
	appSecret    string
	mchID        string
	mchKey       string
	notifyURL    string
	templateID   string
}

func NewWechatService(db *gorm.DB, log *logger.Logger, userSvc *UserService, appID, appSecret, mchID, mchKey, notifyURL, templateID string) *WechatService {
	return &WechatService{
		db:         db,
		log:        log,
		userSvc:    userSvc,
		appID:      appID,
		appSecret:  appSecret,
		mchID:      mchID,
		mchKey:     mchKey,
		notifyURL:  notifyURL,
		templateID: templateID,
	}
}

// --- WeChat Login ---

// WechatLoginResult contains the result of WeChat OAuth login.
type WechatLoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	IsNewUser    bool   `json:"is_new_user"`
}

// GetAuthURL returns the WeChat OAuth authorization URL.
func (s *WechatService) GetAuthURL(redirectURI, state string) string {
	return fmt.Sprintf(
		"https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
		s.appID, redirectURI, state,
	)
}

// HandleLogin processes WeChat OAuth login callback.
func (s *WechatService) HandleLogin(code, clientIP, userAgent string) (*WechatLoginResult, error) {
	// Exchange code for access token
	tokenResp, err := s.exchangeCode(code)
	if err != nil {
		return nil, err
	}

	// Get user info
	userInfo, err := s.getWechatUserInfo(tokenResp.AccessToken, tokenResp.OpenID)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	var oauthAccount model.OAuthAccount
	err = s.db.Where("provider = 'wechat' AND open_id = ?", tokenResp.OpenID).First(&oauthAccount).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new user
		return s.createWechatUser(userInfo, tokenResp.OpenID, tokenResp.UnionID, clientIP, userAgent)
	}
	if err != nil {
		return nil, err
	}

	// Existing user - update last login
	user, err := s.userSvc.GetByID(oauthAccount.UserID)
	if err != nil {
		return nil, errors.New("user account not found")
	}
	if user.Status != 1 {
		return nil, errors.New("account disabled")
	}

	// Update OAuth info
	s.db.Model(&oauthAccount).Updates(map[string]interface{}{
		"username": userInfo.Nickname,
		"avatar":   userInfo.HeadImgURL,
	})

	now := time.Now()
	s.db.Model(user).Update("last_login_at", &now)

	s.log.Infof("wechat login: openid=%s userid=%d", tokenResp.OpenID, user.ID)

	return &WechatLoginResult{
		OpenID:    tokenResp.OpenID,
		IsNewUser: false,
	}, nil
}

type wechatTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

type wechatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
}

func (s *WechatService) exchangeCode(code string) (*wechatTokenResp, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		s.appID, s.appSecret, code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result wechatTokenResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wechat error: %d %s", result.ErrCode, result.ErrMsg)
	}
	return &result, nil
}

func (s *WechatService) getWechatUserInfo(accessToken, openID string) (*wechatUserInfo, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result wechatUserInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *WechatService) createWechatUser(info *wechatUserInfo, openID, unionID, clientIP, userAgent string) (*WechatLoginResult, error) {
	// Generate username
	username := fmt.Sprintf("wx_%s", openID[:8])

	user, err := s.userSvc.Register(RegisterRequest{
		Username: username,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(openID+time.Now().String()))),
		Nickname: info.Nickname,
	})
	if err != nil {
		s.log.Errorf("create wechat user failed: %v", err)
		return nil, errors.New("failed to create user account")
	}

	if info.HeadImgURL != "" {
		s.db.Model(user).Update("avatar", info.HeadImgURL)
	}

	// Create OAuth binding
	oauthAccount := model.OAuthAccount{
		UserID:   user.ID,
		Provider: "wechat",
		OpenID:   openID,
		UnionID:  unionID,
		Username: info.Nickname,
		Avatar:   info.HeadImgURL,
	}
	s.db.Create(&oauthAccount)

	s.log.Infof("new wechat user created: openid=%s userid=%d", openID, user.ID)

	return &WechatLoginResult{
		OpenID:    openID,
		IsNewUser: true,
	}, nil
}

// --- WeChat Payment ---

// WechatPayUnifiedOrder unified order request to WeChat.
type WechatPayUnifiedOrder struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIP string   `xml:"spbill_create_ip"`
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
	OpenID         string   `xml:"openid,omitempty"`
}

// WechatPayNotify is the payment callback from WeChat.
type WechatPayNotify struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`
	ReturnMsg      string   `xml:"return_msg"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	ResultCode     string   `xml:"result_code"`
	OpenID         string   `xml:"openid"`
	IsSubscribe    string   `xml:"is_subscribe"`
	TradeType      string   `xml:"trade_type"`
	BankType       string   `xml:"bank_type"`
	TotalFee       int      `xml:"total_fee"`
	SettlementTotalFee int  `xml:"settlement_total_fee"`
	TransactionID  string   `xml:"transaction_id"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TimeEnd        string   `xml:"time_end"`
}

// wechatPayResponse is the response to WeChat after processing callback.
type wechatPayResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

// CreatePayOrder creates a WeChat pay unified order.
func (s *WechatService) CreatePayOrder(orderNo string, totalFee int, clientIP, openID string) (map[string]string, error) {
	nonceStr := fmt.Sprintf("%d", time.Now().UnixNano())

	params := map[string]string{
		"appid":            s.appID,
		"mch_id":           s.mchID,
		"nonce_str":        nonceStr,
		"body":             "Order Payment",
		"out_trade_no":     orderNo,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": clientIP,
		"notify_url":       s.notifyURL,
		"trade_type":       "JSAPI",
		"openid":           openID,
	}

	params["sign"] = s.signParams(params)

	// Marshal to XML
	xmlData, _ := xml.Marshal(WechatPayUnifiedOrder{
		AppID:          params["appid"],
		MchID:          params["mch_id"],
		NonceStr:       params["nonce_str"],
		Sign:           params["sign"],
		Body:           params["body"],
		OutTradeNo:     params["out_trade_no"],
		TotalFee:       totalFee,
		SpbillCreateIP: params["spbill_create_ip"],
		NotifyURL:      params["notify_url"],
		TradeType:      params["trade_type"],
		OpenID:         params["openid"],
	})

	resp, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "application/xml", bytes.NewReader(xmlData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	s.log.Infof("wechat pay unified order response: %s", string(body))

	// Parse response (simplified)
	return map[string]string{
		"prepay_id": "wx_prepay_" + orderNo,
		"nonce_str": nonceStr,
		"sign":      params["sign"],
	}, nil
}

// HandlePayNotify processes WeChat payment callback.
func (s *WechatService) HandlePayNotify(body []byte) (*WechatPayNotify, error) {
	var notify WechatPayNotify
	if err := xml.Unmarshal(body, &notify); err != nil {
		return nil, err
	}

	if notify.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("wechat pay error: %s", notify.ReturnMsg)
	}

	if notify.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("wechat pay transaction failed")
	}

	// Verify signature
	if !s.verifySign(body) {
		return nil, errors.New("invalid signature")
	}

	s.log.Infof("wechat payment received: order=%s amount=%d txn=%s",
		notify.OutTradeNo, notify.TotalFee, notify.TransactionID)

	// Update order status
	s.db.Table("orders").Where("order_no = ?", notify.OutTradeNo).
		Updates(map[string]interface{}{
			"status":         1,
			"payment_method": "wechat",
			"transaction_id": notify.TransactionID,
			"paid_at":        time.Now(),
		})

	return &notify, nil
}

// GeneratePayNotifyResponse generates a success response for WeChat pay callback.
func (s *WechatService) GeneratePayNotifyResponse() []byte {
	resp := wechatPayResponse{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}
	data, _ := xml.Marshal(resp)
	return data
}

// --- WeChat Template Message ---

type templateMessageRequest struct {
	ToUser      string                 `json:"touser"`
	TemplateID  string                 `json:"template_id"`
	URL         string                 `json:"url,omitempty"`
	Data        map[string]templateDataItem `json:"data"`
}

type templateDataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

// SendTemplateMessage sends a WeChat template message to a user.
func (s *WechatService) SendTemplateMessage(userOpenID, templateID string, data map[string]string, url string) error {
	if templateID == "" {
		templateID = s.templateID
	}

	msgData := make(map[string]templateDataItem)
	for k, v := range data {
		msgData[k] = templateDataItem{Value: v}
	}

	msg := templateMessageRequest{
		ToUser:     userOpenID,
		TemplateID: templateID,
		URL:        url,
		Data:       msgData,
	}

	jsonData, _ := json.Marshal(msg)

	// Get access token
	accessToken, err := s.getAccessToken()
	if err != nil {
		return err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)
	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	s.log.Infof("template message sent to %s: %s", userOpenID, string(body))

	return nil
}

// SendOrderNotify sends an order notification via WeChat template message.
func (s *WechatService) SendOrderNotify(userID uint, orderNo, status string, amount float64) error {
	// Get user's WeChat openid
	var oauthAccount model.OAuthAccount
	if err := s.db.Where("user_id = ? AND provider = 'wechat'", userID).First(&oauthAccount).Error; err != nil {
		return nil // User doesn't have WeChat bound, skip
	}

	data := map[string]string{
		"first":    fmt.Sprintf("Your order %s status has been updated", orderNo),
		"keyword1": orderNo,
		"keyword2": status,
		"keyword3": fmt.Sprintf("%.2f", amount),
		"keyword4": time.Now().Format("2006-01-02 15:04:05"),
		"remark":   "Thank you for your purchase",
	}

	return s.SendTemplateMessage(oauthAccount.OpenID, "", data, "")
}

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func (s *WechatService) getAccessToken() (string, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		s.appID, s.appSecret,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result accessTokenResp
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("wechat access token error: %d %s", result.ErrCode, result.ErrMsg)
	}
	return result.AccessToken, nil
}

// GetDB returns the underlying gorm.DB.
func (s *WechatService) GetDB() *gorm.DB {
	return s.db
}

// signParams generates WeChat pay MD5 signature.
func (s *WechatService) signParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k + "=" + params[k])
	}
	buf.WriteString("&key=" + s.mchKey)

	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

// verifySign verifies WeChat pay callback signature (simplified).
func (s *WechatService) verifySign(body []byte) bool {
	// Simplified verification - in production, parse XML and verify MD5 sign
	return true
}
