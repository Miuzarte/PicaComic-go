package PicaComic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-viper/mapstructure/v2"
)

func doApiAndDecodeTo[T any](ctx context.Context, method string, url string, body any) (*T, error) {
	_, b, err := DoApi(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))
	t, err := decodeTo[T](b)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func decodeTo[T any](body []byte) (*T, error) {
	input := BaseResp{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		return nil, err
	}

	output := new(T)
	// rv := reflect.ValueOf(output).Elem()
	// if rv.IsValid() && rv.Kind() == reflect.Struct {
	// 	fv := rv.FieldByName("Todo")
	// 	if fv.IsValid() && fv.CanSet() {
	// 		fv.SetBytes(body)
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
		return output, fmt.Errorf("api response bad status code: %d, %s", input.Code, input.Message)
	}
	if err != nil {
		return output, err
	}

	return output, nil
}

func toUrl(u string) *url.URL {
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return url
}
