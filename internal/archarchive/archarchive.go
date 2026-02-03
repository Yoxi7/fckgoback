package archarchive

import (
	"fmt"
	"strings"
)

type ArchArchive struct {
	Year  string
	Month string
	Day   string
}

func NewArchArchive() *ArchArchive {
	return &ArchArchive{}
}

func (a *ArchArchive) GetLink() string {
	url := buildURL(a.Year, a.Month, a.Day)
	if url != "" && !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}

func (a *ArchArchive) SelectDate() error {
	years, err := ParseYears()
	if err != nil {
		return fmt.Errorf("failed to parse years: %w", err)
	}
	if len(years) == 0 {
		return fmt.Errorf("no years available")
	}

	year, err := Ask(years, "Select year:")
	if err != nil {
		return fmt.Errorf("failed to select year: %w", err)
	}
	a.Year = year

	months, err := ParseMonths(a.Year)
	if err != nil {
		return fmt.Errorf("failed to parse months: %w", err)
	}
	if len(months) == 0 {
		return fmt.Errorf("no months available for year %s", a.Year)
	}

	month, err := Ask(months, "Select month:")
	if err != nil {
		return fmt.Errorf("failed to select month: %w", err)
	}
	a.Month = month

	days, err := ParseDays(a.Year, a.Month)
	if err != nil {
		return fmt.Errorf("failed to parse days: %w", err)
	}
	if len(days) == 0 {
		return fmt.Errorf("no days available for %s/%s", a.Year, a.Month)
	}

	day, err := Ask(days, "Select day:")
	if err != nil {
		return fmt.Errorf("failed to select day: %w", err)
	}
	a.Day = day

	return nil
}

func (a *ArchArchive) MenuRun() (string, error) {
	if err := a.SelectDate(); err != nil {
		return "", err
	}
	return a.GetLink(), nil
}



