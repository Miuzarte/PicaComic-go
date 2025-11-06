package PicaComic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-viper/mapstructure/v2"
)

func decodeTo[T any](resp *http.Response, err error) (*http.Response, *T, error) {
	if err != nil {
		return resp, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	input := BaseResp{}
	err = json.Unmarshal(data, &input)
	if err != nil {
		return resp, nil, err
	}

	output := new(T)
	// rv := reflect.ValueOf(output).Elem()
	// if rv.IsValid() && rv.Kind() == reflect.Struct {
	// 	fv := rv.FieldByName("Todo")
	// 	if fv.IsValid() && fv.CanSet() {
	// 		fv.SetBytes(data)
	// 		return resp, output, nil
	// 	}
	// }

	config := mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           output,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		panic(err)
	}
	err = decoder.Decode(input.Data)

	// 避免被 decode err 覆盖
	if input.Code != http.StatusOK {
		return resp, output, fmt.Errorf("api response bad status code: %d, %s", input.Code, input.Message)
	}
	if err != nil {
		return resp, output, err
	}

	return resp, output, nil
}

func toUrl(u string) *url.URL {
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return url
}
