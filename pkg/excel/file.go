package excel

import (
	"bytes"
	"context"
	"log"

	"github.com/xuri/excelize/v2"
)

type File interface {
	NewSheet(context.Context, string) Sheet
	Save(context.Context, string) error
	SaveAs(context.Context, string) error
	WriteToBuffer(context.Context) (*bytes.Buffer, error)
	MimeType(context.Context) string
}

type file struct {
	xls      *excelize.File
	mimeType string
}

func newFile(xls *excelize.File) File {
	return &file{
		xls:      xls,
		mimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.shee",
	}
}

func (f *file) NewSheet(ctx context.Context, name string) Sheet {
	numSheet, err := f.xls.NewSheet(name)
	if err != nil {
		log.Fatal(err)
	}
	return newSheet(name, numSheet, f.xls)
}

func (f *file) Save(ctx context.Context, path string) error {
	f.xls.Path = path
	return f.xls.Save()
}

func (f *file) SaveAs(ctx context.Context, path string) error {
	f.xls.Path = path
	return f.xls.SaveAs(f.xls.Path)
}

func (f *file) WriteToBuffer(ctx context.Context) (*bytes.Buffer, error) {
	return f.xls.WriteToBuffer()
}

func (f *file) MimeType(ctx context.Context) string {
	return f.mimeType
}
