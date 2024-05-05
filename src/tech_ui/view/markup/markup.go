package markup_view

import (
	annot_handler "annotater/internal/http-server/handlers/annot"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
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

func DrawBbsOnMarkups(resp *annot_handler.ResponseGetAnnots, folderName string) error {
	if resp.Response != response.OK() {
		return errors.New(resp.Response.Error)
	}
	hashMarkUpType := make(map[uint64]models.MarkupType)

	markups := resp.Markups
	os.Mkdir(folderName, 0777)
	for _, markup := range markups {

		img, _, err := image.Decode(bytes.NewReader(markup.PageData))
		if err != nil {
			return err
		}
		boundingBoxImg := image.NewRGBA(img.Bounds())
		draw.Draw(boundingBoxImg, img.Bounds(), img, image.Point{}, draw.Src)
		imgWidth := float32(img.Bounds().Dx())
		imgHeight := float32(img.Bounds().Dy())
		boundingBoxColor := color.RGBA{255, 0, 0, 255}
		x1, y1, x2, y2 := int(markup.ErrorBB[0]*imgWidth), int(markup.ErrorBB[1]*imgHeight), int(markup.ErrorBB[2]*imgWidth), int(markup.ErrorBB[3]*imgHeight)
		boundingBoxes := []bboxes_utils.BoundingBox{
			{
				XMin: x1,
				YMin: y1,
				XMax: x2,
				YMax: y2,
			},
		}
		bboxes_utils.DrawBoundingBoxes(boundingBoxImg, boundingBoxes, boundingBoxColor)
		bboxes_utils.DrawText(boundingBoxImg, x1, y1, hashMarkUpType[markup.ID].ClassName)
		outputFile, err := os.Create(folderName + "/" + strconv.Itoa(int(markup.ID)) + fileFormat)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		err = png.Encode(outputFile, boundingBoxImg)
		if err != nil {
			return err
		}
	}
	return nil
}
