package youtube

import (
	"strings"

	"google.golang.org/api/youtube/v3"
)

type SearchResult struct {
	*youtube.SearchResult
	Statistics *youtube.VideoStatistics
	Details    *youtube.VideoContentDetails
	Channel    *youtube.ChannelStatistics
}

func (client *Client) Search(q string, count int) ([]SearchResult, error) {
	search, err := client.Service.
		Search.
		List([]string{"id", "snippet"}).
		RegionCode("RU").
		MaxResults(int64(count)).
		Q(q).
		Do()
	if err != nil {
		return nil, err
	}

	var videos, channels []string
	for _, item := range search.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos = append(videos, item.Id.VideoId)
		case "youtube#channel":
			channels = append(channels, item.Id.ChannelId)
		}
	}

	videoStatistics, err := client.Service.
		Videos.
		List([]string{"statistics", "contentDetails"}).
		Id(strings.Join(videos, ",")).
		Do()
	if err != nil {
		return nil, err
	}

	channelStatistics, err := client.Service.
		Channels.
		List([]string{"statistics"}).
		Id(strings.Join(channels, ",")).
		Do()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range search.Items {
		r := SearchResult{
			SearchResult: item,
		}

		switch item.Id.Kind {
		case "youtube#video":
			for _, s := range videoStatistics.Items {
				if s.Id == item.Id.VideoId {
					r.Statistics = s.Statistics
					r.Details = s.ContentDetails
					break
				}
			}
		case "youtube#channel":
			for _, s := range channelStatistics.Items {
				if s.Id == item.Id.ChannelId {
					r.Channel = s.Statistics
					break
				}
			}
		}

		results = append(results, r)
	}
	return results, nil
}
