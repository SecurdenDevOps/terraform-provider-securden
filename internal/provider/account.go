package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Account{}

func account() datasource.DataSource {
	return &Account{}
}

type Account struct {
	client *http.Client
}

type AccountModel struct {
	AccountID    types.Int64  `tfsdk:"account_id"`
	AccountName  types.String `tfsdk:"account_name"`
	AccountTitle types.String `tfsdk:"account_title"`
	AccountType  types.String `tfsdk:"account_type"`
	Account      types.Map    `tfsdk:"account"`
}

func (d *Account) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *Account) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves account details from Securden.",

		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Unique identifier of the account.",
			},
			"account_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name associated with the account.",
			},
			"account_title": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Title or designation of the account.",
			},
			"account_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specifies the type or category of the account.",
			},
			"account": schema.MapAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "A map containing account attributes as keys and their corresponding values.",
			},
		},
	}
}

func (d *Account) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*http.Client)
	if !ok {
		resp.Diagnostics.AddWarning(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *Account) Create(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account AccountModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	var account_id int64
	account_id = 0
	if !account.AccountID.IsNull() {
		account_id = account.AccountID.ValueInt64()
	}
	account_name := account.AccountName.ValueString()
	account_title := account.AccountTitle.ValueString()
	account_type := account.AccountType.ValueString()
	var data AccountModel
	var code int
	var message string
	data, code, message = get_account(ctx, account_id, account_name, account_title, account_type)
	if code != 200 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *Account) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account AccountModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	var account_id int64
	account_id = 0
	if !account.AccountID.IsNull() {
		account_id = account.AccountID.ValueInt64()
	}
	account_name := account.AccountName.ValueString()
	account_title := account.AccountTitle.ValueString()
	account_type := account.AccountType.ValueString()
	var data AccountModel
	var code int
	var message string
	data, code, message = get_account(ctx, account_id, account_name, account_title, account_type)
	if code != 200 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
