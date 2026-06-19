package manufacturing

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

const thickEdgeBandingMaterial = "pvc-thick-3mm"

func addPartEdgeBanding(c *compileContext, partID string, edge domain.EdgeSide, length float64) {
	c.edgeBandingList = append(c.edgeBandingList, domain.EdgeBanding{
		PartID:   partID,
		Edge:     edge,
		Material: c.edgeBanding,
		Length:   length,
	})
}

func addPartThickTopEdgeBanding(c *compileContext, partID string, length float64) {
	c.edgeBandingList = append(c.edgeBandingList, domain.EdgeBanding{
		PartID:   partID,
		Edge:     domain.EdgeTop,
		Material: thickEdgeBandingMaterial,
		Length:   length,
	})
}
