// add by stefan

package logindefine

const (
	//wechat
	weChatAPPID            = "wxd5accb31d9319f12"
	weChatAPPSecret        = "382af0b77ad4e311b0f1e8ff9b56a3b8"
	urlHttpsServerLink     = "http://aoko.test.com"
	urlWeChatAuthorization = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=1#wechat_redirect"
)

// LoadAccessTokenReply 微信拉取access token回复
type AccessTokenReply struct {
	szSessionKey string `json:"session_key"`
	szOpenID     string `json:"openid"`

	nErrorCode int32  `json:"errcode"`
	szErrorMsg string `json:"errmsg"`
}

// UserInfoReply 用户信息数据
type UserInfoReply struct {
	szNickName  string `json:"nickName"`
	szAvatarURL string `json:"avatarURL"`
}
