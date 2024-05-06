package report_creator

import (
	"annotater/internal/models"
	bboxes_utils "annotater/tech_ui/utils/bboxes"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"os/exec"
	"strconv"

	"github.com/google/uuid"
)

var (
	texFileFilename = "check.tex"
	pdfFileFilename = "check.pdf"
	imgFolderPath   = "images/"
	fileFormat      = ".png"
	texHeader       = "\\documentclass{article}\n\\usepackage{graphicx}\n\n\\begin{document}\n\\section{Error Report}\n"
	texTail         = "\\end{document}"
)

type IReportCreator interface {
	CreateReport(reportID uuid.UUID, markups []models.Markup, markupTypes []models.MarkupType) (*models.ErrorReport, error)
}

type PDFReportCreator struct {
	folderPath string
}

func NewPDFReportCreator(workFolderPath string) IReportCreator { // TODO:: think about this
	err := os.MkdirAll(workFolderPath, 0777)
	if err != nil {
		panic(err)
	}
	return &PDFReportCreator{
		folderPath: workFolderPath,
	}
}

func (cr *PDFReportCreator) addImageLatex(imgPath string) string {
	return "\\newpage\n\\noindent\\includegraphics[width=0.9\\textwidth, height=0.9\\textheight]{" + imgPath + "}\n\n"
}

func (cr *PDFReportCreator) saveImagesWithBBs(filePathSave string, markups []models.Markup) ([]string, error) {
	imgPaths := make([]string, len(markups))
	for i, markup := range markups {
		img, _, err := image.Decode(bytes.NewReader(markup.PageData))
		if err != nil {
			return nil, err
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

		imgFilePath := filePathSave + strconv.Itoa(i) + fileFormat

		outputFile, err := os.Create(imgFilePath)
		if err != nil {
			return nil, err
		}

		imgPaths[i] = imgFilePath
		defer outputFile.Close()
		err = png.Encode(outputFile, boundingBoxImg)
		if err != nil {
			return nil, err
		}
	}
	return imgPaths, nil
}

func (cr *PDFReportCreator) CreateReport(reportID uuid.UUID, markups []models.Markup, markupTypes []models.MarkupType) (*models.ErrorReport, error) {
	senderFolderPath := cr.folderPath + "/" + reportID.String() + "/"

	hashMarkUpType := make(map[uint64]models.MarkupType)

	for _, markUpType := range markupTypes {
		hashMarkUpType[markUpType.ID] = markUpType
	}
	err := os.Mkdir(senderFolderPath, 0777)
	if err != nil {
		return nil, err
	}
	texFilePath := senderFolderPath + texFileFilename

	file, err := os.Create(texFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//creating folder for images
	err = os.MkdirAll(senderFolderPath+imgFolderPath, 0777)
	if err != nil {
		return nil, err
	}

	imgPaths, err := cr.saveImagesWithBBs(senderFolderPath+imgFolderPath, markups)
	if err != nil {
		return nil, err
	}
	// Start writing LaTeX content
	content := texHeader

	// Iterate over images and texts to insert each pair on a separate page
	for i := 0; i < len(markups); i++ {
		imgLatex := cr.addImageLatex(imgPaths[i])
		var description string
		if markupType, exists := hashMarkUpType[markups[i].ClassLabel]; exists {
			description = markupType.Description
		} else {
			description = fmt.Sprintf("error not found description for label: %v", markups[i].ClassLabel)
		}
		content += imgLatex + description + "\n"
	}

	// End writing LaTeX content
	content += texTail

	// Write the content to the LaTeX file
	_, err = file.WriteString(content)
	if err != nil {
		return nil, err
	}

	outputDirKey := fmt.Sprintf("-output-directory=%s", senderFolderPath)

	cmd := exec.Command("pdflatex", outputDirKey, texFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run() //run twice for latex reasons
	if err != nil {
		return nil, fmt.Errorf("error running first latex compile: %v", err)
	}

	pdfFilePath := senderFolderPath + pdfFileFilename
	pdfBytes, err := os.ReadFile(pdfFilePath)
	if err != nil {
		return nil, err
	}

	pdfByteSlice := []byte(pdfBytes)

	report := models.ErrorReport{
		DocumentID: reportID,
		ReportData: pdfByteSlice,
	}
	return &report, nil
}
