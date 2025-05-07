package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func logger(data interface{}) error {

	strData := fmt.Sprintf("%v", data)

	file, err := os.OpenFile("securden_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(strData + "\n"); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func isServerReachable(serverURL string) bool {
	client := createInsecureClient()
	resp, err := client.Get(serverURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 400
}

func isValidURL(input string) bool {
	parsedURL, err := url.Parse(input)
	if err != nil || (parsedURL.Scheme != "https" && parsedURL.Scheme != "http") || parsedURL.Host == "" {
		return false
	}

	if strings.HasSuffix(parsedURL.Path, "/") {
		return false
	}

	host := parsedURL.Hostname()
	portStr := parsedURL.Port()

	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			return false
		}
	}

	hostnameRegex := `^(localhost|([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}|(\d{1,3}\.){3}\d{1,3})$`
	matched, _ := regexp.MatchString(hostnameRegex, host)
	return matched
}

func isValidPEMFile(filePath string) bool {
	if filePath == "" || strings.ToLower(filePath) == "none" {
		return false
	}

	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return false
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return false
	}

	_, err = x509.ParseCertificate(block.Bytes)
	return err == nil
}

func setParam(params map[string]any, key string, value attr.Value) {
	if !value.IsNull() && !value.IsUnknown() {
		switch v := value.(type) {
		case types.String:
			if v.ValueString() != "" {
				params[key] = v.ValueString()
			}
		case types.Int64:
			params[key] = v.ValueInt64()
		case types.Int32:
			params[key] = v.ValueInt32()
		case types.Bool:
			params[key] = v.ValueBool()
		}
	}
}

func fetchSSLCertificate(serverURL string) (*x509.Certificate, error) {
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	host := parsedURL.Host
	if parsedURL.Port() == "" {
		host += default_port
	}

	// Custom TLS config that skips hostname verification
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// do nothing here to skip hostname verification
			return nil
		},
	}

	conn, err := tls.Dial("tcp", host, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	return certs[0], nil
}

func createSecureClient(cert *x509.Certificate) *http.Client {
	certPool := x509.NewCertPool()
	certPool.AddCert(cert)

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func createInsecureClient() *http.Client {
	fmt.Sprintf("Creating insecure HTTP client")
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	return &http.Client{
		Transport: transport,
	}
}

func readPEMFile(filePath string) (*x509.Certificate, error) {
	pemData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM file")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert, nil
}

func raise_request(params map[string]any, apiURL string, method string) ([]byte, error) {
	pattern := regexp.MustCompile("^https")
	var client *http.Client
	var err error

	if pattern.MatchString(SecurdenServerURL) {
		if len(SecurdenCertificate) == 0 {
			cert, certErr := fetchSSLCertificate(SecurdenServerURL)
			if certErr != nil {
				client = createInsecureClient()
			} else {
				client = createSecureClient(cert)
			}
		} else if filepath.IsAbs(SecurdenCertificate) {
			cert, certErr := readPEMFile(SecurdenCertificate)
			if certErr != nil {
				client = createInsecureClient()
			} else {
				client = createSecureClient(cert)
			}
		} else {
			client = createInsecureClient()
		}
	} else {
		client = &http.Client{}
	}

	apiURL = SecurdenServerURL + apiURL

	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	if method == http.MethodGet {
		q := reqURL.Query()
		for key, value := range params {
			switch v := value.(type) {
			case string:
				if v != "" {
					q.Add(key, v)
				}
			case int, int64, float64:
				q.Add(key, fmt.Sprintf("%v", v))
			case bool:
				q.Add(key, strconv.FormatBool(v))
			case []int64:
				for _, num := range v {
					q.Add(key, fmt.Sprintf("%d", num))
				}
			default:
			}
		}
		reqURL.RawQuery = q.Encode()
	}

	var apiRequest *http.Request

	if method != http.MethodGet {
		if ids, exists := params["account_ids"]; exists {
			if tfIDs, ok := ids.([]types.Int64); ok {
				var int64IDs []int64
				for _, id := range tfIDs {
					if !id.IsNull() && !id.IsUnknown() {
						int64IDs = append(int64IDs, id.ValueInt64())
					}
				}
				params["account_ids"] = int64IDs
			}
		}

		requestBody, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize request body: %v", err)
		}

		apiRequest, err = http.NewRequest(method, reqURL.String(), bytes.NewBuffer(requestBody))
		apiRequest.Header.Set("Content-Type", "application/json")
	} else {
		apiRequest, err = http.NewRequest(method, reqURL.String(), nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	apiRequest.Header.Set(const_authtoken, SecurdenAuthToken)

	resp, err := client.Do(apiRequest)
	if err != nil {
		client = createInsecureClient()
		resp, err = client.Do(apiRequest)
		if err != nil {
			return nil, fmt.Errorf("request failed even with insecure client: %v", err)
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return body, nil
}

func get_account(ctx context.Context, account_id int64, account_name, account_title, account_type, ticket_id, reason string) (AccountModel, int, string) {
	var account AccountModel
	params := make(map[string]any)
	if account_id != 0 {
		setParam(params, "account_id", types.Int64Value(account_id))
	}
	setParam(params, "account_name", types.StringValue(account_name))
	setParam(params, "account_title", types.StringValue(account_title))
	setParam(params, "account_type", types.StringValue(account_type))
	setParam(params, "ticket_id", types.StringValue(ticket_id))
	setParam(params, "reason", types.StringValue(reason))
	body, err := raise_request(params, "/secretsmanagement/get_account", GET)
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return account, 500, fmt.Sprintf("Error parsing response JSON: %v", err)
	}
	statusCode, ok := response["status_code"].(float64)
	if !ok {
		return account, 500, "Missing or invalid status_code in response"
	}
	if statusCode != 200 {
		if errMsg, exists := response["error"].(map[string]interface{}); exists {
			if msg, ok := errMsg["message"].(string); ok {
				return account, int(statusCode), msg
			}
		}
		if msg, ok := response["message"].(string); ok {
			return account, int(statusCode), msg
		}
		return account, int(statusCode), "Unknown error"
	}
	accountData := make(map[string]attr.Value)
	for key, value := range response {
		switch v := value.(type) {
		case string:
			accountData[key] = types.StringValue(v)
		case float64:
			accountData[key] = types.StringValue(fmt.Sprintf("%v", v))
		case map[string]interface{}:
			nestedData := make(map[string]attr.Value)
			for nestedKey, nestedValue := range v {
				if nestedStr, ok := nestedValue.(string); ok {
					nestedData[nestedKey] = types.StringValue(nestedStr)
				} else {
					nestedData[nestedKey] = types.StringValue(fmt.Sprintf("%v", nestedValue))
				}
			}
			accountData[key], _ = types.MapValue(types.StringType, nestedData)
		default:
			accountData[key] = types.StringValue(fmt.Sprintf("%v", value))
		}
	}
	account.Account, _ = types.MapValue(types.StringType, accountData)
	return account, int(statusCode), "Success"
}

func get_accounts(ctx context.Context, params map[string]any) (map[string]map[string]string, int, string) {
	var accounts_data = make(map[string]any)
	var null map[string]map[string]string

	body, err := raise_request(params, "/secretsmanagement/get_accounts", POST)
	if err != nil {
		return null, 500, fmt.Sprintf("Error in API call: %v", err)
	}

	err = json.Unmarshal(body, &accounts_data)
	if err != nil {
		return null, 500, fmt.Sprintf("Error parsing response: %v", err)
	}

	processedAccounts := make(map[string]map[string]string)

	for key, value := range accounts_data {
		accountMap, ok := value.(map[string]any)
		if !ok {
			continue
		}

		processedEntry := make(map[string]string)
		for k, v := range accountMap {
			if v == nil {
				processedEntry[k] = ""
			} else {
				processedEntry[k] = fmt.Sprintf("%v", v)
			}
		}

		processedAccounts[key] = processedEntry
	}

	return processedAccounts, 200, "Success"
}

func add_account_function(ctx context.Context, params map[string]any) (AddAccountModel, int, string) {
	var account AddAccountModel
	body, err := raise_request(params, "/api/add_account", POST)
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}

	var response struct {
		ID         int64  `json:"ID"`
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Error      struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return account, 500, fmt.Sprintf("Failed to parse response: %v", err)
	}
	if response.StatusCode != 200 && response.StatusCode != 0 {
		errorMessage := response.Message
		if response.Error.Message != "" {
			errorMessage = response.Error.Message
		}
		return account, response.StatusCode, errorMessage
	}
	account.ID = types.Int64Value(response.ID)
	account.Message = types.StringValue(response.Message)
	return account, response.StatusCode, response.Message
}

func delete_accounts_function(ctx context.Context, params map[string]any) (DeleteAccountsModel, int, string) {
	var account DeleteAccountsModel
	body, err := raise_request(params, "/api/delete_accounts", DELETE)
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}

	var response map[string]any
	err = json.Unmarshal(body, &response)
	if err != nil {
		return account, 500, fmt.Sprintf("Error parsing response: %v", err)
	}

	if message, ok := response["message"].(string); ok {
		account.Message = types.StringValue(message)
	}

	if deletedIDs, ok := response["IDs deleted successfully"].([]any); ok {
		var idList []types.Int64
		for _, id := range deletedIDs {
			if idFloat, ok := id.(float64); ok {
				idList = append(idList, types.Int64Value(int64(idFloat)))
			}
		}
		account.DeletedAccounts = idList
	}

	return account, 200, account.Message.ValueString()
}

func edit_account_function(ctx context.Context, params map[string]any) (EditAccountModel, int, string) {
	var account EditAccountModel
	body, err := raise_request(params, "/api/edit_account", PUT)
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}
	var response struct {
		ID         int64  `json:"ID"`
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Error      struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return account, 500, fmt.Sprintf("Failed to parse response: %v", err)
	}
	if response.StatusCode != 200 && response.StatusCode != 0 {
		errorMessage := response.Message
		if response.Error.Message != "" {
			errorMessage = response.Error.Message
		}
		return account, response.StatusCode, errorMessage
	}
	account.Message = types.StringValue(response.Message)
	return account, response.StatusCode, response.Message
}
