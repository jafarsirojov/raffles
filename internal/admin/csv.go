package admin

import (
	"crm/internal/structs"
	"encoding/csv"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
)

func (s *service) createCSV(path string, leads []structs.Lead) error {
	file, err := os.Create(path)
	if err != nil {
		s.logger.Error("internal.admin.createCSV os.Create", zap.Error(err))
		return err
	}
	defer func() {
		err = file.Close()
		s.logger.Error("internal.admin.createCSV file.Close()", zap.Error(err))
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"№",
		"ID",
		"Site",
		"Name",
		"Phone",
		"Stage construction",
		"Region",
		"Type",
		"Purpose of acquisition",
		"Count of rooms",
		"Purchase budget",
		"Email",
		"Communication method",
		"Description",
		"Status",
		"Created at",
	}
	err = writer.Write(header)
	if err != nil {
		s.logger.Error("internal.admin.createCSV writer.Write", zap.Error(err))
		return err
	}

	for i, l := range leads {
		var record []string
		record = append(record,
			strconv.Itoa(i+1),
			strconv.Itoa(l.ID),
			l.Site,
			l.Name,
			l.Phone,
			l.REStageConstruction,
			l.RERegion,
			l.REType,
			l.PurchaseBudget,
			l.RECountOfRooms,
			l.PurchaseBudget,
			l.Email,
			l.CommunicationMethod,
			l.Description,
			l.Status,
			l.CreatedAt,
		)

		err = writer.Write(record)
		if err != nil {
			s.logger.Error("internal.admin.createCSV writer.Write 2",
				zap.Error(err), zap.Any("record", record))
			return err
		}
	}

	return nil
}

func (s *service) generateXLSXFromCSV(csvPath string, XLSXPath string) error {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		s.logger.Error("internal.admin.generateXLSXFromCSV os.Create", zap.Error(err))
		return err
	}
	defer func() {
		err = csvFile.Close()
		if err != nil {
			s.logger.Warn("internal.admin.generateXLSXFromCSV csvFile.Close()", zap.Error(err))
		}
	}()

	reader := csv.NewReader(csvFile)
	reader.Comma = rune(',')

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("leads")
	if err != nil {
		s.logger.Error("internal.admin.generateXLSXFromCSV xlsxFile.AddSheet", zap.Error(err))
		return err
	}

	fields, err := reader.Read()
	for err == nil {
		row := sheet.AddRow()
		for _, field := range fields {
			cell := row.AddCell()
			cell.Value = field
		}
		fields, err = reader.Read()
	}
	if err != io.EOF {
		s.logger.Error("internal.admin.generateXLSXFromCSV reader.Read() err != io.EOF", zap.Error(err))
		return err
	}

	err = xlsxFile.Save(XLSXPath)
	if err != nil {
		s.logger.Error("internal.admin.generateXLSXFromCSV xlsxFile.Save(XLSXPath)",
			zap.String("XLSXPath", XLSXPath), zap.Error(err))
		return err
	}

	return nil
}
