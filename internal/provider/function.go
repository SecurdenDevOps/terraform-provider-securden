package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func logger(data interface{}) error {

	strData := fmt.Sprintf("%v", data)

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(strData + "\n"); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func get_request(params map[string]string, api_url string) ([]byte, error) {
	client := &http.Client{}

	api_request, err := http.NewRequest("GET", api_url, nil)
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
		OracleSID         string `json:"oracle_sid"`
		OracleServiceName string `json:"oracle_service_name"`
		DefaultDatabase   string `json:"default_database"`
		Port              string `json:"port"`
		StatusCode        int    `json:"status_code"`
		Message           string `json:"message"`
		Error             struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	api_url := SecurdenServerURL + "/api/get_account_details_dict"
	body, err := get_request(params, api_url)

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
	account.OracleSID = types.StringValue(Response.OracleSID)
	account.OracleServiceName = types.StringValue(Response.OracleServiceName)
	account.DefaultDatabase = types.StringValue(Response.DefaultDatabase)
	account.Port = types.StringValue(Response.Port)
	return account, Response.StatusCode, Response.Message
}

func post_request(params map[string]interface{}, apiURL string) ([]byte, error) {
	client := &http.Client{}

	requestBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request body: %v", err)
	}

	apiRequest, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("authtoken", SecurdenAuthToken)

	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func get_passwords(ctx context.Context, accountIDs []string) (types.Map, int, string) {
	var accountIDsInt64 []int64
	for _, id := range accountIDs {
		accountID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return types.Map{}, 400, fmt.Sprintf("Invalid account ID format: %v", err)
		}
		accountIDsInt64 = append(accountIDsInt64, accountID)
	}

	params := map[string]interface{}{
		"account_ids": accountIDsInt64,
	}

	apiURL := SecurdenServerURL + "/api/get_multiple_accounts_passwords"

	body, err := post_request(params, apiURL)
	if err != nil {
		return types.Map{}, 500, fmt.Sprintf("Error in API call: %v", err)
	}

	var response struct {
		Passwords  map[string]string `json:"passwords"`
		StatusCode int               `json:"status_code"`
		Message    string            `json:"message"`
		Error      struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return types.Map{}, 500, fmt.Sprintf("Failed to parse response: %v", err)
	}

	if response.StatusCode != 200 {
		errorMessage := response.Message
		if response.Error.Message != "" {
			errorMessage = response.Error.Message
		}
		return types.Map{}, response.StatusCode, errorMessage
	}

	passwordsMap := make(map[string]attr.Value, len(response.Passwords))
	for k, v := range response.Passwords {
		passwordsMap[k] = types.StringValue(v)
	}

	passwords, diags := types.MapValue(types.StringType, passwordsMap)
	if diags.HasError() {
		return types.Map{}, 500, fmt.Sprintf("Error setting map value: %v", diags)
	}

	return passwords, response.StatusCode, "Success"
}
