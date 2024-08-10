package utils

import "time"

func GetMoscowTime() time.Time {
	utcTime := time.Now().UTC()
	return utcTime.Add(3 * time.Hour)
}

func GetEuropeTime(location string) (time.Time, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(loc), nil
}

func GetAccessTokenTTL() time.Duration {
	return time.Hour
}

func GetRefreshTokenTTL() time.Duration {
	return 24 * time.Hour
}

func GetEmailTokenTTL() time.Duration {
	return 24 * time.Hour
}

func GetEmailUpdateTTL() time.Duration {
	return 24 * time.Hour
}

func GetConnectionsCountTTL() time.Duration {
	return time.Second
}
