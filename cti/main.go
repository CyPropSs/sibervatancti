package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run main.go <URL>")
		return
	}
	rawURL := os.Args[1]

	safeFolderName := strings.ReplaceAll(rawURL, "https://", "")
	safeFolderName = strings.ReplaceAll(safeFolderName, "http://", "")

	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	safeFolderName = reg.ReplaceAllString(safeFolderName, "_")

	if len(safeFolderName) > 100 {
		safeFolderName = safeFolderName[:100]
	}

	os.MkdirAll(safeFolderName, 0755)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var htmlContent string
	var screenshot []byte
	var links []string

	fmt.Printf(">> %s adresi taranıyor... \n", rawURL)

	err := chromedp.Run(ctx,
		chromedp.Navigate(rawURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
		chromedp.FullScreenshot(&screenshot, 100),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &links),
	)

	if err != nil {
		fmt.Printf("\n!!! HATA !!!\nSiteye ulaşılamadı (404 veya Bağlantı Hatası): %v \n", err)
		return
	}

	os.WriteFile(filepath.Join(safeFolderName, "icerik.html"), []byte(htmlContent), 0644)
	os.WriteFile(filepath.Join(safeFolderName, "screenshot.png"), screenshot, 0644)

	linkData := strings.Join(links, "\n")
	os.WriteFile(filepath.Join(safeFolderName, "linkler.txt"), []byte(linkData), 0644)

	fmt.Printf("✓ Başarılı! '%s' klasörü oluşturuldu. [cite: 17, 18]\n", safeFolderName)
}
