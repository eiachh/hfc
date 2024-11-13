package service

import "errors"

func ScrapeDataOf(barcode int64, link int) (string, error) {
	if link == 1 {
		return "<html> product_name: WEBSCRAPE,</html>", nil
	} else {
		return "nothing was found", errors.New("webscraping failed no new data")
	}
}
