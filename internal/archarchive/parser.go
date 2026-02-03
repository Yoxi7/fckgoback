package archarchive

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const Endpoint = "https://archive.archlinux.org/repos"

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func getElements(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var elements []string
	doc.Find("pre a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasSuffix(href, "/") && href != "../" {
			element := strings.TrimSuffix(href, "/")
			if element != ".." {
				elements = append(elements, element)
			}
		}
	})

	return elements, nil
}

func buildURL(year, month, day string) string {
	if day != "" && month == "" {
		return ""
	}
	if month != "" && day != "" {
		return fmt.Sprintf("%s/%s/%s/%s", Endpoint, year, month, day)
	}
	if month != "" {
		return fmt.Sprintf("%s/%s/%s", Endpoint, year, month)
	}
	return fmt.Sprintf("%s/%s", Endpoint, year)
}

func reverseStrings(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func ParseYears() ([]string, error) {
	resp, err := httpClient.Get(Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch years: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	elements, err := getElements(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	if len(elements) > 3 {
		elements = elements[:len(elements)-3]
	}

	reverseStrings(elements)

	return elements, nil
}

func ParseMonths(year string) ([]string, error) {
	url := buildURL(year, "", "")
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch months: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	elements, err := getElements(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	reverseStrings(elements)
	return elements, nil
}

func ParseDays(year, month string) ([]string, error) {
	url := buildURL(year, month, "")
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch days: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	elements, err := getElements(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	reverseStrings(elements)
	return elements, nil
}

func CheckArchiveAvailability(url string) error {
	requiredRepos := []string{"core", "extra"}
	optionalRepos := []string{"multilib"}
	
	var missingRepos []string
	
	for _, repo := range requiredRepos {
		testURL := url + "/" + repo + "/os/x86_64/" + repo + ".db"
		resp, err := httpClient.Get(testURL)
		if err != nil {
			return fmt.Errorf("failed to check archive availability for %s: %w", repo, err)
		}
		resp.Body.Close()
		
		if resp.StatusCode == http.StatusNotFound {
			missingRepos = append(missingRepos, repo)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("archive unavailable for %s at %s (status: %d)", repo, url, resp.StatusCode)
		}
	}
	
	if len(missingRepos) > 0 {
		return fmt.Errorf("archive incomplete at %s - missing required repositories: %v. This date may not be fully available in the archive", url, missingRepos)
	}
	
	for _, repo := range optionalRepos {
		testURL := url + "/" + repo + "/os/x86_64/" + repo + ".db"
		resp, err := httpClient.Head(testURL)
		if err == nil {
			resp.Body.Close()
		}
	}
	
	return nil
}

