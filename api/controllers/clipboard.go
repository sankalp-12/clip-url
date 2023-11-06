package controllers


import (
	"context"
	"log"
	"os"
    "fmt"
	"github.com/asaskevich/govalidator"
	"github.com/sankalp-12/clip-url/utils"
	"golang.design/x/clipboard"
)

func SetupWatcher(ctx context.Context) {
	ch := clipboard.Watch(ctx, clipboard.FmtText)
	lastWrite := ""

	for data := range ch {
		if lastWrite == string(data) {
			continue
		}
		if govalidator.IsURL(string(data)) {
			if flag, url := utils.Put(ctx, string(data)); flag {
				url = os.Getenv("DOMAIN") + os.Getenv("API") + url
                fmt.Println(url)
				lastWrite = url
				clipboard.Write(clipboard.FmtText, []byte(url))
			}
		} else {
			log.Printf("Bad request: Clipboard data[%s] is not a URL", string(data))
		}
	}
}
