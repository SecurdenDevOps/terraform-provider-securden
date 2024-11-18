package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SecurdenPasswords{}

func accounts_passwords_source() datasource.DataSource {
	return &SecurdenPasswords{}
}

type SecurdenPasswords struct {
	client *http.Client
}

type SecurdenPasswordsModel struct {
	AccountIDs []types.String `tfsdk:"account_ids"`
	Passwords  types.Map      `tfsdk:"passwords"`
}

func (d *SecurdenPasswords) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_passwords"
}

func (d *SecurdenPasswords) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Securden Accounts Passwords",

		Attributes: map[string]schema.Attribute{
			"account_ids": schema.ListAttribute{ // New list attribute
				ElementType:         types.StringType,
				MarkdownDescription: "List of account ids",
				Required:            true,
			},
			"passwords": schema.MapAttribute{
				ElementType: types.StringType, // Using StringAttribute directly
				Computed:    true,
			},
		},
	}
}

func (d *SecurdenPasswords) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SecurdenPasswords) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account SecurdenPasswordsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)

	// Convert the account IDs into strings
	var accountIDs []string
	for _, id := range account.AccountIDs {
		accountIDs = append(accountIDs, id.ValueString())
	}

	// Call the get_passwords function
	passwords, code, message := get_passwords(ctx, accountIDs)
	if code != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("%d - %s", code, message), "")
		return
	}

	// Assign the fetched passwords back to the model
	account.Passwords = passwords

	// Save the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &account)...)
}

func (d *SecurdenPasswords) Create(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account SecurdenPasswordsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)

	// Convert the account IDs into strings
	var accountIDs []string
	for _, id := range account.AccountIDs {
		accountIDs = append(accountIDs, id.ValueString())
	}

	// Call the get_passwords function
	passwords, code, message := get_passwords(ctx, accountIDs)
	if code != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("%d - %s", code, message), "")
		return
	}

	// Assign the fetched passwords back to the model
	account.Passwords = passwords

	// Save the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &account)...)
}
