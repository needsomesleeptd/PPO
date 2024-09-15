package document_view

import (
	document_handler "annotater/internal/http-server/handlers/document"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	bboxes_utils "annotater/tech_ui/utils/bboxes"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
)

var (
	fileFormat = ".png"
)

func GetCheckDocumentResult(resp *document_handler.ResponseCheckDoucment, folderName string) (string, error) {
	if resp.Response.Status != response.StatusOK {
		return "", errors.New(resp.Response.Error)
	}
	hashMarkUpType := make(map[uint64]models.MarkupType)

	for _, markUpType := range resp.MarkupTypes {
		hashMarkUpType[markUpType.ID] = markUpType
	}

	markups := resp.Markups
	fmt.Print(len(markups))
	err := os.MkdirAll(folderName, 0777)
	if err != nil {
		return "", err
	}
	for i, markup := range markups {

		img, _, err := image.Decode(bytes.NewReader(markup.PageData))
		if err != nil {
			return "", err
		}
		boundingBoxImg := image.NewRGBA(img.Bounds())
		draw.Draw(boundingBoxImg, img.Bounds(), img, image.Point{}, draw.Src)
		boundingBoxColor := color.RGBA{255, 0, 0, 255}
		x1, y1, x2, y2 := int(markup.ErrorBB[0]), int(markup.ErrorBB[1]), int(markup.ErrorBB[2]), int(markup.ErrorBB[3])
		boundingBoxes := []bboxes_utils.BoundingBox{
			{
				XMin: x1,
				YMin: y1,
				XMax: x2,
				YMax: y2,
			},
		}
		bboxes_utils.DrawBoundingBoxes(boundingBoxImg, boundingBoxes, boundingBoxColor)
		fmt.Print(hashMarkUpType, hashMarkUpType[markup.ID].ClassName)
		bboxes_utils.DrawText(boundingBoxImg, x1, y2*2, hashMarkUpType[markup.ID].Description)
		outputFile, err := os.Create(folderName + "/" + strconv.Itoa(i) + fileFormat)
		if err != nil {
			return "", err
		}
		defer outputFile.Close()

		err = png.Encode(outputFile, boundingBoxImg)
		if err != nil {
			return "", err
		}
	}
	return "OK", nil
}

func GetLoadDocumentResult(resp *response.Response) (string, error) {
	if resp.Status == response.StatusOK {
		return resp.Status, nil
	} else {
		return "", errors.New(resp.Error)
	}
}
