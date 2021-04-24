package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

const (
	ACT_ID = "e202102251931481"
	// OS_ROLE_URL   = "https://api-os-takumi.mihoyo.com/auth/api/getUserAccountInfoByLToken?t=%v&ltoken=%v&uid=%v"
	DEFAULT_USER_AGENT = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E150"
)

type Account struct {
	cookie    string
	Lang      string
	UserAgent string
	client    http.Client
	// actId     string
}

func (acc *Account) getRefererUrl() string {
	const OS_REFERER_URL = "https://webstatic-sea.mihoyo.com/ys/event/signin-sea/index.html?act_id=%v"
	return fmt.Sprintf(OS_REFERER_URL, ACT_ID)
}
func (acc *Account) getRewardUrl() string {
	const OS_REWARD_URL = "https://hk4e-api-os.mihoyo.com/event/sol/home?lang=%v&act_id=%v"
	return fmt.Sprintf(OS_REWARD_URL, acc.Lang, ACT_ID)
}
func (acc *Account) getInfoUrl() string {
	const OS_INFO_URL = "https://hk4e-api-os.mihoyo.com/event/sol/info?lang=%v&act_id=%v"
	return fmt.Sprintf(OS_INFO_URL, acc.Lang, ACT_ID)
}
func (acc *Account) getSignUrl() string {
	const OS_SIGN_URL = "https://hk4e-api-os.mihoyo.com/event/sol/sign?lang=%v"
	return fmt.Sprintf(OS_SIGN_URL, acc.Lang)
}

var langRgx = regexp.MustCompile("mi18nLang=([a-zA-Z]{2}-[a-zA-Z]{2})")

// var cookieRgx = regexp.MustCompile(`^account_id=[0-9]{9}; cookie_token=\w{40}; _MHYUUID=\w{8}-\w{4}-\w{4}-\w{4}-\w{12}; ltoken=\w{40}; ltuid=[0-9]{9}; login_ticket=\w{40}; mi18nLang=[a-zA-Z]{2}-[a-zA-Z]{2}$`)

func New(cookie string) (*Account, error) {
	lang := langRgx.FindStringSubmatch(cookie)[1]
	// if !cookieRgx.MatchString(cookie) {
	// 	return nil, fmt.Errorf("cookie did not pass the check")
	// }
	return &Account{cookie: cookie, Lang: lang, UserAgent: DEFAULT_USER_AGENT}, nil
}
func (acc *Account) newRequest(method string, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Add("User-Agent", acc.UserAgent)
	req.Header.Add("Referer", acc.getRefererUrl())
	// req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Cookie", acc.cookie)
	return req
}
func (acc *Account) SignIn() error {
	req := acc.newRequest("POST", acc.getSignUrl())
	jsonBody, err := json.Marshal(struct {
		ActId string `json:"act_id"`
	}{ACT_ID})
	if err != nil {
		log.Panicf("json formatting error: %v", err)
	}

	req.Body = io.NopCloser(bytes.NewReader(jsonBody))
	body, err := acc.doRequest(req)
	if err != nil {
		return fmt.Errorf("HTTP IO POST request error: %w", err)
	}
	var jsonResp SignInError
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return fmt.Errorf("error parsing parsing body as json: %w", err)
	}
	if jsonResp.Retcode != 0 || jsonResp.Data.Code != "ok" || jsonResp.Message != "OK" {
		return &jsonResp
	}
	return nil
}
func (acc *Account) GetInfo() (InfoResponse, error) {
	var ir InfoResponse
	req := acc.newRequest("GET", acc.getInfoUrl())
	body, err := acc.doRequest(req)
	if err != nil {
		return ir, fmt.Errorf("request error: %w", err)
	}
	err = json.Unmarshal(body, &ir)
	if err != nil {
		return ir, fmt.Errorf("error parsing parsing body as json: %w", err)
	}
	if ir.Retcode != 0 && ir.Message != "OK" {
		return ir, fmt.Errorf("mihoyo error: %v", ir.Message)
	}
	return ir, nil
}
func (acc *Account) GetAwards() (AwardsResponse, error) {
	var ar AwardsResponse
	req := acc.newRequest("GET", acc.getRewardUrl())
	body, err := acc.doRequest(req)
	if err != nil {
		return ar, fmt.Errorf("request error: %w", err)
	}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return ar, fmt.Errorf("error parsing parsing body as json: %w", err)
	}
	if ar.Retcode != 0 && ar.Message != "OK" {
		return ar, fmt.Errorf("mihoyo error: %v", ar.Message)
	}
	return ar, nil
}

func (acc *Account) doRequest(req *http.Request) ([]byte, error) {
	resp, err := acc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP IO request error: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP IO body read error: %w", err)
	}
	return body, nil
}

type SignInError struct {
	Data *struct {
		Code string `json:"code"`
	} `json:"data"`
	Message string `json:"message"`
	Retcode int    `json:"retcode"`
}

func (e *SignInError) Error() string {
	return fmt.Sprintf("json: message = %v, retcode = %v", e.Message, e.Retcode)
}

type InfoResponse struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		TotalSignDay int    `json:"total_sign_day"`
		Today        string `json:"today"`
		IsSign       bool   `json:"is_sign"`
		FistBind     bool   `json:"first_bind"`
	} `json:"data"`
}
type AwardsResponse struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Month  int     `json:"month"`
		Awards []Award `json:"awards"`
	} `json:"data"`
}
type Award struct {
	IconUrl string `json:"icon"`
	Name    string `json:"name"`
	Count   int    `json:"cnt"`
}
