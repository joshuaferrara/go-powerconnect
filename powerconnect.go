package powerconnect

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"errors"
	"io/ioutil"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

// Holds info needed to authenticate requests.
type AuthInfo struct {
	ip, username, password, ssid string
	client *http.Client
}

// Calculates an MD5 hash and returns it as a string
func getMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

// Logs into the PowerConnect 2724 switch. Creates a AuthInfo objec to use for authentication
// in other functions.
// ip - ip of PowerConnect switch. Ex: 192.168.2.1
// username - username to login with. Default: admin
// password - password to login with. Default: [blank]
func Login(ip, username, password string) (AuthInfo, error) {
	aInfo := AuthInfo{}

	cookieJar, _ := cookiejar.New(nil)

	aInfo.client = &http.Client{
		Jar: cookieJar,
	}

	aInfo.ip = ip
	aInfo.username = username
	aInfo.password = password

	ssidResp, _ := aInfo.client.Get("http://" + aInfo.ip + "/login11.htm");

	ssidData, _ := ioutil.ReadAll(ssidResp.Body)
	ssidBody := string(ssidData)

	ssidIndex := strings.Index(ssidBody, "Session\" value=\"");
	ssid := ssidBody[ssidIndex + 16 : ssidIndex + 16 + 32];

	aInfo.ssid = ssid;

	loginForm := url.Values{"Username": {aInfo.username}, "Password": {getMD5Hash(aInfo.username + aInfo.password + aInfo.ssid)}, "Session": {aInfo.ssid}}

	loginReq, loginReqErr := http.NewRequest("POST", "http://" + aInfo.ip + "/tgi/login.tgi", strings.NewReader(loginForm.Encode()))

	if loginReqErr != nil {
		return aInfo, loginReqErr
	}

	loginReq.Header.Add("Accept", "text/html, application/xhtml+xml, image/jxr, */*")
	loginReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	loginReq.Header.Add("Referer", "http://" + aInfo.ip + "/login11.htm")
	loginReq.Header.Add("Host", aInfo.ip)
	loginResp, loginErr := aInfo.client.Do(loginReq)

	if loginErr != nil {
		return aInfo, loginErr
	}

	loginData, _ := ioutil.ReadAll(loginResp.Body)
	loginBody := string(loginData)

	if strings.Index(loginBody, "Utilization Summary") == -1 {
		return aInfo, errors.New("Login failed. Please verify that the username and password is correct.")
	}

	return aInfo, nil
}