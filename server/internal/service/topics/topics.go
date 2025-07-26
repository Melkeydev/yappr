package topics

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type TopicService struct {
	client       *http.Client
	redditToken  string
	tokenExpiry  time.Time
}

type Topic struct {
	Title       string
	Description string
	URL         string
	Source      string
}

type redditTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewTopicService() *TopicService {
	return &TopicService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func cleanText(text string) string {
	// Decode HTML entities like &amp;, &lt;, &gt;, &quot;, etc.
	decoded := html.UnescapeString(text)
	// Trim whitespace
	return strings.TrimSpace(decoded)
}

const redditUA = "desktop:gochat:1.0 (by /u/melkeydev)" // Reddit's preferred UA format

// getRedditToken fetches a new OAuth token from Reddit
func (s *TopicService) getRedditToken(ctx context.Context) error {
	// Get Reddit credentials from environment
	clientID := os.Getenv("REDDIT_CLIENT_ID")
	clientSecret := os.Getenv("REDDIT_CLIENT_SECRET")
	
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("REDDIT_CLIENT_ID and REDDIT_CLIENT_SECRET must be set")
	}
	
	// Create token request
	data := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequestWithContext(ctx, "POST", "https://www.reddit.com/api/v1/access_token", data)
	if err != nil {
		return fmt.Errorf("create token request: %w", err)
	}
	
	// Set headers
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("User-Agent", redditUA)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// Make request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch token: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("reddit token request failed with %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var tokenResp redditTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("decode token response: %w", err)
	}
	
	// Store token and expiry
	s.redditToken = tokenResp.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // Subtract 60s for safety
	
	fmt.Printf("Reddit OAuth token obtained, expires at %v\n", s.tokenExpiry)
	return nil
}

func (s *TopicService) getRedditJSON(ctx context.Context, url string, out any) error {
	// Check if we need a new token
	if s.redditToken == "" || time.Now().After(s.tokenExpiry) {
		if err := s.getRedditToken(ctx); err != nil {
			return fmt.Errorf("get reddit token: %w", err)
		}
	}
	
	// Use OAuth endpoint instead of public API
	// Convert public URL to OAuth URL
	oauthURL := strings.Replace(url, "https://www.reddit.com/", "https://oauth.reddit.com/", 1)
	
	req, err := http.NewRequestWithContext(ctx, "GET", oauthURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", redditUA)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.redditToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("reddit %s -> %d: %s", oauthURL, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (s *TopicService) FetchHackerNewsTop(ctx context.Context) (*Topic, error) {
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

	if story.URL == "" {
		story.URL = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ID)
	}

	topic := &Topic{
		Title:       cleanText(story.Title),
		Description: fmt.Sprintf("Top HN story with %d points by %s", story.Score, story.By),
		URL:         story.URL,
		Source:      "HackerNews",
	}

	return topic, nil
}

func (s *TopicService) FetchRedditWorldNews(ctx context.Context) (*Topic, error) {
	url := "https://www.reddit.com/r/worldnews/top.json?limit=1&t=day&raw_json=1"

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

	if err := s.getRedditJSON(ctx, url, &redditResp); err != nil {
		return nil, fmt.Errorf("fetch Reddit worldnews: %w", err)
	}
	if len(redditResp.Data.Children) == 0 {
		return nil, fmt.Errorf("reddit worldnews: empty children")
	}

	post := redditResp.Data.Children[0].Data
	return &Topic{
		Title:       cleanText(post.Title),
		Description: fmt.Sprintf("Top world news with %d upvotes", post.Score),
		URL:         "https://reddit.com" + post.Permalink,
		Source:      "Reddit WorldNews",
	}, nil
}

func (s *TopicService) FetchRedditTIL(ctx context.Context) (*Topic, error) {
	url := "https://www.reddit.com/r/todayilearned/top.json?limit=1&t=day&raw_json=1"

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

	if err := s.getRedditJSON(ctx, url, &redditResp); err != nil {
		return nil, fmt.Errorf("fetch Reddit TIL: %w", err)
	}
	if len(redditResp.Data.Children) == 0 {
		return nil, fmt.Errorf("reddit TIL: empty children")
	}

	post := redditResp.Data.Children[0].Data
	return &Topic{
		Title:       cleanText(post.Title),
		Description: fmt.Sprintf("Today's top learning with %d upvotes", post.Score),
		URL:         "https://reddit.com" + post.Permalink,
		Source:      "Reddit TIL",
	}, nil
}

func (s *TopicService) FetchAllTopics(ctx context.Context) ([]Topic, error) {
	topics := make([]Topic, 0, 3)

	hnTopic, err := s.FetchHackerNewsTop(ctx)
	if err != nil {
		fmt.Printf("Error fetching HackerNews topic: %v\n", err)
		topics = append(topics, Topic{
			Title:       "Tech News Discussion",
			Description: "Discuss today's technology news",
			URL:         "https://news.ycombinator.com",
			Source:      "HackerNews",
		})
	} else {
		topics = append(topics, *hnTopic)
	}

	worldTopic, err := s.FetchRedditWorldNews(ctx)
	if err != nil {
		fmt.Printf("Error fetching Reddit WorldNews topic: %v\n", err)
		topics = append(topics, Topic{
			Title:       "World News Discussion",
			Description: "Discuss today's global news",
			URL:         "https://reddit.com/r/worldnews",
			Source:      "Reddit WorldNews",
		})
	} else {
		topics = append(topics, *worldTopic)
	}

	tilTopic, err := s.FetchRedditTIL(ctx)
	if err != nil {
		fmt.Printf("Error fetching Reddit TIL topic: %v\n", err)
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

