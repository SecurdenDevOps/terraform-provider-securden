package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func api_call(params map[string]string) ([]byte, error) {
	client := &http.Client{}
	url := SecurdenServerURL + "/api/get_account_details_dict"

	api_request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	api_request.Header.Set("authtoken", SecurdenAuthToken)

	q := api_request.URL.Query()
	for key, value := range params {
		if value != "" {
			q.Add(key, value)
		}
	}
	api_request.URL.RawQuery = q.Encode()

	resp, err := client.Do(api_request)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func get_account(ctx context.Context, account_id, account_name, account_title, key_field string) (SecurdenDataSourceModel, int, string) {
	var account SecurdenDataSourceModel
	params := map[string]string{
		"account_id":    account_id,
		"account_name":  account_name,
		"account_title": account_title,
		"key_field":     key_field,
	}
	var Response struct {
		AccountID         int64  `json:"account_id"`
		AccountName       string `json:"account_name"`
		AccountTitle      string `json:"account_title"`
		Password          string `json:"password"`
		KeyValue          string `json:"key_value"`
		PrivateKey        string `json:"private_key"`
		PuTTYPrivateKey   string `json:"putty_private_key"`
		Passphrase        string `json:"passphrase"`
		PPKPassphrase     string `json:"ppk_passphrase"`
		Address           string `json:"address"`
		ClientID          string `json:"client_id"`
		ClientSecret      string `json:"client_secret"`
		AccountAlias      string `json:"account_alias"`
		AccountFile       string `json:"account_file"`
		Port              string `json:"port"`
		OracleSID         string `json:"oracle_sid"`
		OracleServiceName string `json:"oracle_service_name"`
		DefaultDatabase   string `json:"default_database"`
		StatusCode        int    `json:"status_code"`
		Message           string `json:"message"`
		Error             struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	body, err := api_call(params)
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}
	json.Unmarshal(body, &Response)
	if Response.StatusCode != 200 {
		if Response.Error.Message != "" {
			return account, Response.StatusCode, Response.Error.Message
		}
		return account, Response.StatusCode, Response.Message
	}
	account.AccountID = types.Int64Value(Response.AccountID)
	account.AccountName = types.StringValue(Response.AccountName)
	account.AccountTitle = types.StringValue(Response.AccountTitle)
	account.Password = types.StringValue(Response.Password)
	account.KeyValue = types.StringValue(Response.KeyValue)
	account.PrivateKey = types.StringValue(Response.PrivateKey)
	account.PuTTYPrivateKey = types.StringValue(Response.PuTTYPrivateKey)
	account.Passphrase = types.StringValue(Response.Passphrase)
	account.PPKPassphrase = types.StringValue(Response.PPKPassphrase)
	account.Address = types.StringValue(Response.Address)
	account.ClientID = types.StringValue(Response.ClientID)
	account.ClientSecret = types.StringValue(Response.ClientSecret)
	account.AccountAlias = types.StringValue(Response.AccountAlias)
	account.AccountFile = types.StringValue(Response.AccountFile)
	account.Port = types.StringValue(Response.Port)
	account.OracleSID = types.StringValue(Response.OracleSID)
	account.OracleServiceName = types.StringValue(Response.OracleServiceName)
	account.DefaultDatabase = types.StringValue(Response.DefaultDatabase)
	return account, Response.StatusCode, Response.Message
}
