package helper

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// ExportFormat represents export file format
type ExportFormat string

const (
	FormatCSV   ExportFormat = "csv"
	FormatExcel ExportFormat = "excel"
	FormatPDF   ExportFormat = "pdf"
)

// ExportConfig represents export configuration
type ExportConfig struct {
	Format    ExportFormat
	Filename  string
	Headers   []string
	SheetName string // For Excel
	Title     string // For PDF
}

// DefaultExportConfig returns default export configuration
func DefaultExportConfig(format ExportFormat) ExportConfig {
	timestamp := time.Now().Format("20060102_150405")
	return ExportConfig{
		Format:    format,
		Filename:  fmt.Sprintf("export_%s.%s", timestamp, format),
		Headers:   []string{},
		SheetName: "Sheet1",
		Title:     "Data Export",
	}
}

// ExportToCSV exports data to CSV file
func ExportToCSV(data [][]string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return &InternalServerError{
			Message: fmt.Sprintf("Failed to create CSV file: %v", err),
			Order:   "H-Export-CSV-1",
		}
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return &InternalServerError{
				Message: fmt.Sprintf("Failed to write CSV row: %v", err),
				Order:   "H-Export-CSV-2",
			}
		}
	}

	return nil
}

// ExportToExcel exports data to Excel file
func ExportToExcel(data [][]string, config ExportConfig) error {
	f := excelize.NewFile()
	defer f.Close()

	// Create sheet
	sheetName := config.SheetName
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return &InternalServerError{
			Message: fmt.Sprintf("Failed to create sheet: %v", err),
			Order:   "H-Export-Excel-1",
		}
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Write headers with styling
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	// Write data
	for rowIdx, row := range data {
		for colIdx, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			f.SetCellValue(sheetName, cell, value)
			
			// Apply header style to first row
			if rowIdx == 0 {
				f.SetCellStyle(sheetName, cell, cell, headerStyle)
			}
		}
	}

	// Auto-fit columns
	for colIdx := range data[0] {
		col, _ := excelize.ColumnNumberToName(colIdx + 1)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// Save file
	if err := f.SaveAs(config.Filename); err != nil {
		return &InternalServerError{
			Message: fmt.Sprintf("Failed to save Excel file: %v", err),
			Order:   "H-Export-Excel-2",
		}
	}

	return nil
}

// ExportToPDF exports data to PDF file
func ExportToPDF(data [][]string, config ExportConfig) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, config.Title)
	pdf.Ln(12)

	// Table headers
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(68, 114, 196)
	pdf.SetTextColor(255, 255, 255)
	
	if len(data) > 0 {
		// Calculate column widths
		colWidth := 190.0 / float64(len(data[0]))
		
		// Write headers
		for _, header := range data[0] {
			pdf.CellFormat(colWidth, 7, header, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)

		// Table data
		pdf.SetFont("Arial", "", 10)
		pdf.SetFillColor(240, 240, 240)
		pdf.SetTextColor(0, 0, 0)
		fill := false

		for i, row := range data[1:] {
			for _, value := range row {
				pdf.CellFormat(colWidth, 7, value, "1", 0, "L", fill, 0, "")
			}
			pdf.Ln(-1)
			fill = !fill
			
			// Add new page if needed
			if i%35 == 0 && i > 0 {
				pdf.AddPage()
				// Repeat headers on new page
				pdf.SetFont("Arial", "B", 10)
				pdf.SetFillColor(68, 114, 196)
				pdf.SetTextColor(255, 255, 255)
				for _, header := range data[0] {
					pdf.CellFormat(colWidth, 7, header, "1", 0, "C", true, 0, "")
				}
				pdf.Ln(-1)
				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(0, 0, 0)
			}
		}
	}

	// Add footer with page numbers
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})

	// Save PDF
	if err := pdf.OutputFileAndClose(config.Filename); err != nil {
		return &InternalServerError{
			Message: fmt.Sprintf("Failed to save PDF file: %v", err),
			Order:   "H-Export-PDF-1",
		}
	}

	return nil
}

// ConvertToStringSlice converts struct slice to string slice for export
func ConvertToStringSlice(data interface{}, headers []string) [][]string {
	result := make([][]string, 0)
	
	// Add headers
	result = append(result, headers)
	
	// Use reflection to convert struct to string slice
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return result
	}
	
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		row := make([]string, 0)
		
		if item.Kind() == reflect.Struct {
			for _, header := range headers {
				field := item.FieldByName(header)
				if field.IsValid() {
					row = append(row, formatValue(field))
				} else {
					row = append(row, "")
				}
			}
		} else if item.Kind() == reflect.Map {
			for _, header := range headers {
				value := item.MapIndex(reflect.ValueOf(header))
				if value.IsValid() {
					row = append(row, formatValue(value))
				} else {
					row = append(row, "")
				}
			}
		}
		
		result = append(result, row)
	}
	
	return result
}

// formatValue formats reflect.Value to string
func formatValue(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}
	
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 2, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Struct:
		// Handle time.Time
		if t, ok := v.Interface().(time.Time); ok {
			return t.Format("2006-01-02 15:04:05")
		}
		return fmt.Sprintf("%v", v.Interface())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// ExportData is a generic function to export data in any format
func ExportData(data interface{}, headers []string, config ExportConfig) (string, error) {
	// Convert data to string slice
	stringData := ConvertToStringSlice(data, headers)
	
	if len(stringData) <= 1 {
		return "", &BadRequestError{
			Message: "No data to export",
			Order:   "H-Export-Data-1",
		}
	}
	
	// Export based on format
	switch config.Format {
	case FormatCSV:
		if err := ExportToCSV(stringData, config.Filename); err != nil {
			return "", err
		}
	case FormatExcel:
		if err := ExportToExcel(stringData, config); err != nil {
			return "", err
		}
	case FormatPDF:
		if err := ExportToPDF(stringData, config); err != nil {
			return "", err
		}
	default:
		return "", &BadRequestError{
			Message: "Invalid export format",
			Order:   "H-Export-Data-2",
		}
	}
	
	return config.Filename, nil
}
