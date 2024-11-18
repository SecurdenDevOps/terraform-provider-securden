package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SecurdenDataSource{}

func account_data_source() datasource.DataSource {
	return &SecurdenDataSource{}
}

type SecurdenDataSource struct {
	client *http.Client
}

type SecurdenDataSourceModel struct {
	AccountID         types.Int64  `tfsdk:"account_id"`
	AccountName       types.String `tfsdk:"account_name"`
	AccountTitle      types.String `tfsdk:"account_title"`
	Password          types.String `tfsdk:"password"`
	KeyField          types.String `tfsdk:"key_field"`
	KeyValue          types.String `tfsdk:"key_value"`
	PrivateKey        types.String `tfsdk:"private_key"`
	PuTTYPrivateKey   types.String `tfsdk:"putty_private_key"`
	Passphrase        types.String `tfsdk:"passphrase"`
	PPKPassphrase     types.String `tfsdk:"ppk_passphrase"`
	Address           types.String `tfsdk:"address"`
	ClientID          types.String `tfsdk:"client_id"`
	ClientSecret      types.String `tfsdk:"client_secret"`
	AccountAlias      types.String `tfsdk:"account_alias"`
	AccountFile       types.String `tfsdk:"account_file"`
	OracleSID         types.String `tfsdk:"oracle_sid"`
	OracleServiceName types.String `tfsdk:"oracle_service_name"`
	DefaultDatabase   types.String `tfsdk:"default_database"`
	Port              types.String `tfsdk:"port"`
}

func (d *SecurdenDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keyvalue"
}

func (d *SecurdenDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Securden data source",

		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the account",
				Optional:            true,
				Computed:            true,
			},
			"account_name": schema.StringAttribute{
				MarkdownDescription: "Name of the account",
				Optional:            true,
				Computed:            true,
			},
			"account_title": schema.StringAttribute{
				MarkdownDescription: "Title of the account",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password of the account",
				Computed:            true,
			},
			"key_field": schema.StringAttribute{
				MarkdownDescription: "Key field for the required field",
				Optional:            true,
			},
			"key_value": schema.StringAttribute{
				MarkdownDescription: "Key value of the required field",
				Computed:            true,
			},
			"private_key": schema.StringAttribute{
				Computed: true,
			},
			"putty_private_key": schema.StringAttribute{
				Computed: true,
			},
			"passphrase": schema.StringAttribute{
				MarkdownDescription: "Passphrase for Private Key",
				Computed:            true,
			},
			"ppk_passphrase": schema.StringAttribute{
				MarkdownDescription: "Passphrase for PuTTY Private Key",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				Computed: true,
			},
			"client_id": schema.StringAttribute{
				Computed: true,
			},
			"client_secret": schema.StringAttribute{
				Computed: true,
			},
			"account_alias": schema.StringAttribute{
				Computed: true,
			},
			"account_file": schema.StringAttribute{
				Computed: true,
			},
			"oracle_sid": schema.StringAttribute{
				Computed: true,
			},
			"oracle_service_name": schema.StringAttribute{
				Computed: true,
			},
			"default_database": schema.StringAttribute{
				Computed: true,
			},
			"port": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *SecurdenDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*http.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *SecurdenDataSource) Create(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account SecurdenDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	account_id := account.AccountID.String()
	account_name := account.AccountName.ValueString()
	account_title := account.AccountTitle.ValueString()
	account_field := account.KeyField.ValueString()
	data, code, message := get_account(ctx, account_id, account_name, account_title, account_field)
	if code != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *SecurdenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account SecurdenDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	account_id := account.AccountID.String()
	account_name := account.AccountName.ValueString()
	account_title := account.AccountTitle.ValueString()
	account_field := account.KeyField.ValueString()
	data, code, message := get_account(ctx, account_id, account_name, account_title, account_field)
	if code != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
