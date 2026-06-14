package export

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-pdf/fpdf"
)

//go:embed fonts/DejaVuSans.ttf
//go:embed fonts/DejaVuSans-Bold.ttf
var embeddedFonts embed.FS

const pdfFontFamily = "DejaVu"

var (
	fontDirOnce sync.Once
	fontDir     string
	fontDirErr  error
)

func prepareFontDir() (string, error) {
	fontDirOnce.Do(func() {
		dir, err := os.MkdirTemp("", "inmobilia-pdf-fonts-*")
		if err != nil {
			fontDirErr = err
			return
		}
		for _, name := range []string{"DejaVuSans.ttf", "DejaVuSans-Bold.ttf"} {
			data, err := embeddedFonts.ReadFile("fonts/" + name)
			if err != nil {
				fontDirErr = fmt.Errorf("read font %s: %w", name, err)
				return
			}
			if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
				fontDirErr = err
				return
			}
		}
		fontDir = dir
	})
	return fontDir, fontDirErr
}

func setupPDFFonts(pdf *fpdf.Fpdf) error {
	dir, err := prepareFontDir()
	if err != nil {
		return err
	}
	pdf.AddUTF8Font(pdfFontFamily, "", filepath.Join(dir, "DejaVuSans.ttf"))
	pdf.AddUTF8Font(pdfFontFamily, "B", filepath.Join(dir, "DejaVuSans-Bold.ttf"))
	return nil
}

func setPDFFont(pdf *fpdf.Fpdf, style string, size float64) {
	pdf.SetFont(pdfFontFamily, style, size)
}
