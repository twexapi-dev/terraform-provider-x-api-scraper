package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDmQueryFromRequired(t *testing.T) {
	q := dmQueryFrom(dmModel{
		Recipient: types.StringValue("elonmusk"),
		Text:      types.StringValue("hello"),
		Cookie:    types.StringValue("ct0=abc"),
	})
	if q.Recipient != "elonmusk" || q.Text != "hello" || q.Cookie != "ct0=abc" {
		t.Fatalf("unexpected query: %+v", q)
	}
	if q.MediaUrls != nil || q.VideoURL != nil || q.Pin != nil || q.Proxy != nil {
		t.Fatalf("optional fields should stay unset: %+v", q)
	}
}

func TestDmQueryFromMedia(t *testing.T) {
	q := dmQueryFrom(dmModel{
		Recipient: types.StringValue("1"),
		Text:      types.StringValue("pic"),
		Cookie:    types.StringValue("cookie"),
		MediaURL:  types.StringValue("https://example.com/a.jpg"),
	})
	urls, ok := q.MediaUrls.GetOrZero()
	if !ok || len(urls) != 1 || urls[0] != "https://example.com/a.jpg" {
		t.Fatalf("media_urls: %+v", q.MediaUrls)
	}
}
