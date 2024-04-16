package document_view

import (
	document_handler "annotater/internal/http-server/handlers/document"
	response "annotater/internal/lib/api"
	bboxes_utils "annotater/tech_ui/utils/bboxes"
	"bytes"
	"errors"
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
	markups := resp.Markups
	os.Mkdir(folderName, 0777)
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
	return "Success", nil
}
