package service

import "errors"

// TODO make this into an actual webscraper
func ScrapeDataOf(barcode int, link int) (string, error) {
	if link == 1 {
		return "<html>     categories_hierarchy: [en:snacks,en:sweet-nacks,en:confecioneries,en:candie,en:gummi-candies],</html>", nil
	} else {
		return "nothing was found", errors.New("webscraping failed no new data")
	}
}
