package mcmap

type Biome uint8

// Names and values from: http://www.minecraftwiki.net/wiki/Data_values

// Valid values for Biome
const (
	BioOcean                = 0
	BioPlains               = 1
	BioDesert               = 2
	BioExtremeHills         = 3
	BioForest               = 4
	BioTaiga                = 5
	BioSwampland            = 6
	BioRiver                = 7
	BioHell                 = 8
	BioSky                  = 9
	BioFrozenOcean          = 10
	BioFrozenRiver          = 11
	BioIcePlains            = 12
	BioIceMountains         = 13
	BioMushroomIsland       = 14
	BioMushroomIslandShore  = 15
	BioBeach                = 16
	BioDesertHills          = 17
	BioForestHills          = 18
	BioTaigaHills           = 19
	BioExtremeHillsEdge     = 20
	BioJungle               = 21
	BioJungleHills          = 22
	BioJungleEdge           = 23
	BioDeepOcean            = 24
	BioStoneBeach           = 25
	BioColdBeach            = 26
	BioBirchForest          = 27
	BioBirchForestHills     = 28
	BioRoofedForest         = 29
	BioColdTaiga            = 30
	BioColdTaigaHills       = 31
	BioMegaTaiga            = 32
	BioMegaTaigaHills       = 33
	BioExtremeHillsPlus     = 34
	BioSavanna              = 35
	BioSavannaPlateau       = 36
	BioMesa                 = 37
	BioMesaPlateauF         = 38
	BioMesaPlateau          = 39
	BioSunflowerPlains      = 129
	BioDesertM              = 130
	BioExtremeHillsM        = 131
	BioFlowerForest         = 132
	BioTaigaM               = 133
	BioSwamplandM           = 134
	BioIcePlainsSpikes      = 140
	BioIceMountainsSpikes   = 141
	BioJungleM              = 149
	BioJungleEdgeM          = 151
	BioBirchForestM         = 155
	BioBirchForestHillsM    = 156
	BioRoofedForestM        = 157
	BioColdTaigaM           = 158
	BioMegaSpruceTaiga      = 160
	BioMegaSpruceTaigaHills = 161
	BioExtremeHillsPlusM    = 162
	BioSavannaM             = 163
	BioSavannaPlateauM      = 164
	BioMesaBryce            = 165
	BioMesaPlateauFM        = 166
	BioMesaPlateauM         = 167
	BioUncalculated         = 0xff // (-1)
)

var biomeNames = map[Biome]string{
	BioOcean:                "Ocean",
	BioPlains:               "Plains",
	BioDesert:               "Desert",
	BioExtremeHills:         "Extreme Hills",
	BioForest:               "Forest",
	BioTaiga:                "Taiga",
	BioSwampland:            "Swampland",
	BioRiver:                "River",
	BioHell:                 "Hell",
	BioSky:                  "Sky",
	BioFrozenOcean:          "Frozen Ocean",
	BioFrozenRiver:          "Frozen River",
	BioIcePlains:            "Ice Plains",
	BioIceMountains:         "Ice Mountains",
	BioMushroomIsland:       "Mushroom Island",
	BioMushroomIslandShore:  "Mushroom Island Shore",
	BioBeach:                "Beach",
	BioDesertHills:          "Desert Hills",
	BioForestHills:          "Forest Hills",
	BioTaigaHills:           "Taiga Hills",
	BioExtremeHillsEdge:     "Extreme Hills Edge",
	BioJungle:               "Jungle",
	BioJungleHills:          "Jungle Hills",
	BioJungleEdge:           "Jungle Edge",
	BioDeepOcean:            "Deep Ocean",
	BioStoneBeach:           "Stone Beach",
	BioColdBeach:            "Cold Beach",
	BioBirchForest:          "Birch Forest",
	BioBirchForestHills:     "Birch Forest Hills",
	BioRoofedForest:         "Roofed Forest",
	BioColdTaiga:            "Cold Taiga",
	BioColdTaigaHills:       "Cold Taiga Hills",
	BioMegaTaiga:            "Mega Taiga",
	BioMegaTaigaHills:       "Mega Taiga Hills",
	BioExtremeHillsPlus:     "Extreme Hills+",
	BioSavanna:              "Savanna",
	BioSavannaPlateau:       "Savanna Plateau",
	BioMesa:                 "Mesa",
	BioMesaPlateauF:         "Mesa Plateau F",
	BioMesaPlateau:          "Mesa Plateau",
	BioSunflowerPlains:      "Sunflower Plains",
	BioDesertM:              "Desert M",
	BioExtremeHillsM:        "Extreme Hills M",
	BioFlowerForest:         "Flower Forest",
	BioTaigaM:               "Taiga M",
	BioSwamplandM:           "Swampland M",
	BioIcePlainsSpikes:      "Ice Plains Spikes",
	BioIceMountainsSpikes:   "Ice Mountains Spikes",
	BioJungleM:              "Jungle M",
	BioJungleEdgeM:          "JungleEdge M",
	BioBirchForestM:         "Birch Forest M",
	BioBirchForestHillsM:    "Birch Forest Hills M",
	BioRoofedForestM:        "Roofed Forest M",
	BioColdTaigaM:           "Cold Taiga M",
	BioMegaSpruceTaiga:      "Mega Spruce Taiga",
	BioMegaSpruceTaigaHills: "Mega Spruce Taiga Hills",
	BioExtremeHillsPlusM:    "Extreme Hills+ M",
	BioSavannaM:             "Savanna M",
	BioSavannaPlateauM:      "Savanna Plateau M",
	BioMesaBryce:            "Mesa (Bryce)",
	BioMesaPlateauFM:        "Mesa Plateau F M",
	BioMesaPlateauM:         "Mesa Plateau M",
	BioUncalculated:         "(Uncalculated)",
}

func (b Biome) String() string {
	if s, ok := biomeNames[b]; ok {
		return s
	}
	return "(Unknown)"
}
