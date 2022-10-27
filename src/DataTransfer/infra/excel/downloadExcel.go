package excel

import (
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

func downloadExcel(spreadsheetID string, filename string) {
	ctx := context.Background()

	driveService, err := drive.NewService(ctx, option.WithCredentialsFile("secret.json"))
	if err != nil {
		log.Fatal("service error: ", err)
	}

	response, err := driveService.Files.Export(spreadsheetID, mimeType).Download()

	if err != nil {
		log.Fatal("http response error: ", err)
	}

	errDownloading := download(response, filename)

	if errDownloading != nil {
		log.Fatal("error downloading: ", err)
	}

}

func download(resp *http.Response, filename string) error {

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, b, 0644) //0644 meaning: User: read & write, Group: read, Other: read
	if err != nil {
		return err
	}

	return nil
}
