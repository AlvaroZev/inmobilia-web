package domain

// RoomGeometry (Layer 1) represents the physical environment.

type RoomGeometry struct {
	ID        string     `json:"id"`
	Name      string     `json:"name,omitempty"`
	Perimeter Polygon2D  `json:"perimeter"`
	Floor     Plane      `json:"floor"`
	Ceiling   Plane      `json:"ceiling"`
	Walls     []Wall     `json:"walls"`
	Openings  []Opening  `json:"openings"`
	Obstacles []Obstacle `json:"obstacles"`
}

type Wall struct {
	ID          string    `json:"id"`
	Vertices    []Point3D `json:"vertices"`
	Thickness   float64   `json:"thickness"`
	OutOfPlumb  *Point3D  `json:"outOfPlumb,omitempty"`
}

type OpeningType string

const (
	OpeningDoor   OpeningType = "door"
	OpeningWindow OpeningType = "window"
)

type Opening struct {
	ID       string      `json:"id"`
	Type     OpeningType `json:"type"`
	WallID   string      `json:"wallId"`
	Origin   Point3D     `json:"origin"`
	Width    float64     `json:"width"`
	Height   float64     `json:"height"`
}

type ObstacleType string

const (
	ObstacleColumn   ObstacleType = "column"
	ObstacleBeam     ObstacleType = "beam"
	ObstaclePipe     ObstacleType = "pipe"
	ObstacleSkirting ObstacleType = "skirting"
	ObstacleOther    ObstacleType = "other"
)

type Obstacle struct {
	ID       string       `json:"id"`
	Type     ObstacleType `json:"type"`
	Label    string       `json:"label,omitempty"`
	Bounds   BoundingBox3D `json:"bounds"`
	Profile  []Point3D    `json:"profile,omitempty"`
}
