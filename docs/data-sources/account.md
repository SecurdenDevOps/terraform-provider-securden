---
page_title: "securden_account Data Source - terraform-provider-securden"
subcategory: ""
description: |-
  Securden data source
---

# securden_account (Data Source)

Account Data


## Schema

### Optional

- `account_id` (Number) ID of the account
- `account_name` (String) Name of the account
- `account_title` (String) Title of the account
- `key_field` (String) Key field for the required additional field

### Read-Only

- `password` (String) Password of the account
- `account_file` (String) File Content of the account
- `key_value` (String) Value of the additional field
- `address` (String) Address of the account
- `client_id` (String) Client ID of the account
- `client_secret` (String) Client Secret of the account
- `private_key` (String) Private Key of the account
- `passphrase` (String) Passphrase of the Private Key
- `putty_private_key` (String) PuTTY Private Key of the account
- `ppk_passphrase` (String) Passphrase of the PuTTY Private Key
- `default_database` (String) Default Database of the account
- `account_alias` (String) Required for AWS IAM Account
- `oracle_service_name` (String) Oracle Service Name of the account
- `oracle_sid` (String) Oracle SID of the account
- `port` (String) Account Port
