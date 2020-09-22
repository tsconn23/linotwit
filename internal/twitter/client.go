package twitter

import (
	"context"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/tsconn23/linotwit/internal/config"
	"log"
	"sync"
)

type TwitterClient struct {
	configuration *config.ConfigInfo
	client *twitter.Client
}

func NewClient(cfg *config.ConfigInfo) *TwitterClient {
	return &TwitterClient{
		configuration: cfg,
	}
}

func (c *TwitterClient) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup) (success bool) {
	wg.Add(1)

	credentials := c.configuration.Credentials //for shorthand
	cfgAuth := oauth1.NewConfig(credentials.ConsumerKey, credentials.ConsumerSecret)
	token := oauth1.NewToken(credentials.AccessToken, credentials.AccessSecret)
	httpClient := cfgAuth.Client(oauth1.NoContext, token)
	c.client = twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
	}
	/*
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}*/

	var terms []string
	terms = append(terms, c.configuration.Subscriptions.Handles...)
	terms = append(terms, c.configuration.Subscriptions.Hashtags...)
	filterParams := &twitter.StreamFilterParams{
		Track: terms,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := c.client.Streams.Filter(filterParams)
	if err != nil {
		success = false
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	return true
}