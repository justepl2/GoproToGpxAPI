package request

import "errors"

type Terrain string

const (
	TerrainRoad    Terrain = "road"
	TerrainOffroad Terrain = "offroad"
)

type LinkVideos struct {
	VideoIds []string `json:"videoIds" example:"5f5e3e4e-3e4e-5f5e-3e4e-5f5e3e4e3e4e"`
	Terrain  Terrain  `json:"terrain" example:"offroad"`
}

func (lv *LinkVideos) Validate() error {
	if len(lv.VideoIds) < 2 {
		return errors.New("at least 2 video ids are required")
	}

	if lv.Terrain == "" {
		return errors.New("terrain is required")
	}

	if !lv.Terrain.IsValid() {
		return errors.New("invalid terrain")
	}

	return nil
}

func (t Terrain) IsValid() bool {
	switch t {
	case TerrainRoad, TerrainOffroad:
		return true
	}
	return false
}
