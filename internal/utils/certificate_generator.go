package utils

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
)

func GenerateCertificate(username, courseTitle, instructorName, completionDate string) (image.Image, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(courseTitle) == "" ||
		strings.TrimSpace(instructorName) == "" || strings.TrimSpace(completionDate) == "" {
		return nil, errors.New("all of username, courseTitle, instructorName, completionDate must be provided")
	}

	width, height := 1200, 800
	dc := gg.NewContext(width, height)

	dc.SetColor(color.White)
	dc.Clear()

	dc.SetLineWidth(8)
	dc.SetColor(color.Black)
	dc.DrawRectangle(20, 20, float64(width-40), float64(height-40))
	dc.Stroke()

	dc.SetColor(color.Black)
	dc.DrawStringAnchored("CERTIFICATE OF COMPLETION", float64(width)/2, 150, 0.5, 0.5)
	dc.DrawStringAnchored("This is to certify that", float64(width)/2, 250, 0.5, 0.5)
	dc.DrawStringAnchored(username, float64(width)/2, 320, 0.5, 0.5)
	dc.DrawStringAnchored("has successfully completed the course", float64(width)/2, 400, 0.5, 0.5)
	dc.DrawStringAnchored(courseTitle, float64(width)/2, 460, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("Instructor: %s", instructorName), float64(width)*0.25, 600, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("Date: %s", completionDate), float64(width)*0.75, 600, 0.5, 0.5)

	return dc.Image(), nil
}
