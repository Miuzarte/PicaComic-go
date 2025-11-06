package PicaComic

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Miuzarte/PicaComic-go/internal/constant"
)

var authorization = ""

func SetToken(token string) {
	authorization = token
}

type BaseResp struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"` // delay decode using [mapstructure.Decode]
}

func DoApi(ctx context.Context, method string, url string, body any) (*http.Response, error) {
	const (
		nonce     = "b1ab87b4800d4d4590a11701b8551afa"
		apiKey    = "C69BAF41DA5ABD1FFEDC6D2FEA56B"
		secretKey = "~d}$Q7$eIni=V)9\\RK/P.RM4;9[7|@/CA}b~OW!3?EV`:<>M7pddUBL5n|0/*Cn"
	)

	pathQuery := strings.ReplaceAll(url, API_URL+"/", "")
	time := strconv.FormatInt(time.Now().Unix(), 10)

	msg := strings.ToLower(pathQuery + time + nonce + method + apiKey)
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	signature := hex.EncodeToString(mac.Sum(nil))

	var bodyReader io.Reader
	if body != nil {
		switch v := body.(type) {
		case io.Reader:
			bodyReader = v
		default:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(b)
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	for k, v := range constant.Header {
		req.Header.Set(k, v)
	}
	req.Header.Set("signature", signature)
	req.Header.Set("time", time)
	if authorization != "" {
		req.Header.Set("authorization", authorization)
	}

	return http.DefaultClient.Do(req)
}

type Image struct {
	FileServer   string `json:"fileServer"`
	OriginalName string `json:"originalName"`
	Path         string `json:"path"`
}

func (i *Image) Download(ctx context.Context) (resp *http.Response, body []byte, err error) {
	u, err := url.Parse(i.FileServer)
	if err != nil {
		return nil, nil, err
	}
	if strings.Contains(u.Path, "static") {
		u.Path = path.Join(u.Path, i.Path)
	} else {
		u.Path = path.Join(u.Path, IMG_STATIC, i.Path)
	}

	// 图片下载不需要鉴权
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp, body, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, body, fmt.Errorf("Image.Download: bad http status code: %s", resp.Status)
	}
	return resp, body, nil
}
