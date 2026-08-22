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
)

var _ resource.Resource = &followResource{}

type followResource struct {
	client *xapiscraper.XapiScraper
}

type followModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	Cookie   types.String `tfsdk:"cookie"`
}

func NewFollowResource() resource.Resource {
	return &followResource{}
}

func (r *followResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_follow"
}

func (r *followResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Follow an X account via `POST /twitter/user/follow`. Destroy unfollows.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Screen name to follow, without `@`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cookie": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Twitter cookie or `auth_token` for the acting account.",
			},
		},
	}
}

func (r *followResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *followResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan followModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Users.Follow(ctx, components.FollowUserQuery{
		Username: plan.Username.ValueString(),
		Cookie:   plan.Cookie.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI follow failed", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Username.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *followResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state followModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *followResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan followModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(plan.Username.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *followResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state followModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Users.Unfollow(ctx, components.UnfollowUserQuery{
		Username: state.Username.ValueString(),
		Cookie:   state.Cookie.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI unfollow failed", err.Error())
	}
}
