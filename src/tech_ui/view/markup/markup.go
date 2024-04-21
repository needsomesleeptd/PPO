package markup_view

import (
	annot_handler "annotater/internal/http-server/handlers/annot"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	bboxes_utils "annotater/tech_ui/utils/bboxes"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"
)

var (
	fileFormat = ".png"
)

func GetCheckDocumentResult(resp *annot_handler.ResponseGetByUserID, folderName string) error {
	if resp.Response != response.OK() {
		return errors.New(resp.Response.Error)
	}
	hashMarkUpType := make(map[uint64]models.MarkupType)

	markups := resp.Markups
	os.Mkdir(folderName, 0777)
	for i, markup := range markups {

		img, _, err := image.Decode(bytes.NewReader(markup.PageData))
		if err != nil {
			return err
		}
		boundingBoxImg := image.NewRGBA(img.Bounds())
		draw.Draw(boundingBoxImg, img.Bounds(), img, image.Point{}, draw.Src)
		//boundingBoxColor := color.RGBA{255, 0, 0, 255}
		x1, y1, _, _ := int(markup.ErrorBB[0]), int(markup.ErrorBB[1]), int(markup.ErrorBB[2]), int(markup.ErrorBB[3])
		/*boundingBoxes := []bboxes_utils.BoundingBox{
			{
				XMin: x1,
				YMin: y1,
				XMax: x2,
				YMax: y2,
			},
		}*/
		//bboxes_utils.DrawBoundingBoxes(boundingBoxImg, boundingBoxes, boundingBoxColor)
		fmt.Print(hashMarkUpType, hashMarkUpType[markup.ID].ClassName)
		bboxes_utils.DrawText(boundingBoxImg, x1, y1, hashMarkUpType[markup.ID].ClassName)
		outputFile, err := os.Create(folderName + "/" + strconv.Itoa(i) + fileFormat)
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
