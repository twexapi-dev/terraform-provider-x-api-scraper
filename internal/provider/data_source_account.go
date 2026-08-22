package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	xapiscraper "github.com/twexapi-dev/x-api-scraper-go"
)

var _ datasource.DataSource = &accountDataSource{}

type accountDataSource struct {
	client *xapiscraper.XapiScraper
}

type accountModel struct {
	ID       types.String `tfsdk:"id"`
	Code     types.Int64  `tfsdk:"code"`
	Message  types.String `tfsdk:"message"`
	DataJSON types.String `tfsdk:"data_json"`
}

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Read the current TwexAPI account balance via `GET /balance`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static id for this data source.",
			},
			"code":    schema.Int64Attribute{Computed: true},
			"message": schema.StringAttribute{Computed: true},
			"data_json": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Raw balance payload as JSON.",
			},
		},
	}
}

func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	res, err := d.client.Account.Balance(ctx)
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI account balance failed", err.Error())
		return
	}
	if res.StandardResponse == nil {
		resp.Diagnostics.AddError("TwexAPI account balance failed", "empty response")
		return
	}

	dataJSON := ""
	if payload, ok := res.StandardResponse.Data.GetOrZero(); ok {
		raw, err := json.Marshal(payload)
		if err != nil {
			resp.Diagnostics.AddError("TwexAPI account balance failed", err.Error())
			return
		}
		dataJSON = string(raw)
	}

	state := accountModel{
		ID:       types.StringValue("account"),
		Code:     types.Int64Value(res.StandardResponse.Code),
		Message:  types.StringValue(res.StandardResponse.Msg),
		DataJSON: types.StringValue(dataJSON),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
