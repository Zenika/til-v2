package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (t *testContext) iHaveAJWTTokenForMyUser(ctx context.Context, userKind string) (context.Context, error) {

	code := "my_user_is_granted"
	if userKind == "admin" {
		code = "my_user_is_admin"
	}

	response, err := http.Get(fmt.Sprintf("%s/api/auth?code=%s&state=redirect_uri=http://127.0.0.1:8000", *t.Endpoint, code))
	if err != nil {
		return ctx, err
	}

	if response.StatusCode != 204 {
		return ctx, fmt.Errorf("expected status code 204, got %d", response.StatusCode)
	}

	return context.WithValue(ctx, "token", response.Header.Get("Authorization")), nil
}

func (t *testContext) iClearTheCurrentJWTToken(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, "token", ""), nil
}

func (t *testContext) iSendARequest(ctx context.Context, verb string, endpoint string) (context.Context, error) {
	return t.iSendARequestWithPayload(ctx, verb, endpoint, "")
}

func (t *testContext) iSendARequestWithPayload(ctx context.Context, verb string, endpoint string, payload string) (context.Context, error) {
	endpoint = t.parseStringToInjectContextValues(ctx, endpoint)
	payload = t.parseStringToInjectContextValues(ctx, payload)

	p := strings.NewReader(payload)

	req, err := http.NewRequest(verb, fmt.Sprintf("%s%s", *t.Endpoint, endpoint), p)
	if err != nil {
		return ctx, err
	}
	if ctx.Value("token") != "" {
		req.Header.Set("Authorization", fmt.Sprintf("%s", ctx.Value("token")))
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, "httpResponse", res), nil
}

func (t *testContext) theResponseCodeShouldBe(ctx context.Context, code int) error {
	response := ctx.Value("httpResponse")
	if response == nil {
		return fmt.Errorf("expected response to be non-nil")
	}

	if response.(*http.Response).StatusCode != code {
		return fmt.Errorf("expected status code %d, got %d", code, response.(*http.Response).StatusCode)
	}
	return nil
}

func (t *testContext) iSaveHeaderForSuite(ctx context.Context, header string, name string) (context.Context, error) {
	response := ctx.Value("httpResponse")
	if response == nil {
		return ctx, fmt.Errorf("expected response to be non-nil")
	}

	return context.WithValue(ctx, name, response.(*http.Response).Header.Get(header)), nil
}

func (t *testContext) resetDatabase(ctx context.Context) (context.Context, error) {
	ctx = context.WithValue(ctx, "token", nil)
	ctx = context.WithValue(ctx, "httpResponse", nil)

	if err := pool.Client.RestartContainer(server.Container.ID, 0); err != nil {
		return ctx, err
	}

	c, _ := pool.Client.InspectContainer(server.Container.ID)

	endpoint := fmt.Sprintf("http://172.17.0.1:%s", c.NetworkSettings.Ports["8000/tcp"][0].HostPort)
	t.Endpoint = &endpoint
	time.Sleep(time.Second * 5)

	return ctx, nil
}

func (t *testContext) parseStringToInjectContextValues(ctx context.Context, str string) string {
	r := regexp.MustCompile("({{[^}]+}})")

	for _, val := range r.FindAllString(str, -1) {
		str = strings.ReplaceAll(str, val, ctx.Value(strings.ReplaceAll(strings.ReplaceAll(val, "}}", ""), "{{", "")).(string))
	}
	return str
}

func (t *testContext) theResponseShouldHaveItemsCountInPath(ctx context.Context, nb int, path string) error {
	rawData, err := t.getResponseBody(ctx)
	if err != nil {
		return err
	}

	if rawData, err = t.navigateToPath(rawData, path); err != nil {
		return err
	}

	var data []interface{}
	if err = json.NewDecoder(bytes.NewReader(rawData)).Decode(&data); err != nil {
		return err
	}
	if len(data) != nb {
		return fmt.Errorf("expected %d items in path %s, got %d", nb, path, len(data))
	}

	return nil
}

func (t *testContext) theResponseShouldHaveTheFollowingContentInPath(ctx context.Context, path string, content *godog.Table) error {
	rawData, err := t.getResponseBody(ctx)
	if err != nil {
		return err
	}

	if rawData, err = t.navigateToPath(rawData, path); err != nil {
		return err
	}

	var data map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(rawData)).Decode(&data); err != nil {
		return err
	}
	for key, value := range content.Rows[0].Cells {
		if fmt.Sprintf("%v", data[value.Value]) != t.parseStringToInjectContextValues(ctx, content.Rows[1].Cells[key].Value) {
			return fmt.Errorf("expected %s[%s] to be %s, but got %s", path, value.Value, t.parseStringToInjectContextValues(ctx, content.Rows[1].Cells[key].Value), data[value.Value])
		}
	}

	return nil
}

func (t *testContext) theResponseShouldHaveTheFollowingItemsInPath(ctx context.Context, path string, content *godog.Table) error {
	rawData, err := t.getResponseBody(ctx)
	if err != nil {
		return err
	}

	if rawData, err = t.navigateToPath(rawData, path); err != nil {
		return err
	}

	var data []map[string]interface{}
	if err = json.Unmarshal(rawData, &data); err != nil {
		var data []interface{}
		if err = json.Unmarshal(rawData, &data); err != nil {
			return err
		}

		for key, value := range content.Rows[0].Cells {
			if fmt.Sprintf("%v", data[key]) != t.parseStringToInjectContextValues(ctx, content.Rows[key].Cells[0].Value) {
				return fmt.Errorf("expected %s[%s] to be %s, but got %s", path, value.Value, t.parseStringToInjectContextValues(ctx, content.Rows[key].Cells[0].Value), data[key])
			}
		}
	} else {
		for dataKey, _ := range data {
			for key, value := range content.Rows[0].Cells {
				if fmt.Sprintf("%v", data[dataKey][value.Value]) != t.parseStringToInjectContextValues(ctx, content.Rows[dataKey+1].Cells[key].Value) {
					return fmt.Errorf("expected %s[%s] to be %s, but got %s", path, value.Value, t.parseStringToInjectContextValues(ctx, content.Rows[dataKey+1].Cells[key].Value), data[dataKey][value.Value])
				}
			}
		}
	}

	return nil
}

func (t *testContext) theResponseShouldNotHaveTheKeyInPath(ctx context.Context, key string, path string) error {
	if err := t.theResponseShouldHaveTheKeyInPath(ctx, key, path); err == nil {
		return fmt.Errorf("key %s should not have a value in path %s", key, path)
	}
	return nil
}

func (t *testContext) theResponseShouldHaveTheKeyInPath(ctx context.Context, key string, path string) error {
	rawData, err := t.getResponseBody(ctx)
	if err != nil {
		return err
	}

	if rawData, err = t.navigateToPath(rawData, path); err != nil {
		return err
	}

	var data map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(rawData)).Decode(&data); err != nil {
		return err
	}

	if _, ok := data[key]; ok {
		return nil
	}

	return fmt.Errorf("key %s not found in path %s", key, path)
}

func (t *testContext) iSaveTheValueInPath(ctx context.Context, value string, path string, key string) (context.Context, error) {
	rawData, err := t.getResponseBody(ctx)
	if err != nil {
		return ctx, err
	}

	if rawData, err = t.navigateToPath(rawData, path); err != nil {
		return ctx, err
	}

	var data map[string]interface{}
	if err = json.NewDecoder(bytes.NewReader(rawData)).Decode(&data); err != nil {
		return ctx, err
	}

	if val, ok := data[value]; ok {
		return context.WithValue(ctx, key, val), nil
	}

	return ctx, fmt.Errorf("key %s not found in path %s", key, path)
}

func (t *testContext) getResponseBody(ctx context.Context) ([]byte, error) {
	response := ctx.Value("httpResponse")
	if response == nil {
		return nil, fmt.Errorf("expected response to be non-nil")
	}

	rawData, err := io.ReadAll(response.(*http.Response).Body)
	if err != nil {
		return nil, err
	}

	response.(*http.Response).Body = io.NopCloser(bytes.NewReader(rawData))

	return rawData, nil
}

func (t *testContext) navigateToPath(rawData []byte, path string) ([]byte, error) {
	if path == "@" {
		return rawData, nil
	}

	for _, part := range strings.Split(strings.TrimPrefix(path, "@."), ".") {
		if index, err := strconv.Atoi(part); err == nil {
			var data []interface{}
			if err = json.NewDecoder(bytes.NewReader(rawData)).Decode(&data); err != nil {
				return nil, err
			}
			rawData, _ = json.Marshal(data[index])
		} else {
			var data map[string]interface{}
			if err = json.NewDecoder(bytes.NewReader(rawData)).Decode(&data); err != nil {
				var data []interface{}
				if err = json.NewDecoder(bytes.NewReader(rawData)).Decode(&data); err != nil {
					return nil, err
				}
				rawData, _ = json.Marshal(data)
			} else {
				rawData, _ = json.Marshal(data[part])
			}
		}
	}

	return rawData, nil
}
