# Securden Terraform Provider

The `securden` provider allows Terraform to interact with Securden's API, enabling you to retrieve account details such as passwords and ports based on an `account_id`, or a combination of `account_name` and `account_title`.

## Terraform provider block and data block naming metrics

The naming convention for the Terraform blocks is based on the following structure:

- **Provider Block**: The provider block name must be `"securden"` as shown below:
  
  ```hcl
  provider "securden" {
    authtoken  = var.authtoken
    server_url = var.server_url
  }
  ```

- **Data Block**: The data block name must be `"securden_keyvalue"`, as this data source retrieves key-value data from the Securden system. Example:

  ```hcl
  data "securden_keyvalue" "account_data" {
    account_id = 2000000004406
  }
  ```

  In the above example:
  - `"securden_keyvalue"`: This is the fixed name of the data block to fetch key-value data from Securden.
  - `"account_data"`: This is the user-defined name for the resource that can be referenced throughout the Terraform configuration. You can replace `"account_data"` with any meaningful name according to your use case.

The naming convention ensures that the Terraform configuration remains consistent and intuitive while interacting with the Securden provider.


## Example Usage

### Provider Configuration

```hcl
provider "securden" {
  authtoken  = var.authtoken
  server_url = var.server_url
}
```

- `authtoken`: The API authentication token used to access Securden.
- `server_url`: The URL for your Securden server.

These values can be set via variables, environment variables, or directly in the provider block.

### Data Source Configuration

To fetch account data from Securden, use the `securden_keyvalue` data block. You can provide either an `account_id` or a combination of `account_name` and `account_title`.

#### Fetching Account by `account_id`:

```hcl
data "securden_keyvalue" "account_data" {
  account_id = 2000000004406
}
```

#### Fetching Account by `account_name` and `account_title`:

```hcl
data "securden_keyvalue" "account_data" {
  account_name  = "my_account"
  account_title = "production"
}
```

> **Note:** If both `account_id` and the combination of `account_name` and `account_title` are provided, `account_id` will take priority.

### Accessing Account Data

Once the data is fetched, you can access the account details as follows:

```hcl
output "password" {
  value = data.securden_keyvalue.account_data.password
}

output "port" {
  value = data.securden_keyvalue.account_data.port
}

output "account_name" {
  value = data.securden_keyvalue.account_data.account_name
}
```

## Argument Reference

### Provider

- `authtoken` (Required): The API token used for authentication with Securden.
- `server_url` (Required): The URL for the Securden server.

### Data Source: `securden_keyvalue`

- `account_id` (Optional): The unique identifier of the account.
- `account_name` (Optional): The name of the account.
- `account_title` (Optional): The title of the account.

> **Note:** If `account_id` is provided, it will take precedence over `account_name` and `account_title`.

## Attributes Reference

The following fields are available and can be accessed after retrieving account data:

- `account_id`: The unique identifier of the account.
- `account_name`: The name of the account.
- `account_title`: The title of the account.
- `password`: The password for the account.
- `key_field`: A field key related to the account.
- `key_value`: The value corresponding to the key field.
- `private_key`: The private key associated with the account.
- `putty_private_key`: The PuTTY-format private key.
- `passphrase`: The passphrase used for the private key.
- `ppk_passphrase`: The passphrase for the PuTTY private key.
- `address`: The address associated with the account.
- `client_id`: The client ID for the account.
- `client_secret`: The client secret for the account.
- `account_alias`: The alias or alternative name for the account.
- `account_file`: The file associated with the account.
- `oracle_sid`: The Oracle SID for the account.
- `oracle_service_name`: The Oracle service name for the account.
- `default_database`: The default database associated with the account.
- `port`: The port number for the account.