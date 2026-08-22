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

var _ resource.Resource = &tweetResource{}

type tweetResource struct {
	client *xapiscraper.XapiScraper
}

type tweetModel struct {
	ID           types.String `tfsdk:"id"`
	TweetContent types.String `tfsdk:"tweet_content"`
	Username     types.String `tfsdk:"username"`
	Cookie       types.String `tfsdk:"cookie"`
	ReplyTweetID types.String `tfsdk:"reply_tweet_id"`
	TweetID      types.String `tfsdk:"tweet_id"`
}

func NewTweetResource() resource.Resource {
	return &tweetResource{}
}

func (r *tweetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tweet"
}

func (r *tweetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Post a tweet via `POST /twitter/tweets/create`. Destroy deletes that tweet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tweet_content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Tweet text.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Acting account username, used when deleting the tweet.",
			},
			"cookie": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Twitter cookie or `auth_token` for the acting account.",
			},
			"reply_tweet_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Tweet id to reply to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tweet_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Created tweet id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *tweetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *tweetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tweetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := components.PostTweetQuery{
		TweetContent: plan.TweetContent.ValueString(),
		Cookie:       plan.Cookie.ValueString(),
	}
	if !plan.ReplyTweetID.IsNull() && !plan.ReplyTweetID.IsUnknown() {
		replyID := plan.ReplyTweetID.ValueString()
		query.ReplyTweetID = optionalnullable.From(&replyID)
	}

	res, err := r.client.Tweets.Actions.Create(ctx, query)
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI tweet create failed", err.Error())
		return
	}
	if res.PostTweetResponse == nil {
		resp.Diagnostics.AddError("TwexAPI tweet create failed", "empty response")
		return
	}
	posted, ok := res.PostTweetResponse.Data.GetOrZero()
	if !ok {
		resp.Diagnostics.AddError("TwexAPI tweet create failed", "missing tweet data")
		return
	}
	tweetID, ok := posted.TweetID.GetOrZero()
	if !ok || tweetID == "" {
		resp.Diagnostics.AddError("TwexAPI tweet create failed", "missing tweet id")
		return
	}

	plan.TweetID = types.StringValue(tweetID)
	plan.ID = types.StringValue(tweetID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tweetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tweetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *tweetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tweetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tweetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tweetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	targetID := state.TweetID.ValueString()
	_, err := r.client.Tweets.Actions.DeleteBatch(ctx, components.DeleteTweetQuery{
		Cookie:   state.Cookie.ValueString(),
		Username: state.Username.ValueString(),
		TargetID: optionalnullable.From(&targetID),
	})
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI tweet delete failed", err.Error())
	}
}
