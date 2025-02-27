package models

type UrlData struct {
	OriginalUrl string
	ShortUrl    string
}

type OrigUrlData struct {
	OriginalUrl string `json:"original_url" validate:"required"`
}

type ShortUrlData struct {
	ShortUrl string `json:"shortened_url"`
}
