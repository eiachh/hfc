package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/eiachh/hfc/logger"
)

var scrapeLinkDict map[int64][]string

func ScrapeDataOf(barcode int64, desiredLinkInd int) string {
	errorNoDataTxt := "Could not get any extra data, and won't be able to get more."
	barcodeStr := strconv.FormatInt(barcode, 15)

	if scrapeLinkDict == nil {
		scrapeLinkDict = make(map[int64][]string)
	}

	if desiredLinkInd < 0 {
		return errorNoDataTxt
	}

	if links, haskey := scrapeLinkDict[barcode]; haskey {
		if desiredLinkInd >= len(links) {
			return errorNoDataTxt
		}
		scrapedText, scrapeError := GetTextFomrWebsite(links[desiredLinkInd])
		for i := desiredLinkInd + 1; scrapeError != nil; i++ {
			if i >= len(links) {
				return errorNoDataTxt
			}
			scrapedText, scrapeError = GetTextFomrWebsite(links[i])
		}
		return scrapedText
	}

	scrapedLinks, duckduckError := DuckDuckGoSearch(barcodeStr)
	if duckduckError != nil {
		return errorNoDataTxt
	}
	scrapeLinkDict[barcode] = scrapedLinks

	scrapedText, scrapeError := GetTextFomrWebsite(scrapedLinks[0])
	for i := 0 + 1; scrapeError != nil; i++ {
		if i >= len(scrapedLinks) {
			return errorNoDataTxt
		}
		scrapedText, scrapeError = GetTextFomrWebsite(scrapedLinks[i])
	}

	return scrapedText
}

func GetTextFomrWebsite(url string) (string, error) {
	logger.Log().Debugf("Scraping website: %s", url)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx, cancel = chromedp.NewExecAllocator(ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("start-maximized", false),
		)...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var bodyAsStr string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`//body`),
		chromedp.Text(`//body`, &bodyAsStr),
	}

	err := chromedp.Run(ctx, tasks)
	if err != nil {
		logger.Log().Warnf("timed out waiting for website: %s", url)
		return "", errors.New("timed out waiting for website")
	}

	return bodyAsStr, nil
}

func DuckDuckGoSearch(searchVerb string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx, cancel = chromedp.NewExecAllocator(ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("start-maximized", false),
		)...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	duckduckgoUrl := "https://duckduckgo.com/"
	searchbarSelector := `//input[@name="q"]`

	linksSelector := `//h2/child::a`
	var projects []*cdp.Node

	tasks := chromedp.Tasks{
		chromedp.Navigate(duckduckgoUrl),
		chromedp.WaitVisible(searchbarSelector),
		chromedp.SendKeys(searchbarSelector, searchVerb),
		chromedp.Submit(searchbarSelector),
		chromedp.WaitVisible(`//a[contains(text(), 'Images')]`),
		chromedp.Nodes(linksSelector, &projects),
	}

	var links []string
	err := chromedp.Run(ctx, tasks)
	if err != nil {
		logger.Log().Warnf("timed out waiting for duckduckgo with searchverb: %s", searchVerb)
		return links, errors.New("timed out waiting for duckduckgo")
	}

	for i := 0; i < len(projects); i++ {
		links = append(links, projects[i].AttributeValue("href"))
	}

	return links, nil
}
