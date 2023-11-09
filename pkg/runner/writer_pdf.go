package runner

import (
	_ "embed"

	"fmt"
	"github.com/flopp/go-findfont"
	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/signintech/gopdf"
	"io"
	"log"
	"strconv"
)

type pdfResultWriter struct {
	writer io.Writer
}

// NewPDFResultWriter creates a new PDFResultWriter
func NewPDFResultWriter(writer io.Writer) ReportResultWriter {
	return &pdfResultWriter{writer: writer}
}

// Output writes the PDF base report to target writer
func (w *pdfResultWriter) Output(result []ReportResult) (err error) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	fmt.Println(findfont.List()[len(findfont.List())-1])
	fontPath, err := findfont.Find("DejaVuSerif.ttf")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found 'ttf' in '%s'\n", fontPath)
	err = pdf.AddTTFFont("wts11", fontPath)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = pdf.SetFont("wts11", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	pdf.AddHeader(func() {

	})
	pdf.AddFooter(func() {
		const X_bias float64 = 101
		pdf.SetXY(X_bias, 825)

		pdf.Text("Generated by github.com/LinuxSuRen/api-testing")
		pdf.AddExternalLink("https://github.com/LinuxSuRen/api-testing", 95+X_bias, 813, 200, 15)

		// pdf.Image("../pkg/runner/data/imgs/gopher.jpg", 500, 780, nil) //print image
	})

	const line_bias float64 = 50
	const Y_start float64 = 100
	for _, api := range result {
		pdf.AddPage()
		pdf.SetXY(50, Y_start)
		pdf.Cell(nil, "API:    "+string(api.API))
		pdf.SetXY(50, Y_start+line_bias)
		pdf.Cell(nil, "Count:  "+strconv.Itoa(api.Count))
		pdf.SetXY(50, Y_start+line_bias*2)
		pdf.Cell(nil, "Average:"+api.Average.String())
		pdf.SetXY(50, Y_start+line_bias*3)
		pdf.Cell(nil, "Max:    "+api.Max.String())
		pdf.SetXY(50, Y_start+line_bias*4)
		pdf.Cell(nil, "Min:    "+api.Min.String())
		pdf.SetXY(50, Y_start+line_bias*5)
		pdf.Cell(nil, "QPS:    "+strconv.Itoa(api.QPS))
		pdf.SetXY(50, Y_start+line_bias*6)
		pdf.Cell(nil, "Error:  "+strconv.Itoa(api.Error))
		pdf.SetXY(50, Y_start+line_bias*7)
		pdf.Cell(nil, "LastErrorMessage:")
		pdf.SetXY(50, Y_start+line_bias*8)
		pdf.Cell(nil, api.LastErrorMessage)

		if api.Error != 0 {
			pdf.Image("../pkg/runner/data/imgs/warn.jpg", 30, Y_start+line_bias*6-5, nil)
		}

	}

	fmt.Fprint(w.writer, "Report is OK!")
	pdf.WritePdf("Report.pdf")
	return
}

// WithAPIConverage sets the api coverage
func (w *pdfResultWriter) WithAPIConverage(apiConverage apispec.APIConverage) ReportResultWriter {
	return w
}
