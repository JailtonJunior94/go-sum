package excel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/xuri/excelize/v2"
)

type Sheet interface {
	Write(context.Context, interface{}) error
}

type sheet struct {
	name     string
	numSheet int

	line      int
	hasHeader bool

	xls *excelize.File
}

func newSheet(name string, numSheet int, xls *excelize.File) Sheet {
	return &sheet{
		name:      name,
		numSheet:  numSheet,
		line:      1,
		hasHeader: false,
		xls:       xls,
	}
}

func (s *sheet) Write(ctx context.Context, data interface{}) error {
	if reflect.ValueOf(data).Kind() != reflect.Struct {
		return errors.New("invalid data")
	}

	if !s.hasHeader {
		s.writeHeader(data)
	}

	s.writeData(data)
	return nil
}

func (s *sheet) writeHeader(data interface{}) {
	t := reflect.TypeOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		column, ok := field.Tag.Lookup("column")
		if !ok {
			continue
		}

		header, ok := field.Tag.Lookup("header")
		if !ok {
			continue
		}

		axis := fmt.Sprintf("%s1", column)
		err := s.xls.SetCellValue(s.name, axis, header)
		if err != nil {
			log.Fatal(err)
		}
	}

	s.hasHeader = true
	s.line++
}

func (s *sheet) writeData(data interface{}) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		column, ok := field.Tag.Lookup("column")
		if !ok {
			continue
		}

		axis := fmt.Sprintf("%s%d", column, s.line)
		err := s.xls.SetCellValue(s.name, axis, v.Field(i).Interface())
		if err != nil {
			log.Fatal(err)
		}
	}

	s.line++
}
