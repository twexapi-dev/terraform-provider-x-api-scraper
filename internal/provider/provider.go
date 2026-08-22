package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	xapiscraper "github.com/twexapi-dev/x-api-scraper-go"
)

var (
	_ provider.Provider = &XapiScraperProvider{}
)

type XapiScraperProvider struct {
	version string
}

type providerModel struct {
	BearerAuth types.String `tfsdk:"bearer_auth"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &XapiScraperProvider{version: version}
	}
}

func (p *XapiScraperProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "x-api-scraper"
	resp.Version = p.version
}

func (p *XapiScraperProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for TwexAPI Twitter search, followers, DMs, and X automation. Not affiliated with X Corp.",
		Attributes: map[string]schema.Attribute{
			"bearer_auth": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "TwexAPI API key. Falls back to `X_API_SCRAPER_KEY` or `X_API_SCRAPER_BEARER_AUTH`.",
			},
		},
	}
}

func (p *XapiScraperProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bearer := config.BearerAuth.ValueString()
	if bearer == "" {
		bearer = os.Getenv("X_API_SCRAPER_KEY")
	}
	if bearer == "" {
		bearer = os.Getenv("X_API_SCRAPER_BEARER_AUTH")
	}
	if bearer == "" {
		resp.Diagnostics.AddError(
			"Missing TwexAPI API key",
			"Set bearer_auth on the provider, or export X_API_SCRAPER_KEY.",
		)
		return
	}

	client := xapiscraper.New(xapiscraper.WithSecurity(bearer))
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *XapiScraperProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFollowResource,
		NewTweetResource,
	}
}

func (p *XapiScraperProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserAboutDataSource,
		NewAccountDataSource,
	}
}

func clientFrom(data any) (*xapiscraper.XapiScraper, bool) {
	client, ok := data.(*xapiscraper.XapiScraper)
	return client, ok && client != nil
}
