package export

import (
	"bytes"
	"fmt"

	"github.com/go-pdf/fpdf"
	"github.com/inmobilia/inmobilia-web/backend/internal/domain"
)

func GeneratePDF(bom domain.BillOfMaterials, cutPlan domain.CutPlan) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	if err := setupPDFFonts(pdf); err != nil {
		return nil, err
	}
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	setPDFFont(pdf, "B", 16)
	title := bom.FurnitureName
	if title == "" {
		title = bom.FurnitureID
	}
	pdf.Cell(0, 10, fmt.Sprintf("BOM / Planos de corte — %s", title))
	pdf.Ln(12)

	setPDFFont(pdf, "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("ID: %s", bom.FurnitureID))
	pdf.Ln(5)
	pdf.Cell(0, 6, fmt.Sprintf("Generado: %s", bom.GeneratedAt.Format("2006-01-02 15:04 UTC")))
	pdf.Ln(5)
	pdf.Cell(0, 6, fmt.Sprintf("Piezas: %d | Tableros: %.2f m2 | Tapacantos: %.2f m",
		bom.Summary.PartCount, bom.Summary.TotalBoardM2, bom.Summary.TotalEdgeM))
	pdf.Ln(10)

	writeSectionTitle(pdf, "Lista de materiales (piezas)")
	writePartsTable(pdf, bom.Parts)

	if len(bom.Hardware) > 0 {
		pdf.AddPage()
		writeSectionTitle(pdf, "Herrajes")
		writeHardwareTable(pdf, bom.Hardware)
	}

	if len(bom.EdgeBanding) > 0 {
		pdf.Ln(8)
		writeSectionTitle(pdf, "Tapacantos")
		writeEdgeTable(pdf, bom.EdgeBanding)
	}

	pdf.AddPage()
	writeSectionTitle(pdf, "Planos de corte")
	for _, sheet := range cutPlan.Sheets {
		setPDFFont(pdf, "B", 11)
		pdf.Cell(0, 7, fmt.Sprintf("%s — %.0f mm (%.2f m2)", sheet.MaterialName, sheet.Thickness, sheet.TotalAreaM2))
		pdf.Ln(8)
		writeCutTable(pdf, sheet.Parts)
		pdf.Ln(6)
	}

	if bom.Cost != nil {
		pdf.AddPage()
		writeSectionTitle(pdf, "Resumen de costos")
		writeCostSummary(pdf, *bom.Cost)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeSectionTitle(pdf *fpdf.Fpdf, title string) {
	setPDFFont(pdf, "B", 12)
	pdf.Cell(0, 8, title)
	pdf.Ln(8)
}

func writePartsTable(pdf *fpdf.Fpdf, parts []domain.BOMPartLine) {
	headers := []string{"Pieza", "Tipo", "Ancho", "Alto", "Esp.", "Material", "m2"}
	widths := []float64{42, 22, 18, 18, 14, 38, 16}
	writeTableHeader(pdf, headers, widths)

	setPDFFont(pdf, "", 8)
	for _, part := range parts {
		writeTableRow(pdf, widths, []string{
			truncate(part.Name, 24),
			part.Type,
			fmt.Sprintf("%.0f", part.Width),
			fmt.Sprintf("%.0f", part.Height),
			fmt.Sprintf("%.0f", part.Thickness),
			truncate(part.MaterialName, 20),
			fmt.Sprintf("%.3f", part.AreaM2),
		})
	}
}

func writeHardwareTable(pdf *fpdf.Fpdf, lines []domain.BOMHardwareLine) {
	headers := []string{"Herraje", "Cantidad", "Unit.", "Total"}
	widths := []float64{60, 25, 30, 30}
	writeTableHeader(pdf, headers, widths)

	setPDFFont(pdf, "", 9)
	for _, line := range lines {
		writeTableRow(pdf, widths, []string{
			line.Name,
			fmt.Sprintf("%d", line.Quantity),
			fmt.Sprintf("%.2f", line.UnitCost),
			fmt.Sprintf("%.2f", line.Total),
		})
	}
}

func writeEdgeTable(pdf *fpdf.Fpdf, lines []domain.BOMEdgeLine) {
	headers := []string{"Material", "Metros"}
	widths := []float64{80, 30}
	writeTableHeader(pdf, headers, widths)

	setPDFFont(pdf, "", 9)
	for _, line := range lines {
		writeTableRow(pdf, widths, []string{line.Material, fmt.Sprintf("%.2f", line.TotalLengthM)})
	}
}

func writeCutTable(pdf *fpdf.Fpdf, parts []domain.CutPartLine) {
	headers := []string{"Pieza", "Corte (mm)", "Esp.", "Veta", "Cant."}
	widths := []float64{50, 45, 18, 22, 15}
	writeTableHeader(pdf, headers, widths)

	setPDFFont(pdf, "", 8)
	for _, part := range parts {
		writeTableRow(pdf, widths, []string{
			truncate(part.Name, 28),
			fmt.Sprintf("%.0f x %.0f", part.Width, part.Height),
			fmt.Sprintf("%.0f", part.Thickness),
			part.Grain,
			fmt.Sprintf("%d", part.Quantity),
		})
	}
}

func writeCostSummary(pdf *fpdf.Fpdf, cost domain.CostResult) {
	setPDFFont(pdf, "", 10)
	for _, line := range cost.Materials {
		pdf.Cell(0, 6, fmt.Sprintf("Material: %s — %.2f m2 — %s %.2f", line.Name, line.AreaM2, cost.Currency, line.Total))
		pdf.Ln(5)
	}
	for _, line := range cost.Hardware {
		pdf.Cell(0, 6, fmt.Sprintf("Herraje: %s x%d — %s %.2f", line.Name, line.Quantity, cost.Currency, line.Total))
		pdf.Ln(5)
	}
	pdf.Cell(0, 6, fmt.Sprintf("Mano de obra: %.1f h — %s %.2f", cost.Labor.Hours, cost.Currency, cost.Labor.Total))
	pdf.Ln(5)
	pdf.Cell(0, 6, fmt.Sprintf("Desperdicio (%.0f%%): %s %.2f", cost.Waste.Percentage, cost.Currency, cost.Waste.Total))
	pdf.Ln(8)
	setPDFFont(pdf, "B", 12)
	pdf.Cell(0, 8, fmt.Sprintf("TOTAL: %s %.2f", cost.Currency, cost.Total))
	pdf.Ln(8)
}

func writeTableHeader(pdf *fpdf.Fpdf, headers []string, widths []float64) {
	setPDFFont(pdf, "B", 8)
	for i, header := range headers {
		pdf.CellFormat(widths[i], 7, header, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)
}

func writeTableRow(pdf *fpdf.Fpdf, widths []float64, values []string) {
	for i, value := range values {
		pdf.CellFormat(widths[i], 6, value, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max-3] + "..."
}
