package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Accounts{}

func accounts() datasource.DataSource {
	return &Accounts{}
}

type Accounts struct {
	client *http.Client
}

type AccountsModel struct {
	AccountIDs []types.Int64                `tfsdk:"account_ids"`
	Accounts   map[string]map[string]string `tfsdk:"accounts"`
}

func (d *Accounts) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accounts"
}

func (d *Accounts) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves multiple account details from Securden.",

		Attributes: map[string]schema.Attribute{
			"account_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				MarkdownDescription: "A list of account IDs to fetch details for.",
				Required:            true,
			},
			"accounts": schema.MapAttribute{
				ElementType:         types.MapType{ElemType: types.StringType},
				Computed:            true,
				MarkdownDescription: "A map containing multiple account details, where each key represents an account ID and the value is a map of account attributes.",
			},
		},
	}
}

func (d *Accounts) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *Accounts) Create(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var accounts AccountsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &accounts)...)

	params := make(map[string]any)
	params["account_ids"] = accounts.AccountIDs

	accountsData, code, message := get_accounts(ctx, params)
	if code != 200 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	accounts.Accounts = accountsData

	resp.Diagnostics.Append(resp.State.Set(ctx, &accounts)...)
}

func (d *Accounts) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var accounts AccountsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &accounts)...)

	params := make(map[string]any)
	params["account_ids"] = accounts.AccountIDs

	accountsData, code, message := get_accounts(ctx, params)
	if code != 200 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	accounts.Accounts = accountsData

	resp.Diagnostics.Append(resp.State.Set(ctx, &accounts)...)
}
