package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	xapiscraper "github.com/twexapi-dev/x-api-scraper-go"
	"github.com/twexapi-dev/x-api-scraper-go/models/components"
	"github.com/twexapi-dev/x-api-scraper-go/optionalnullable"
)

var _ resource.Resource = &dmResource{}

type dmResource struct {
	client *xapiscraper.XapiScraper
}

type dmModel struct {
	ID        types.String `tfsdk:"id"`
	Recipient types.String `tfsdk:"recipient"`
	Text      types.String `tfsdk:"text"`
	Cookie    types.String `tfsdk:"cookie"`
	MediaURL  types.String `tfsdk:"media_url"`
	VideoURL  types.String `tfsdk:"video_url"`
	MessageID types.String `tfsdk:"message_id"`
	SentAt    types.String `tfsdk:"sent_at"`
}

func NewDmResource() resource.Resource {
	return &dmResource{}
}

func (r *dmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dm"
}

func (r *dmResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Send a DM via `POST /v3/twitter/send-dm`. Destroy removes Terraform state only; the API cannot unsend.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"recipient": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User id, `@handle`, or an existing group id (`g...`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"text": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Message body.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cookie": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Twitter cookie or `auth_token` for the acting account.",
			},
			"media_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Public image URL (max 1). Mutually exclusive with `video_url`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"video_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Public video URL. Mutually exclusive with `media_url`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"message_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Sent message id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sent_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Send timestamp from the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *dmResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := clientFrom(req.ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected a TwexAPI client.")
		return
	}
	r.client = client
}

func (r *dmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dmModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Dm.Send(ctx, dmQueryFrom(plan))
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI send DM failed", err.Error())
		return
	}
	if res.SendDmResponse == nil {
		resp.Diagnostics.AddError("TwexAPI send DM failed", "empty response")
		return
	}

	messageID := plan.Recipient.ValueString()
	sentAt := types.StringNull()
	if data := res.SendDmResponse.Data; data != nil {
		if data.ID != "" {
			messageID = data.ID
		}
		if data.Time != "" {
			sentAt = types.StringValue(data.Time)
		}
	}

	plan.MessageID = types.StringValue(messageID)
	plan.SentAt = sentAt
	plan.ID = types.StringValue(messageID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dmModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dmModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"DM not unsent",
		"TwexAPI has no unsend-DM route. Destroy only removes this resource from Terraform state.",
	)
}

func dmQueryFrom(plan dmModel) components.SendDmV3Query {
	query := components.SendDmV3Query{
		Recipient: plan.Recipient.ValueString(),
		Text:      plan.Text.ValueString(),
		Cookie:    plan.Cookie.ValueString(),
	}
	if !plan.MediaURL.IsNull() && !plan.MediaURL.IsUnknown() && plan.MediaURL.ValueString() != "" {
		urls := []string{plan.MediaURL.ValueString()}
		query.MediaUrls = optionalnullable.From(&urls)
	}
	if !plan.VideoURL.IsNull() && !plan.VideoURL.IsUnknown() && plan.VideoURL.ValueString() != "" {
		video := plan.VideoURL.ValueString()
		query.VideoURL = optionalnullable.From(&video)
	}
	return query
}
