package mcmap

type Biome int8

// Names and values from: http://www.minecraftwiki.net/wiki/Data_values

// Valid values for Biome
const (
	BioOcean               = 0
	BioPlains              = 1
	BioDesert              = 2
	BioExtremeHills        = 3
	BioForest              = 4
	BioTaiga               = 5
	BioSwampland           = 6
	BioRiver               = 7
	BioHell                = 8
	BioSky                 = 9
	BioFrozenOcean         = 10
	BioFrozenRiver         = 11
	BioIcePlains           = 12
	BioIceMountains        = 13
	BioMushroomIsland      = 14
	BioMushroomIslandShore = 15
	BioBeach               = 16
	BioDesertHills         = 17
	BioForestHills         = 18
	BioTaigaHills          = 19
	BioExtremeHillsEdge    = 20
	BioJungle              = 21
	BioJungleHills         = 22
	BioUncalculated        = -1
)

var biomeNames = map[Biome]string{
	BioOcean:               "Ocean",
	BioPlains:              "Plains",
	BioDesert:              "Desert",
	BioExtremeHills:        "Extreme Hills",
	BioForest:              "Forest",
	BioTaiga:               "Taiga",
	BioSwampland:           "Swampland",
	BioRiver:               "River",
	BioHell:                "Hell",
	BioSky:                 "Sky",
	BioFrozenOcean:         "Frozen Ocean",
	BioFrozenRiver:         "Frozen River",
	BioIcePlains:           "Ice Plains",
	BioIceMountains:        "Ice Mountains",
	BioMushroomIsland:      "Mushroom Island",
	BioMushroomIslandShore: "Mushroom Island Shore",
	BioBeach:               "Beach",
	BioDesertHills:         "Desert Hills",
	BioForestHills:         "Forest Hills",
	BioTaigaHills:          "Taiga Hills",
	BioExtremeHillsEdge:    "Extreme Hills Edge",
	BioJungle:              "Jungle",
	BioJungleHills:         "Jungle Hills",
	BioUncalculated:        "(Uncalculated)",
}

func (b Biome) String() string {
	if s, ok := biomeNames[b]; ok {
		return s
	}
	return "(Unknown)"
}
