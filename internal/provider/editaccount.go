package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &EditAccount{}

func edit_account() datasource.DataSource {
	return &EditAccount{}
}

type EditAccount struct {
	client *http.Client
}

type EditAccountModel struct {
	AccountID                 types.Int64  `tfsdk:"account_id"`
	AccountName               types.String `tfsdk:"account_name"`
	AccountTitle              types.String `tfsdk:"account_title"`
	AccountType               types.String `tfsdk:"account_type"`
	IPAddress                 types.String `tfsdk:"ipaddress"`
	Notes                     types.String `tfsdk:"notes"`
	Tags                      types.String `tfsdk:"tags"`
	FolderID                  types.Int64  `tfsdk:"folder_id"`
	OverwriteAdditionalFields types.Bool   `tfsdk:"overwrite_additional_fields"`
	AccountExpirationDate     types.String `tfsdk:"account_expiration_date"`
	DistinguishedName         types.String `tfsdk:"distinguished_name"`
	AccountAlias              types.String `tfsdk:"account_alias"`
	DomainName                types.String `tfsdk:"domain_name"`
	Message                   types.String `tfsdk:"message"`
}

func (d *EditAccount) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_edit_account"
}

func (d *EditAccount) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Defines the structure for managing account updates in Securden.",

		Attributes: map[string]schema.Attribute{
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "Unique identifier of the account.",
				Required:            true,
			},
			"account_type": schema.StringAttribute{
				MarkdownDescription: "Specifies the type of the account.",
				Optional:            true,
			},
			"account_title": schema.StringAttribute{
				MarkdownDescription: "The title associated with the account.",
				Optional:            true,
			},
			"account_name": schema.StringAttribute{
				MarkdownDescription: "The name associated with the account.",
				Optional:            true,
			},
			"ipaddress": schema.StringAttribute{
				MarkdownDescription: "The IP address of the account (if applicable).",
				Optional:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Additional notes related to the account.",
				Optional:            true,
			},
			"tags": schema.StringAttribute{
				MarkdownDescription: "Tags associated with the account.",
				Optional:            true,
			},
			"folder_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the folder where the account belongs to.",
				Optional:            true,
			},
			"account_expiration_date": schema.StringAttribute{
				MarkdownDescription: "The expiration date of the account (format: DD/MM/YYYY).",
				Optional:            true,
			},
			"distinguished_name": schema.StringAttribute{
				MarkdownDescription: "Required for LDAP domain accounts.",
				Optional:            true,
			},
			"overwrite_additional_fields": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether additional fields should be overwritten (true/false).",
				Optional:            true,
			},
			"account_alias": schema.StringAttribute{
				MarkdownDescription: "Required for AWS IAM accounts.",
				Optional:            true,
			},
			"domain_name": schema.StringAttribute{
				MarkdownDescription: "Required for Google Workspace accounts.",
				Optional:            true,
			},
			"message": schema.StringAttribute{
				MarkdownDescription: "Response message indicating the result of the operation.",
				Computed:            true,
			},
		},
	}
}

func (d *EditAccount) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EditAccount) Create(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account EditAccountModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	params := make(map[string]any)
	setParam(params, "account_id", account.AccountID)
	setParam(params, "account_title", account.AccountTitle)
	setParam(params, "account_name", account.AccountName)
	setParam(params, "account_type", account.AccountType)
	setParam(params, "ipaddress", account.IPAddress)
	setParam(params, "notes", account.Notes)
	setParam(params, "tags", account.Tags)
	setParam(params, "folder_id", account.FolderID)
	setParam(params, "overwrite_additional_fields", account.OverwriteAdditionalFields)
	setParam(params, "account_expiration_date", account.AccountExpirationDate)
	setParam(params, "distinguished_name", account.DistinguishedName)
	setParam(params, "account_alias", account.AccountAlias)
	setParam(params, "domain_name", account.DomainName)
	edit_account, code, message := edit_account_function(ctx, params)
	if code != 200 && code != 0 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &edit_account)...)
}

func (d *EditAccount) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var account EditAccountModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &account)...)
	params := make(map[string]any)
	setParam(params, "account_id", account.AccountID)
	setParam(params, "account_title", account.AccountTitle)
	setParam(params, "account_name", account.AccountName)
	setParam(params, "account_type", account.AccountType)
	setParam(params, "ipaddress", account.IPAddress)
	setParam(params, "notes", account.Notes)
	setParam(params, "tags", account.Tags)
	setParam(params, "folder_id", account.FolderID)
	setParam(params, "overwrite_additional_fields", account.OverwriteAdditionalFields)
	setParam(params, "account_expiration_date", account.AccountExpirationDate)
	setParam(params, "distinguished_name", account.DistinguishedName)
	setParam(params, "account_alias", account.AccountAlias)
	setParam(params, "domain_name", account.DomainName)
	edit_account, code, message := edit_account_function(ctx, params)
	if code != 200 && code != 0 {
		resp.Diagnostics.AddWarning(fmt.Sprintf("%d - %s", code, message), "")
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &edit_account)...)
}
