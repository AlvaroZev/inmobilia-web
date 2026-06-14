package domain

// InstallationConstraints (Layer 3) connects room geometry and furniture definition.

type InstallationConstraints struct {
	ID          string                  `json:"id"`
	Zone        InstallationZone        `json:"zone"`
	Clearances  Clearances              `json:"clearances"`
	Tolerances  Tolerances              `json:"tolerances"`
	References  InstallationReferences  `json:"references"`
}

type InstallationZone struct {
	AnchorWallIDs []string         `json:"anchorWallIds,omitempty"`
	Bounds        *BoundingBox3D   `json:"bounds,omitempty"`
}

type Clearances struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
	Back   float64 `json:"back"`
	Front  float64 `json:"front"`
}

type Tolerances struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Depth  float64 `json:"depth"`
}

type InstallationReferences struct {
	FloorOffset       float64 `json:"floorOffset"`
	CeilingOffset     float64 `json:"ceilingOffset,omitempty"`
	WallOffset        float64 `json:"wallOffset,omitempty"`
	ReferenceWallID   string  `json:"referenceWallId,omitempty"`
}
