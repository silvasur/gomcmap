package mcmap

import (
	"github.com/kch42/gonbt/nbt"
)

type BlockID uint16

type Block struct {
	ID                   BlockID
	Data                 byte            // Actually only a half-byte.
	BlockLight, SkyLight byte            // Also, only half-bytes.
	TileEntity           nbt.TagCompound // The x, y and z values in here can be ignored, will automatically be fixed on saving. Will be nil, if no TileEntity is available.
	Tick                 *TileTick       // If nil, no TileTick info is available for this block
}

type TileTick struct {
	i, t, p int32
	hasP    bool
}

func (tt *TileTick) I() int32 { return tt.i }
func (tt *TileTick) T() int32 { return tt.t }
func (tt *TileTick) P() int32 { return tt.p }

func (tt *TileTick) SetI(i int32) { tt.i = i }
func (tt *TileTick) SetT(t int32) { tt.t = t }

func (tt *TileTick) SetP(p int32) {
	tt.p = p
	tt.hasP = true
}

// Names and values from: http://www.minecraftwiki.net/wiki/Data_values

// Valid values for BlockID
const (
	BlkAir                        = 0
	BlkStone                      = 1
	BlkGrassBlock                 = 2
	BlkDirt                       = 3
	BlkCobblestone                = 4
	BlkWoodPlanks                 = 5
	BlkSaplings                   = 6
	BlkBedrock                    = 7
	BlkWater                      = 8
	BlkStationaryWater            = 9
	BlkLava                       = 10
	BlkStationaryLava             = 11
	BlkSand                       = 12
	BlkGravel                     = 13
	BlkGoldOre                    = 14
	BlkIronOre                    = 15
	BlkCoalOre                    = 16
	BlkWood                       = 17
	BlkLeaves                     = 18
	BlkSponge                     = 19
	BlkGlass                      = 20
	BlkLapisLazuliOre             = 21
	BlkLapisLazuliBlock           = 22
	BlkDispenser                  = 23
	BlkSandstone                  = 24
	BlkNoteBlock                  = 25
	BlkBed                        = 26
	BlkPoweredRail                = 27
	BlkDetectorRail               = 28
	BlkStickyPiston               = 29
	BlkCobweb                     = 30
	BlkGrass                      = 31
	BlkDeadBush                   = 32
	BlkPiston                     = 33
	BlkPistonExtension            = 34
	BlkWool                       = 35
	BlkBlockMovedByPiston         = 36
	BlkDandelion                  = 37
	BlkFlower                     = 38
	BlkBrownMushroom              = 39
	BlkRedMushroom                = 40
	BlkBlockOfGold                = 41
	BlkBlockOfIron                = 42
	BlkDoubleSlabs                = 43
	BlkSlabs                      = 44
	BlkBricks                     = 45
	BlkTNT                        = 46
	BlkBookshelf                  = 47
	BlkMossStone                  = 48
	BlkObsidian                   = 49
	BlkTorch                      = 50
	BlkFire                       = 51
	BlkMonsterSpawner             = 52
	BlkOakWoodStairs              = 53
	BlkChest                      = 54
	BlkRedstoneWire               = 55
	BlkDiamondOre                 = 56
	BlkBlockOfDiamond             = 57
	BlkCraftingTable              = 58
	BlkWheat                      = 59
	BlkFarmland                   = 60
	BlkFurnace                    = 61
	BlkBurningFurnace             = 62
	BlkSignPost                   = 63
	BlkWoodenDoor                 = 64
	BlkLadders                    = 65
	BlkRail                       = 66
	BlkCobblestoneStairs          = 67
	BlkWallSign                   = 68
	BlkLever                      = 69
	BlkStonePressurePlate         = 70
	BlkIronDoor                   = 71
	BlkWoodenPressurePlate        = 72
	BlkRedstoneOre                = 73
	BlkGlowingRedstoneOre         = 74
	BlkRedstoneTorchInactive      = 75
	BlkRedstoneTorchActive        = 76
	BlkStoneButton                = 77
	BlkSnow                       = 78
	BlkIce                        = 79
	BlkSnowBlock                  = 80
	BlkCactus                     = 81
	BlkClay                       = 82
	BlkSugarCane                  = 83
	BlkJukebox                    = 84
	BlkFence                      = 85
	BlkPumpkin                    = 86
	BlkNetherrack                 = 87
	BlkSoulSand                   = 88
	BlkGlowstone                  = 89
	BlkNetherPortal               = 90
	BlkJackOLantern               = 91
	BlkCakeBlock                  = 92
	BlkRedstoneRepeaterInactive   = 93
	BlkRedstoneRepeaterActive     = 94
	BlkStainedGlass               = 95
	BlkLockedChest                = 95
	BlkTrapdoor                   = 96
	BlkMonsterEgg                 = 97
	BlkStoneBricks                = 98
	BlkHugeBrownMushroom          = 99
	BlkHugeRedMushroom            = 100
	BlkIronBars                   = 101
	BlkGlassPane                  = 102
	BlkMelon                      = 103
	BlkPumpkinStem                = 104
	BlkMelonStem                  = 105
	BlkVines                      = 106
	BlkFenceGate                  = 107
	BlkBrickStairs                = 108
	BlkStoneBrickStairs           = 109
	BlkMycelium                   = 110
	BlkLilyPad                    = 111
	BlkNetherBrick                = 112
	BlkNetherBrickFence           = 113
	BlkNetherBrickStairs          = 114
	BlkNetherWart                 = 115
	BlkEnchantmentTable           = 116
	BlkBrewingStand               = 117
	BlkCauldron                   = 118
	BlkEndPortal                  = 119
	BlkEndPortalBlock             = 120
	BlkEndStone                   = 121
	BlkDragonEgg                  = 122
	BlkRedstoneLampInactive       = 123
	BlkRedstoneLampActive         = 124
	BlkWoodenDoubleSlab           = 125
	BlkWoodenSlab                 = 126
	BlkCocoa                      = 127
	BlkSandstoneStairs            = 128
	BlkEmeraldOre                 = 129
	BlkEnderChest                 = 130
	BlkTripwireHook               = 131
	BlkTripwire                   = 132
	BlkBlockOfEmerald             = 133
	BlkSpruceWoodStairs           = 134
	BlkBirchWoodStairs            = 135
	BlkJungleWoodStairs           = 136
	BlkCommandBlock               = 137
	BlkBeacon                     = 138
	BlkCobblestoneWall            = 139
	BlkFlowerPot                  = 140
	BlkCarrots                    = 141
	BlkPotatoes                   = 142
	BlkWoodenButton               = 143
	BlkMobHead                    = 144
	BlkAnvil                      = 145
	BlkTrappedChest               = 146
	BlkWeightedPressurePlateLight = 147
	BlkWeightedPressurePlateHeavy = 148
	BlkRedstoneComparatorInactive = 149
	BlkRedstoneComparatorActive   = 150
	BlkDaylightSensor             = 151
	BlkBlockOfRedstone            = 152
	BlkNetherQuartzOre            = 153
	BlkHopper                     = 154
	BlkBlockOfQuartz              = 155
	BlkQuartzStairs               = 156
	BlkActivatorRail              = 157
	BlkDropper                    = 158
	BlkStainedClay                = 159
	BlkStainedGlassPane           = 160
	BlkWood2                      = 162
	BlkAcaciaWoodStairs           = 163
	BlkDarkOakWoodStairs          = 164
	BlkSlimeBlock                 = 165
	BlkHayBlock                   = 170
	BlkCarpet                     = 171
	BlkHardenedClay               = 172
	BlkBlockOfCoal                = 173
	BlkPackedIce                  = 174
	BlkLargeFlower                = 175
	BlkBarrier                    = 166

	// Aliases
	BlkRose      = BlkFlower
	BlkPoppy     = BlkFlower
	BlkInvisible = BlkBarrier
)

var blockNames = map[BlockID]string{
	BlkAir:                        "Air",
	BlkStone:                      "Stone",
	BlkGrassBlock:                 "Grass Block",
	BlkDirt:                       "Dirt",
	BlkCobblestone:                "Cobblestone",
	BlkWoodPlanks:                 "Wood Planks",
	BlkSaplings:                   "Saplings",
	BlkBedrock:                    "Bedrock",
	BlkWater:                      "Water",
	BlkStationaryWater:            "Stationary water",
	BlkLava:                       "Lava",
	BlkStationaryLava:             "Stationary lava",
	BlkSand:                       "Sand",
	BlkGravel:                     "Gravel",
	BlkGoldOre:                    "Gold Ore",
	BlkIronOre:                    "Iron Ore",
	BlkCoalOre:                    "Coal Ore",
	BlkWood:                       "Wood",
	BlkLeaves:                     "Leaves",
	BlkSponge:                     "Sponge",
	BlkGlass:                      "Glass",
	BlkLapisLazuliOre:             "Lapis Lazuli Ore",
	BlkLapisLazuliBlock:           "Lapis Lazuli Block",
	BlkDispenser:                  "Dispenser",
	BlkSandstone:                  "Sandstone",
	BlkNoteBlock:                  "Note Block",
	BlkBed:                        "Bed",
	BlkPoweredRail:                "Powered Rail",
	BlkDetectorRail:               "Detector Rail",
	BlkStickyPiston:               "Sticky Piston",
	BlkCobweb:                     "Cobweb",
	BlkGrass:                      "Grass",
	BlkDeadBush:                   "Dead Bush",
	BlkPiston:                     "Piston",
	BlkPistonExtension:            "Piston Extension",
	BlkWool:                       "Wool",
	BlkBlockMovedByPiston:         "Block moved by Piston",
	BlkDandelion:                  "Dandelion",
	BlkFlower:                     "Flower",
	BlkBrownMushroom:              "Brown Mushroom",
	BlkRedMushroom:                "Red Mushroom",
	BlkBlockOfGold:                "Block of Gold",
	BlkBlockOfIron:                "Block of Iron",
	BlkDoubleSlabs:                "Double Slabs",
	BlkSlabs:                      "Slabs",
	BlkBricks:                     "Bricks",
	BlkTNT:                        "TNT",
	BlkBookshelf:                  "Bookshelf",
	BlkMossStone:                  "Moss Stone",
	BlkObsidian:                   "Obsidian",
	BlkTorch:                      "Torch",
	BlkFire:                       "Fire",
	BlkMonsterSpawner:             "Monster Spawner",
	BlkOakWoodStairs:              "Oak Wood Stairs",
	BlkChest:                      "Chest",
	BlkRedstoneWire:               "Redstone Wire",
	BlkDiamondOre:                 "Diamond Ore",
	BlkBlockOfDiamond:             "Block of Diamond",
	BlkCraftingTable:              "Crafting Table",
	BlkWheat:                      "Wheat",
	BlkFarmland:                   "Farmland",
	BlkFurnace:                    "Furnace",
	BlkBurningFurnace:             "Burning Furnace",
	BlkSignPost:                   "Sign Post",
	BlkWoodenDoor:                 "Wooden Door",
	BlkLadders:                    "Ladders",
	BlkRail:                       "Rail",
	BlkCobblestoneStairs:          "Cobblestone Stairs",
	BlkWallSign:                   "Wall Sign",
	BlkLever:                      "Lever",
	BlkStonePressurePlate:         "Stone Pressure Plate",
	BlkIronDoor:                   "Iron Door",
	BlkWoodenPressurePlate:        "Wooden Pressure Plate",
	BlkRedstoneOre:                "Redstone Ore",
	BlkGlowingRedstoneOre:         "Glowing Redstone Ore",
	BlkRedstoneTorchInactive:      "Redstone Torch (inactive)",
	BlkRedstoneTorchActive:        "Redstone Torch (active)",
	BlkStoneButton:                "Stone Button",
	BlkSnow:                       "Snow",
	BlkIce:                        "Ice",
	BlkSnowBlock:                  "Snow Block",
	BlkCactus:                     "Cactus",
	BlkClay:                       "Clay",
	BlkSugarCane:                  "Sugar Cane",
	BlkJukebox:                    "Jukebox",
	BlkFence:                      "Fence",
	BlkPumpkin:                    "Pumpkin",
	BlkNetherrack:                 "Netherrack",
	BlkSoulSand:                   "Soul Sand",
	BlkGlowstone:                  "Glowstone",
	BlkNetherPortal:               "Nether Portal",
	BlkJackOLantern:               "Jack 'o' Lantern",
	BlkCakeBlock:                  "Cake Block",
	BlkRedstoneRepeaterInactive:   "Redstone Repeater (inactive)",
	BlkRedstoneRepeaterActive:     "Redstone Repeater (active)",
	BlkStainedGlass:               "Stained Glass",
	BlkTrapdoor:                   "Trapdoor",
	BlkMonsterEgg:                 "Monster Egg",
	BlkStoneBricks:                "Stone Bricks",
	BlkHugeBrownMushroom:          "Huge Brown Mushroom",
	BlkHugeRedMushroom:            "Huge Red Mushroom",
	BlkIronBars:                   "Iron Bars",
	BlkGlassPane:                  "Glass Pane",
	BlkMelon:                      "Melon",
	BlkPumpkinStem:                "Pumpkin Stem",
	BlkMelonStem:                  "Melon Stem",
	BlkVines:                      "Vines",
	BlkFenceGate:                  "Fence Gate",
	BlkBrickStairs:                "Brick Stairs",
	BlkStoneBrickStairs:           "Stone Brick Stairs",
	BlkMycelium:                   "Mycelium",
	BlkLilyPad:                    "Lily Pad",
	BlkNetherBrick:                "Nether Brick",
	BlkNetherBrickFence:           "Nether Brick Fence",
	BlkNetherBrickStairs:          "Nether Brick Stairs",
	BlkNetherWart:                 "Nether Wart",
	BlkEnchantmentTable:           "Enchantment Table",
	BlkBrewingStand:               "Brewing Stand",
	BlkCauldron:                   "Cauldron",
	BlkEndPortal:                  "End Portal",
	BlkEndPortalBlock:             "End Portal Block",
	BlkEndStone:                   "End Stone",
	BlkDragonEgg:                  "Dragon Egg",
	BlkRedstoneLampInactive:       "Redstone Lamp (inactive)",
	BlkRedstoneLampActive:         "Redstone Lamp (active)",
	BlkWoodenDoubleSlab:           "Wooden Double Slab",
	BlkWoodenSlab:                 "Wooden Slab",
	BlkCocoa:                      "Cocoa",
	BlkSandstoneStairs:            "Sandstone Stairs",
	BlkEmeraldOre:                 "Emerald Ore",
	BlkEnderChest:                 "Ender Chest",
	BlkTripwireHook:               "Tripwire Hook",
	BlkTripwire:                   "Tripwire",
	BlkBlockOfEmerald:             "Block of Emerald",
	BlkSpruceWoodStairs:           "Spruce Wood Stairs",
	BlkBirchWoodStairs:            "Birch Wood Stairs",
	BlkJungleWoodStairs:           "Jungle Wood Stairs",
	BlkCommandBlock:               "Command Block",
	BlkBeacon:                     "Beacon",
	BlkCobblestoneWall:            "Cobblestone Wall",
	BlkFlowerPot:                  "Flower Pot",
	BlkCarrots:                    "Carrots",
	BlkPotatoes:                   "Potatoes",
	BlkWoodenButton:               "Wooden Button",
	BlkMobHead:                    "Mob Head",
	BlkAnvil:                      "Anvil",
	BlkTrappedChest:               "Trapped Chest",
	BlkWeightedPressurePlateLight: "Weighted Pressure Plate (Light)",
	BlkWeightedPressurePlateHeavy: "Weighted Pressure Plate (Heavy)",
	BlkRedstoneComparatorInactive: "Redstone Comparator (inactive)",
	BlkRedstoneComparatorActive:   "Redstone Comparator (active)",
	BlkDaylightSensor:             "Daylight Sensor",
	BlkBlockOfRedstone:            "Block of Redstone",
	BlkNetherQuartzOre:            "Nether Quartz Ore",
	BlkHopper:                     "Hopper",
	BlkBlockOfQuartz:              "Block of Quartz",
	BlkQuartzStairs:               "Quartz Stairs",
	BlkActivatorRail:              "Activator Rail",
	BlkDropper:                    "Dropper",
	BlkStainedClay:                "Stained Clay",
	BlkHayBlock:                   "Hay Block",
	BlkCarpet:                     "Carpet",
	BlkHardenedClay:               "Hardened Clay",
	BlkBlockOfCoal:                "Block of Coal",
	BlkPackedIce:                  "Packed Ice",
	BlkLargeFlower:                "Large Flower",
	BlkStainedGlassPane:           "Stained Glass Pane",
	BlkWood2:                      "Wood 2",
	BlkAcaciaWoodStairs:           "Acacia Wood Stairs",
	BlkDarkOakWoodStairs:          "Dark Oak Wood Stairs",
	BlkSlimeBlock:                 "Slime Block",
	BlkBarrier:                    "Barrier",
}

func (b BlockID) String() string {
	if s, ok := blockNames[b]; ok {
		return s
	}

	return "(unused)"
}
