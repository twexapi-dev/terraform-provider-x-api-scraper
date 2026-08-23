package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	xapiscraper "github.com/twexapi-dev/x-api-scraper-go"
)

var _ datasource.DataSource = &tweetDataSource{}

type tweetDataSource struct {
	client *xapiscraper.XapiScraper
}

type tweetDataModel struct {
	TweetID       types.String `tfsdk:"tweet_id"`
	Text          types.String `tfsdk:"text"`
	CreatedAt     types.String `tfsdk:"created_at"`
	ScreenName    types.String `tfsdk:"screen_name"`
	UserID        types.String `tfsdk:"user_id"`
	FavoriteCount types.Int64  `tfsdk:"favorite_count"`
	RetweetCount  types.Int64  `tfsdk:"retweet_count"`
	ReplyCount    types.Int64  `tfsdk:"reply_count"`
	DataJSON      types.String `tfsdk:"data_json"`
}

func NewTweetDataSource() datasource.DataSource {
	return &tweetDataSource{}
}

func (d *tweetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tweet"
}

func (d *tweetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Read a tweet by id via tweet lookup. Maps to the competitor `x_tweet` data source.",
		Attributes: map[string]schema.Attribute{
			"tweet_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Tweet id to look up.",
			},
			"text":           schema.StringAttribute{Computed: true},
			"created_at":     schema.StringAttribute{Computed: true},
			"screen_name":    schema.StringAttribute{Computed: true},
			"user_id":        schema.StringAttribute{Computed: true},
			"favorite_count": schema.Int64Attribute{Computed: true},
			"retweet_count":  schema.Int64Attribute{Computed: true},
			"reply_count":    schema.Int64Attribute{Computed: true},
			"data_json": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Raw tweet payload as JSON.",
			},
		},
	}
}

func (d *tweetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tweetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data tweetDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.Tweets.Lookup(ctx, []string{data.TweetID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI tweet lookup failed", err.Error())
		return
	}
	if res.GetTweetsByIDResponse == nil {
		resp.Diagnostics.AddError("TwexAPI tweet lookup failed", "empty response")
		return
	}
	tweets, ok := res.GetTweetsByIDResponse.Data.GetOrZero()
	if !ok || len(tweets) == 0 {
		resp.Diagnostics.AddError("TwexAPI tweet lookup failed", "tweet not found")
		return
	}
	tweet := tweets[0]
	data.Text = types.StringValue(tweet.Text)
	data.CreatedAt = optString(tweet.CreatedAt)
	data.FavoriteCount = optInt64(tweet.FavoriteCount)
	data.RetweetCount = optInt64(tweet.RetweetCount)
	data.ReplyCount = optInt64(tweet.ReplyCount)
	if user, ok := tweet.User.GetOrZero(); ok && !tweet.User.IsNull() {
		data.ScreenName = types.StringValue(user.ScreenName)
		data.UserID = types.StringValue(user.ID)
	} else {
		data.ScreenName = types.StringNull()
		data.UserID = types.StringNull()
	}
	raw, err := asJSON(tweet)
	if err != nil {
		resp.Diagnostics.AddError("TwexAPI tweet lookup failed", err.Error())
		return
	}
	data.DataJSON = raw
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
