package manufacturing

import "github.com/inmobilia/inmobilia-web/backend/internal/domain"

func addPartEdgeBanding(c *compileContext, partID string, edge domain.EdgeSide, length float64) {
	c.edgeBandingList = append(c.edgeBandingList, domain.EdgeBanding{
		PartID:   partID,
		Edge:     edge,
		Material: c.edgeBanding,
		Length:   length,
	})
}
