package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	xapiscraper "github.com/twexapi-dev/x-api-scraper-go"
	"github.com/twexapi-dev/x-api-scraper-go/models/operations"
)

var _ datasource.DataSource = &userAboutDataSource{}

type userAboutDataSource struct {
	client *xapiscraper.XapiScraper
}

type userAboutModel struct {
	ScreenName     types.String `tfsdk:"screen_name"`
	UserID         types.String `tfsdk:"user_id"`
	Name           types.String `tfsdk:"name"`
	Avatar         types.String `tfsdk:"avatar"`
	CreatedAt      types.String `tfsdk:"created_at"`
	IsBlueVerified types.Bool   `tfsdk:"is_blue_verified"`
}

func NewUserAboutDataSource() datasource.DataSource {
	return &userAboutDataSource{}
}

func (d *userAboutDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_about"
}

func (d *userAboutDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Read an X profile by screen name via `GET /twitter/{screen_name}/about`.",
		Attributes: map[string]schema.Attribute{
			"screen_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "X/Twitter screen name without `@`.",
			},
			"user_id":          schema.StringAttribute{Computed: true},
			"name":             schema.StringAttribute{Computed: true},
			"avatar":           schema.StringAttribute{Computed: true},
			"created_at":       schema.StringAttribute{Computed: true},
			"is_blue_verified": schema.BoolAttribute{Computed: true},
		},
	}
}

func (d *userAboutDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := clientFrom(req.ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected a TwexAPI client.")
		return
	}
	d.client = client
}

func (d *userAboutDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userAboutModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.Users.GetAbout(ctx, operations.UsersGetAboutRequest{
		ScreenName: data.ScreenName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI user about failed", err.Error())
		return
	}
	if res.TwitterUserAboutStandardResponse == nil {
		resp.Diagnostics.AddError("TwexAPI user about failed", "empty response")
		return
	}
	about, ok := res.TwitterUserAboutStandardResponse.Data.GetOrZero()
	if !ok {
		resp.Diagnostics.AddError("TwexAPI user about failed", "missing profile data")
		return
	}

	data.UserID = types.StringValue(about.UserID)
	data.Name = types.StringValue(about.Name)
	data.Avatar = types.StringValue(about.Avatar)
	data.CreatedAt = types.StringValue(about.CreatedAt)
	data.IsBlueVerified = types.BoolValue(about.IsBlueVerified)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
