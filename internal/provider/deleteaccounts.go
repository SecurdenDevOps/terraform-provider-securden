package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DeleteAccounts{}

func delete_accounts() datasource.DataSource {
	return &DeleteAccounts{}
}

type DeleteAccounts struct {
	client *http.Client
}

type DeleteAccountsModel struct {
	AccountIDs        []types.Int64 `tfsdk:"account_ids"`
	Reason            types.String  `tfsdk:"reason"`
	DeletePermanently types.Bool    `tfsdk:"delete_permanently"`
	Message           types.String  `tfsdk:"message"`
	DeletedAccounts   []types.Int64 `tfsdk:"deleted_accounts"`
}

func (d *DeleteAccounts) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delete_accounts"
}

func (d *DeleteAccounts) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Defines the structure for managing account deletions in Securden.",

		Attributes: map[string]schema.Attribute{
			"account_ids": schema.ListAttribute{
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of account IDs to be deleted.",
				Required:            true,
			},
			"reason": schema.StringAttribute{
				MarkdownDescription: "Reason for deleting the accounts.",
				Optional:            true,
			},
			"delete_permanently": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the accounts should be permanently deleted (true/false).",
				Optional:            true,
			},
			"message": schema.StringAttribute{
				MarkdownDescription: "Response message indicating the result of the deletion operation.",
				Computed:            true,
			},
			"deleted_accounts": schema.ListAttribute{
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of account IDs that were successfully deleted.",
				Computed:            true,
			},
		},
	}
}

func (d *DeleteAccounts) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DeleteAccounts) Create(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account DeleteAccountsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	params := make(map[string]any)
	params["account_ids"] = account.AccountIDs
	if account.Reason.ValueString() != "" {
		params["reason"] = account.Reason
	}
	if account.DeletePermanently.ValueBool() {
		params["delete_permanently"] = account.DeletePermanently
	}
	delete_accounts, code, message := delete_accounts_function(ctx, params)
	if code != 200 && code != 0 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &delete_accounts)...)
}

func (d *DeleteAccounts) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account DeleteAccountsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	params := make(map[string]any)
	params["account_ids"] = account.AccountIDs
	if account.Reason.ValueString() != "" {
		params["reason"] = account.Reason
	}
	if account.DeletePermanently.ValueBool() {
		params["delete_permanently"] = account.DeletePermanently
	}
	delete_accounts, code, message := delete_accounts_function(ctx, params)
	if code != 200 && code != 0 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &delete_accounts)...)
}
