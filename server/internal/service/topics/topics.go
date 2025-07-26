package topics

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

type TopicService struct {
	client *http.Client
}

type Topic struct {
	Title       string
	Description string
	URL         string
	Source      string
}

func NewTopicService() *TopicService {
	return &TopicService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// cleanText decodes HTML entities and trims whitespace from text
func cleanText(text string) string {
	// Decode HTML entities like &amp;, &lt;, &gt;, &quot;, etc.
	decoded := html.UnescapeString(text)
	// Trim whitespace
	return strings.TrimSpace(decoded)
}

// FetchHackerNewsTop fetches the top story from Hacker News
func (s *TopicService) FetchHackerNewsTop(ctx context.Context) (*Topic, error) {
	// Get top story IDs
	resp, err := s.client.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, fmt.Errorf("fetch HN top stories: %w", err)
	}
	defer resp.Body.Close()

	var storyIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&storyIDs); err != nil {
		return nil, fmt.Errorf("decode HN story IDs: %w", err)
	}

	if len(storyIDs) == 0 {
		return nil, fmt.Errorf("no HN stories found")
	}

	// Get the top story details
	storyURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyIDs[0])
	resp, err = s.client.Get(storyURL)
	if err != nil {
		return nil, fmt.Errorf("fetch HN story: %w", err)
	}
	defer resp.Body.Close()

	var story struct {
		Title string `json:"title"`
		URL   string `json:"url"`
		Score int    `json:"score"`
		By    string `json:"by"`
		ID    int    `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return nil, fmt.Errorf("decode HN story: %w", err)
	}

	fmt.Printf("HackerNews API Response: %+v\n", story)

	// If no URL, it's a text post - link to HN discussion
	if story.URL == "" {
		story.URL = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ID)
	}

	topic := &Topic{
		Title:       cleanText(story.Title),
		Description: fmt.Sprintf("Top HN story with %d points by %s", story.Score, story.By),
		URL:         story.URL,
		Source:      "HackerNews",
	}

	fmt.Printf("HackerNews Topic: %+v\n", topic)

	return topic, nil
}

// FetchRedditWorldNews fetches the top post from r/worldnews
func (s *TopicService) FetchRedditWorldNews(ctx context.Context) (*Topic, error) {
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/worldnews/top.json?limit=1&t=day", nil)
	if err != nil {
		return nil, fmt.Errorf("create Reddit worldnews request: %w", err)
	}
	req.Header.Set("User-Agent", "GoChat/1.0 (by /u/melkeydev)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch Reddit worldnews: %w", err)
	}
	defer resp.Body.Close()

	var redditResp struct {
		Data struct {
			Children []struct {
				Data struct {
					Title     string `json:"title"`
					URL       string `json:"url"`
					Subreddit string `json:"subreddit"`
					Score     int    `json:"score"`
					Permalink string `json:"permalink"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&redditResp); err != nil {
		return nil, fmt.Errorf("decode Reddit response: %w", err)
	}

	if len(redditResp.Data.Children) == 0 {
		return nil, fmt.Errorf("no Reddit posts found")
	}

	post := redditResp.Data.Children[0].Data
	fmt.Printf("Reddit WorldNews API Response: %+v\n", post)

	// Use Reddit URL for discussion
	redditURL := "https://reddit.com" + post.Permalink

	topic := &Topic{
		Title:       cleanText(post.Title),
		Description: fmt.Sprintf("Top world news with %d upvotes", post.Score),
		URL:         redditURL,
		Source:      "Reddit WorldNews",
	}

	fmt.Printf("Reddit WorldNews Topic: %+v\n", topic)

	return topic, nil
}

// FetchRedditTIL fetches the top post from r/todayilearned
func (s *TopicService) FetchRedditTIL(ctx context.Context) (*Topic, error) {
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/todayilearned/top.json?limit=1&t=day", nil)
	if err != nil {
		return nil, fmt.Errorf("create Reddit TIL request: %w", err)
	}
	req.Header.Set("User-Agent", "GoChat/1.0 (by /u/melkeydev)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch Reddit TIL: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("this is response: ", resp)

	var redditResp struct {
		Data struct {
			Children []struct {
				Data struct {
					Title     string `json:"title"`
					URL       string `json:"url"`
					Score     int    `json:"score"`
					Permalink string `json:"permalink"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&redditResp); err != nil {
		return nil, fmt.Errorf("decode Reddit TIL response: %w", err)
	}

	if len(redditResp.Data.Children) == 0 {
		return nil, fmt.Errorf("no Reddit TIL posts found")
	}

	post := redditResp.Data.Children[0].Data
	fmt.Printf("Reddit TIL API Response: %+v\n", post)

	// Use Reddit URL for discussion
	redditURL := "https://reddit.com" + post.Permalink

	topic := &Topic{
		Title:       cleanText(post.Title),
		Description: fmt.Sprintf("Today's top learning with %d upvotes", post.Score),
		URL:         redditURL,
		Source:      "Reddit TIL",
	}

	fmt.Printf("Reddit TIL Topic: %+v\n", topic)

	return topic, nil
}

// FetchAllTopics fetches topics from all sources
func (s *TopicService) FetchAllTopics(ctx context.Context) ([]Topic, error) {
	topics := make([]Topic, 0, 3)

	// Fetch HackerNews
	hnTopic, err := s.FetchHackerNewsTop(ctx)
	if err != nil {
		fmt.Printf("Error fetching HackerNews topic: %v\n", err)
		// Use fallback topic
		topics = append(topics, Topic{
			Title:       "Tech News Discussion",
			Description: "Discuss today's technology news",
			URL:         "https://news.ycombinator.com",
			Source:      "HackerNews",
		})
	} else {
		topics = append(topics, *hnTopic)
	}

	// Fetch Reddit World News
	worldTopic, err := s.FetchRedditWorldNews(ctx)
	if err != nil {
		fmt.Printf("Error fetching Reddit WorldNews topic: %v\n", err)
		// Use fallback topic
		topics = append(topics, Topic{
			Title:       "World News Discussion",
			Description: "Discuss today's global news",
			URL:         "https://reddit.com/r/worldnews",
			Source:      "Reddit WorldNews",
		})
	} else {
		topics = append(topics, *worldTopic)
	}

	// Fetch Reddit TIL
	tilTopic, err := s.FetchRedditTIL(ctx)
	if err != nil {
		fmt.Printf("Error fetching Reddit TIL topic: %v\n", err)
		// Use fallback topic
		topics = append(topics, Topic{
			Title:       "Today I Learned",
			Description: "Share interesting facts",
			URL:         "https://reddit.com/r/todayilearned",
			Source:      "Reddit TIL",
		})
	} else {
		topics = append(topics, *tilTopic)
	}

	return topics, nil
}

