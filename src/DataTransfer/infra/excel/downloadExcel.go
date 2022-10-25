package excel

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func downloadExcel(spreadsheetID string, filename string) {
	ctx := context.Background()

	driveService, err := drive.NewService(ctx, option.WithCredentialsFile("secret.json"))
	if err != nil {
		log.Fatal(err, " client")
	}

	mimeType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	response, err := driveService.Files.Export(spreadsheetID, mimeType).Download()

	if err != nil {
		log.Fatal(err, "sheet")
	}

	errDownloading := download(response, filename)

	if errDownloading != nil {
		log.Fatal(err)
	}

}

func download(resp *http.Response, filename string) error {

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		checkError(err)
		return err
	}

	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		checkError(err)
		return err
	}

	fmt.Println("Doc downloaded in ", filename)

	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
