package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	xapiscraper "github.com/twexapi-dev/x-api-scraper-go"
	"github.com/twexapi-dev/x-api-scraper-go/models/components"
	"github.com/twexapi-dev/x-api-scraper-go/models/operations"
)

var _ resource.Resource = &tweetActionResource{}

type tweetActionKind string

const (
	actionLike     tweetActionKind = "like"
	actionRetweet  tweetActionKind = "retweet"
	actionBookmark tweetActionKind = "bookmark"
)

type tweetActionResource struct {
	client *xapiscraper.XapiScraper
	kind   tweetActionKind
}

type tweetActionModel struct {
	ID      types.String `tfsdk:"id"`
	TweetID types.String `tfsdk:"tweet_id"`
	Cookie  types.String `tfsdk:"cookie"`
}

func NewLikeResource() resource.Resource {
	return &tweetActionResource{kind: actionLike}
}

func NewRetweetResource() resource.Resource {
	return &tweetActionResource{kind: actionRetweet}
}

func NewBookmarkResource() resource.Resource {
	return &tweetActionResource{kind: actionBookmark}
}

func (r *tweetActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + string(r.kind)
}

func (r *tweetActionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("%s a tweet. Destroy reverses the action. Maps to the competitor `x_write` / `x_write_action` surface.", titleKind(r.kind)),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tweet_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Target tweet id.",
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

func (r *tweetActionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *tweetActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tweetActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.apply(ctx, plan, false); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("TwexAPI %s failed", r.kind), err.Error())
		return
	}
	plan.ID = types.StringValue(string(r.kind) + ":" + plan.TweetID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tweetActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tweetActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *tweetActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tweetActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(string(r.kind) + ":" + plan.TweetID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tweetActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tweetActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.apply(ctx, state, true); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("TwexAPI undo %s failed", r.kind), err.Error())
	}
}

func (r *tweetActionResource) apply(ctx context.Context, data tweetActionModel, undo bool) error {
	tweetID := data.TweetID.ValueString()
	cookie := data.Cookie.ValueString()
	switch r.kind {
	case actionLike:
		if undo {
			_, err := r.client.Tweets.Actions.Unlike(ctx, operations.TweetsActionsUnlikeRequest{
				TweetID: tweetID,
				Body:    components.UnfavoriteTweetQuery{Cookie: cookie},
			})
			return err
		}
		_, err := r.client.Tweets.Actions.Like(ctx, operations.TweetsActionsLikeRequest{
			TweetID: tweetID,
			Body:    components.LikeActionBody{Cookie: cookie},
		})
		return err
	case actionRetweet:
		if undo {
			_, err := r.client.Tweets.Actions.Unretweet(ctx, operations.TweetsActionsUnretweetRequest{
				TweetID: tweetID,
				Body:    components.DeleteRetweetQuery{Cookie: cookie},
			})
			return err
		}
		_, err := r.client.Tweets.Actions.Retweet(ctx, operations.TweetsActionsRetweetRequest{
			TweetID: tweetID,
			Body:    components.RetweetActionBody{Cookie: cookie},
		})
		return err
	case actionBookmark:
		if undo {
			_, err := r.client.Tweets.Actions.Unbookmark(ctx, operations.TweetsActionsUnbookmarkRequest{
				TweetID: tweetID,
				Body:    components.DeleteBookmarkQuery{Cookie: cookie},
			})
			return err
		}
		_, err := r.client.Tweets.Actions.Bookmark(ctx, operations.TweetsActionsBookmarkRequest{
			TweetID: tweetID,
			Body:    components.BookmarkBody{Cookie: cookie},
		})
		return err
	default:
		return fmt.Errorf("unknown action %s", r.kind)
	}
}

func titleKind(kind tweetActionKind) string {
	switch kind {
	case actionLike:
		return "Like"
	case actionRetweet:
		return "Retweet"
	case actionBookmark:
		return "Bookmark"
	default:
		return string(kind)
	}
}
