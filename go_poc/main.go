// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	email := `iriversland@gmail.com`
	psw := `pesQyr-zekta0-havwip`

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
		chromedp.WithDebugf(log.Printf),
		chromedp.WithErrorf(log.Fatalf),
	)
	defer cancel()
	// start the browser
	// see https://github.com/chromedp/chromedp/issues/513#issuecomment-558122963
	if err := chromedp.Run(ctx); err != nil {
		log.Println("start:", err)
	}

	// create a timeout
	// ctx, cancel = context.WithTimeout(ctx, 600*time.Second)
	// defer cancel()

	// navigate to a page, wait for an element, click
	log.Printf(`Start scraping...`)
	var title string
	var body string
	var text string
	var inputCode string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.linkedin.com/login`),
		chromedp.SendKeys(`#username`, email),
		chromedp.SendKeys(`#password`, psw),
		chromedp.Click(`button[type=submit][aria-label="Sign in"]`),
		chromedp.Sleep(2*time.Second),
		chromedp.Title(&title),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(`Title: %s\n`, title)
	if strings.Contains(title, "Login") {
		log.Fatal(`Stuck at login page`)
	}
	if strings.Contains(title, "Security") {
		err = chromedp.Run(ctx,
			chromedp.OuterHTML(`body`, &body),
			chromedp.Title(&title),
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Println(text)
				log.Println(body)
				log.Println(title)

				var code string
				fmt.Printf("Enter code:")
				fmt.Scanln(&code)
				fmt.Printf("You entered code `%s`, len=%d\n", code, len(code))
				return chromedp.SendKeys(`#input__email_verification_pin`, code, chromedp.ByID).Do(ctx)
			}),
			chromedp.Value(`#input__email_verification_pin`, &inputCode, chromedp.ByID),
			chromedp.Click(`#email-pin-submit-button`, chromedp.NodeVisible, chromedp.ByID),
			chromedp.Sleep(2*time.Second),
			chromedp.Title(&title),
			// wait for footer element is visible (ie, page is loaded)
			// chromedp.WaitVisible(`body > footer`),
			// find and click "Expand All" link
			// chromedp.Click(`#pkg-examples > a`, chromedp.NodeVisible),
			// retrieve the value of the textarea
			// chromedp.Value(`#example-After div.Documentation-exampleDetailsBody textarea.code`, &example),
		)
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(title, "Security") {
			log.Fatal(`Cannot complete MFA, page remains at MFA page`)
		}
	}

	// locate search bar
	// <input class="search-global-typeahead__input always-show-placeholder" placeholder="Search" role="combobox" aria-autocomplete="list" aria-label="Search" aria-activedescendant="" aria-expanded="false" aria-owns="" type="text">
	body = ``
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// chromedp.OuterHTML(`body`, &body).Do(ctx)
			// log.Printf("Body:%s\n", body)

			log.Printf(`Page is loading...`)
			chromedp.WaitNotVisible(`div.initial-load-animation`, chromedp.ByQuery).Do(ctx)
			log.Printf(`Page loading done.`)

			// chromedp.Text(`body`, &body).Do(ctx)
			// log.Printf("Body:%s\n", body)

			// log.Printf(`Waiting main...`)
			// chromedp.WaitVisible(`#main`, chromedp.ByID).Do(ctx)
			// chromedp.OuterHTML(`#main`, &body, chromedp.ByID).Do(ctx)
			// log.Printf(`Main: %s`, body)
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// top search bar
			//   <input class="search-global-typeahead__input always-show-placeholder" placeholder="Search" role="combobox" aria-autocomplete="list" aria-label="Search" aria-activedescendant="" aria-expanded="false" aria-owns="" type="text">

			// log.Println(`Locating search bar...`)
			// chromedp.WaitVisible(`input[aria-label="Search"]`, chromedp.ByQuery).Do(ctx) // stucks here...
			// log.Println(`Sending keys...`)
			// chromedp.SendKeys(`input[aria-label="Search"]`, "bbc", chromedp.ByQuery).Do(ctx)
			// log.Println(`Keys sent.`)
			// chromedp.Sleep(2 * time.Second).Do(ctx)

			log.Println(`Print out nav search`)
			printHTML(`#global-nav-search`).Do(ctx)
			log.Println(body)
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			searchBarSel := `body > div.application-outlet > #global-nav > div.global-nav__content > div#global-nav-search > div.search-global-typeahead > div#global-nav-typeahead > input.search-global-typeahead__input`
			printDocumentTree(searchBarSel).Do(ctx)
			log.Println(`Focusing input...`) // <--- stuck here - cannot compute box model
			chromedp.Focus(fmt.Sprintf(`document.querySelector("%s")`, searchBarSel), chromedp.ByJSPath).Do(ctx)
			log.Println(`Sending keys...`)
			chromedp.SendKeys(fmt.Sprintf(`document.querySelector("%s")`, searchBarSel), "bbc", chromedp.ByJSPath).Do(ctx)
			printHTML(`#global-nav-search`).Do(ctx)
			return nil
		}),
		// chromedp.WaitVisible(`#main`, chromedp.ByID).Do(ctx)
		// chromedp.WaitVisible(`#global-nav-typeahead search-global-typeahead__hit`),
		// chromedp.OuterHTML(`#global-nav-typeahead search-global-typeahead__hit`, &body, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	// err = chromedp.Run(ctx,
	// 	researchCompany(),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Printf(`Done\n`)
}

func researchCompany() chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Printf(`Visiting company search page...`)
		chromedp.Navigate(`https://www.linkedin.com/search/results/companies/?keywords=bbc`).Do(ctx)
		printDocumentTree(`body`).Do(ctx)
		return nil
	})
}

func printDocumentTree(path string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var nodes []*cdp.Node
		log.Printf(`Search tree of %s:\n`, path)
		chromedp.Nodes(fmt.Sprintf(`document.querySelector("%s")`, path), &nodes, chromedp.ByJSPath).Do(ctx)
		fmt.Println(nodes[0].Dump("  ", "  ", false))
		return nil
	})
}

func printHTML(path string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var body string
		log.Printf("Printing HTML for %s\n", path)
		chromedp.OuterHTML(fmt.Sprintf(`document.querySelector("%s")`, path), &body, chromedp.ByJSPath).Do(ctx)
		fmt.Println(body)
		return nil
	})
}
