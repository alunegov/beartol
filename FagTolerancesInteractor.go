package beartol

import (
	"fmt"
)

// FagTolerancesInteractor implements TolerancesInteractor.
type FagTolerancesInteractor struct{}

func (thiz *FagTolerancesInteractor) GetInnerDiameterTolerance(rb RollingBearing) (int, int, error) {
	id, err := thiz.toToleranceId(rb.Type)
	if err != nil {
		return 0, 0, err
	}
	tolerances := gDiameterTolerances[id][rb.ClassTochn].Inner
	return thiz.findTolerance(tolerances, rb.InnerDiameter)
}

func (thiz *FagTolerancesInteractor) GetOuterDiameterTolerance(rb RollingBearing) (int, int, error) {
	id, err := thiz.toToleranceId(rb.Type)
	if err != nil {
		return 0, 0, err
	}
	tolerances := gDiameterTolerances[id][rb.ClassTochn].Outer
	return thiz.findTolerance(tolerances, rb.OuterDiameter)
}

func (thiz *FagTolerancesInteractor) toToleranceId(rbType RollingBearingType) (int, error) {
	switch rbType {
	case RollingBearingType_BallRadial:
		fallthrough
	case RollingBearingType_BallRadialUpor:
		// у 72B и 73B допуски внутреннего диаметра для подшипников всех классов точности - по P5 (без дополнительного
		// обозначения)
		fallthrough
	case RollingBearingType_Shpindel:
		fallthrough
	case RollingBearingType_4ptContact:
		fallthrough
	case RollingBearingType_BallSphere_Cylynder:
		fallthrough
	case RollingBearingType_BallSphere_Kone:
		fallthrough
	case RollingBearingType_RollerCylynder:
		fallthrough
	case RollingBearingType_RollerSphere_Cylynder:
		// отдельная таблица на стр. 366 for T41A
		fallthrough
	case RollingBearingType_RollerSphere_Kone:
		// ограничение только на наружный диаметр
		fallthrough
	case RollingBearingType_RollerSphere_Bochka_Cylynder:
		fallthrough
	case RollingBearingType_RollerSphere_Bochka_Kone:
		return 0, nil
	case RollingBearingType_RollerKone:
		// default P6X for 320X, 329, 330, 331, 332 and d < 200 mm
		return 1, nil
	case RollingBearingType_BallUpor:
		fallthrough
	case RollingBearingType_BallRadialUpor2:
		// отдельная таблица на стр. 469 for 7602, 7603
		fallthrough
	case RollingBearingType_RollerCylynderUpor:
		fallthrough
	case RollingBearingType_RollerSphereUpor:
		return 2, nil
	}

	return 0, fmt.Errorf("unsupp bearing type %d", rbType)
}

func (thiz *FagTolerancesInteractor) findTolerance(tolerances []diameterTolerance, diameter int) (int, int, error) {
	for _, it := range tolerances {
		if (it.DiameterLow < diameter) && (diameter <= it.DiameterHigh) {
			return it.ToleranceA, it.ToleranceB, nil
		}
	}
	return 0, 0, fmt.Errorf("unsupp diameter %d", diameter)
}

func (thiz *FagTolerancesInteractor) GetClearance(rb RollingBearing) (int, int, error) {
	id, err := thiz.toClearanceId(rb.Type, rb.TypeSpecial)
	if err != nil {
		return 0, 0, err
	}
	clearances := gClearances[id]
	clearance, err := thiz.findClearance(clearances, rb.InnerDiameter)
	if err != nil {
		return 0, 0, err
	}
	clearanceRange := clearance[rb.ClearanceGroup]
	return clearanceRange.Min, clearanceRange.Max, nil
}

func (thiz *FagTolerancesInteractor) toClearanceId(rbType RollingBearingType, rbTypeSpecial RollingBearingTypeSpecial) (int, error) {
	switch rbType {
	case RollingBearingType_BallRadial:
		return 0, nil
	case RollingBearingType_BallRadialUpor:
		if rbTypeSpecial == RollingBearingTypeSpecial_72B_73B {
			return 0, fmt.Errorf("no clearance for bearing type %d (%d)", rbType, rbTypeSpecial)
		} else if rbTypeSpecial == RollingBearingTypeSpecial_33DA {
			return 4, nil
		} else {
			return 3, nil
		}
	case RollingBearingType_4ptContact:
		return 5, nil
	case RollingBearingType_BallSphere_Cylynder:
		return 1, nil
	case RollingBearingType_BallSphere_Kone:
		return 2, nil
	case RollingBearingType_RollerCylynder:
		if rbTypeSpecial == RollingBearingTypeSpecial_NN30ASK {
			// default C1NA for NN30ASK
			return 7, nil
		} else {
			return 6, nil
		}
	case RollingBearingType_RollerSphere_Cylynder:
		// default C4 for T41A
		return 8, nil
	case RollingBearingType_RollerSphere_Kone:
		return 9, nil
	case RollingBearingType_RollerSphere_Bochka_Cylynder:
		return 10, nil
	case RollingBearingType_RollerSphere_Bochka_Kone:
		return 11, nil
	case RollingBearingType_Shpindel:
		fallthrough
	case RollingBearingType_RollerKone:
		fallthrough
	case RollingBearingType_BallUpor:
		fallthrough
	case RollingBearingType_BallRadialUpor2:
		fallthrough
	case RollingBearingType_RollerCylynderUpor:
		fallthrough
	case RollingBearingType_RollerSphereUpor:
		return 0, fmt.Errorf("no clearance for bearing type %d (%d)", rbType, rbTypeSpecial)
	}

	return 0, fmt.Errorf("unsupp bearing type %d (%d)", rbType, rbTypeSpecial)
}

func (thiz *FagTolerancesInteractor) findClearance(clearances []clearance, diameter int) (map[string]clearanceRange, error) {
	for _, it := range clearances {
		if (it.DiameterLow < diameter) && (diameter <= it.DiameterHigh) {
			return it.Clearance, nil
		}
	}
	return nil, fmt.Errorf("unsupp diameter %d", diameter)
}

type diameterTolerance struct {
	DiameterLow  int
	DiameterHigh int
	ToleranceA   int
	ToleranceB   int
}

/*
0 допуски радиальных подшипников (кроме конических роликоподшипников)
  P4S шпиндельные
  SP, UP двухрядные роликоподшипники с цилиндрическими роликами
1 допуски конических роликоподшипников
2 допуски упорных подшипников
  SP радиально-упорные шарикоподшипники, серии 2344 и 2347
*/
var gDiameterTolerances = map[int]map[string]struct {
	Inner []diameterTolerance
	Outer []diameterTolerance
}{
	0: { // допуски радиальных подшипников (кроме конических роликоподшипников)
		"PN": {
			[]diameterTolerance{
				{2 /*2.5*/, 10, 0, -8},
				{10, 18, 0, -8},
				{18, 30, 0, -10},
				{30, 50, 0, -12},
				{50, 80, 0, -15},
				{80, 120, 0, -20},
				{120, 180, 0, -25},
				{180, 250, 0, -30},
				{250, 315, 0, -35},
				{315, 400, 0, -40},
				{400, 500, 0, -45},
				{500, 630, 0, -50},
				{630, 800, 0, -75},
				{800, 1000, 0, -100},
				{1000, 1250, 0, -125},
				{1250, 1600, 0, -160},
				{1600, 2000, 0, -200},
			},
			[]diameterTolerance{
				{6, 18, 0, -8},
				{18, 30, 0, -9},
				{30, 50, 0, -11},
				{50, 80, 0, -13},
				{80, 120, 0, -15},
				{120, 150, 0, -18},
				{150, 180, 0, -25},
				{180, 250, 0, -30},
				{250, 315, 0, -35},
				{315, 400, 0, -40},
				{400, 500, 0, -45},
				{500, 630, 0, -50},
				{630, 800, 0, -75},
				{800, 1000, 0, -100},
				{1000, 1250, 0, -125},
				{1250, 1600, 0, -160},
				{1600, 2000, 0, -200},
				{2000, 2500, 0, -250},
			},
		},
		"P6": {
			[]diameterTolerance{
				{2 /*2.5*/, 10, 0, -7},
				{10, 18, 0, -7},
				{18, 30, 0, -8},
				{30, 50, 0, -10},
				{50, 80, 0, -12},
				{80, 120, 0, -15},
				{120, 180, 0, -18},
				{180, 250, 0, -22},
				{250, 315, 0, -25},
				{315, 400, 0, -30},
				{400, 500, 0, -35},
				{500, 630, 0, -40},
				{630, 800, 0, -50},
				{800, 1000, 0, -65},
				{1000, 1250, 0, -80},
				{1250, 1600, 0, -100},
				{1600, 2000, 0, -130},
			},
			[]diameterTolerance{
				{6, 18, 0, -7},
				{18, 30, 0, -8},
				{30, 50, 0, -9},
				{50, 80, 0, -11},
				{80, 120, 0, -13},
				{120, 150, 0, -15},
				{150, 180, 0, -18},
				{180, 250, 0, -20},
				{250, 315, 0, -25},
				{315, 400, 0, -28},
				{400, 500, 0, -33},
				{500, 630, 0, -38},
				{630, 800, 0, -45},
				{800, 1000, 0, -60},
				{1000, 1250, 0, -80},
				{1250, 1600, 0, -100},
				{1600, 2000, 0, -140},
				{2000, 2500, 0, -180},
			},
		},
		"P5": {
			[]diameterTolerance{
				{2 /*2.5*/, 10, 0, -5},
				{10, 18, 0, -5},
				{18, 30, 0, -6},
				{30, 50, 0, -8},
				{50, 80, 0, -9},
				{80, 120, 0, -10},
				{120, 180, 0, -13},
				{180, 250, 0, -15},
				{250, 315, 0, -18},
				{315, 400, 0, -23},
				{400, 500, 0, -27},
				{500, 630, 0, -33},
				{630, 800, 0, -40},
			},
			[]diameterTolerance{
				{6, 18, 0, -5},
				{18, 30, 0, -6},
				{30, 50, 0, -7},
				{50, 80, 0, -9},
				{80, 120, 0, -10},
				{120, 150, 0, -11},
				{150, 180, 0, -13},
				{180, 250, 0, -15},
				{250, 315, 0, -18},
				{315, 400, 0, -20},
				{400, 500, 0, -23},
				{500, 630, 0, -28},
				{630, 800, 0, -35},
				{800, 1000, 0, -40},
				{1000, 1250, 0, -50},
				{1250, 1600, 0, -65},
			},
		},
		"P4": {
			[]diameterTolerance{
				{2 /*2.5*/, 10, 0, -4},
				{10, 18, 0, -4},
				{18, 30, 0, -5},
				{30, 50, 0, -6},
				{50, 80, 0, -7},
				{80, 120, 0, -8},
				{120, 180, 0, -10},
				{180, 250, 0, -12},
				{250, 315, 0, -15},
				{315, 400, 0, -19},
				{400, 500, 0, -23},
				{500, 630, 0, -26},
				{630, 800, 0, -34},
			},
			[]diameterTolerance{
				{6, 18, 0, 0}, // nop
				{18, 30, 0, -5},
				{30, 50, 0, -6},
				{50, 80, 0, -7},
				{80, 120, 0, -8},
				{120, 150, 0, -9},
				{150, 180, 0, -10},
				{180, 250, 0, -11},
				{250, 315, 0, -13},
				{315, 400, 0, -15},
				{400, 500, 0, -20},
				{500, 630, 0, -25},
				{630, 800, 0, -28},
				{800, 1000, 0, -35},
				{1000, 1250, 0, -40},
				{1250, 1600, 0, -55},
			},
		},
		"P4S": { // шпиндельные
			[]diameterTolerance{
				{0, 10, 0, -4},
				{10, 18, 0, -4},
				{18, 30, 0, -5},
				{30, 50, 0, -6},
				{50, 80, 0, -7},
				{80, 120, 0, -8},
				{120, 150, 0, -10},
				{150, 180, 0, -10},
				{180, 250, 0, -12},
			},
			[]diameterTolerance{
				{18, 30, 0, -5},
				{30, 50, 0, -6},
				{50, 80, 0, -7},
				{80, 120, 0, -8},
				{120, 150, 0, -9},
				{150, 180, 0, -10},
				{180, 250, 0, -11},
				{250, 315, 0, -13},
				{315, 400, 0, -15},
			},
		},
		"SP": { // двухрядные роликоподшипники с цилиндрическими роликами
			[]diameterTolerance{
				{18, 30, 0, -6},
				{30, 50, 0, -8},
				{50, 80, 0, -9},
				{80, 120, 0, -10},
				{120, 180, 0, -13},
				{180, 250, 0, -15},
				{250, 315, 0, -18},
				{315, 400, 0, -23},
				{400, 500, 0, -27},
				{500, 630, 0, -30},
				{630, 800, 0, -40},
				{800, 1000, 0, -50},
				{1000, 1250, 0, -65},
			},
			[]diameterTolerance{
				{30, 50, 0, -7},
				{50, 80, 0, -9},
				{80, 120, 0, -10},
				{120, 150, 0, -11},
				{150, 180, 0, -13},
				{180, 250, 0, -15},
				{250, 315, 0, -18},
				{315, 400, 0, -20},
				{400, 500, 0, -23},
				{500, 630, 0, -28},
				{630, 800, 0, -35},
				{800, 1000, 0, -40},
				{1000, 1250, 0, -50},
				{1250, 1600, 0, -65},
			},
		},
		"UP": { // двухрядные роликоподшипники с цилиндрическими роликами
			[]diameterTolerance{
				{18, 30, 0, -5},
				{30, 50, 0, -6},
				{50, 80, 0, -7},
				{80, 120, 0, -8},
				{120, 180, 0, -10},
				{180, 250, 0, -12},
				{250, 315, 0, -15},
				{315, 400, 0, -19},
				{400, 500, 0, -23},
				{500, 630, 0, -26},
				{630, 800, 0, -34},
				{800, 1000, 0, -40},
				{1000, 1250, 0, -55},
			},
			[]diameterTolerance{
				{30, 50, 0, -5},
				{50, 80, 0, -6},
				{80, 120, 0, -7},
				{120, 150, 0, -8},
				{150, 180, 0, -9},
				{180, 250, 0, -10},
				{250, 315, 0, -12},
				{315, 400, 0, -14},
				{400, 500, 0, -17},
				{500, 630, 0, -20},
				{630, 800, 0, -25},
				{800, 1000, 0, -30},
				{1000, 1250, 0, -36},
				{1250, 1600, 0, -48},
			},
		},
	},
	1: { // допуски конических роликоподшипников
		"PN":  {},
		"P6X": {},
		"P5":  {},
		"P4":  {},
	},
	2: { // допуски упорных подшипников
		"PN": {},
		"P6": {},
		"P5": {},
		"P4": {},
		"SP": { // радиально-упорные шарикоподшипники, серии 2344 и 2347
			[]diameterTolerance{
				{0, 18, 0, 0}, // nop
				{18, 30, 0, -8},
				{30, 50, 0, -10},
				{50, 80, 0, -12},
				{80, 120, 0, -15},
				{120, 180, 0, -18},
				{180, 250, 0, -22},
				{250, 315, 0, -25},
				{315, 400, 0, -30},
				// >400 nop
			},
			[]diameterTolerance{
				{18, 30, 0, 0}, // nop
				{30, 50, 0, 0}, // nop
				{50, 80, -24, -43},
				{80, 120, -28, -50},
				{120, 180, -33, -58},
				{180, 250, -37, -66},
				{250, 315, -41, -73},
				{315, 400, -46, -82},
				{400, 500, -50, -90},
				{500, 630, -55, -99},
				// >630 nop
			},
		},
	},
}

type clearanceRange struct {
	Min int
	Max int
}

type clearance struct {
	DiameterLow  int
	DiameterHigh int
	Clearance    map[string]clearanceRange
}

/*
0 радиальный зазор шарикоподшипников с цилиндрическим отверстием
1 радиальный зазор сферических шарикоподшипников с цилиндрическим отверстием
2 -|- с коническим отверстием
3 осевой зазор двухрядных радиально-упорных шарикоподшипников серий 32, 32B, 33B
4 -|- серии 33DA
5 осевой зазор подшипников с четырёхточечным контактом
6 радиальный зазор однорядных и двухрядных роликоподшипников с цилиндрическими роликами с цилиндрическим отверстием
7 -|- с коническим отверстием
8 радиальный зазор сферических роликоподшипников с цилиндрическим отверстием
9 -|- с коническим отверстием
10 радиальный зазор сферических подшипников с бочкообразными роликами с цилиндрическим отверстием
11 -|- с коническим отверстием
*/
var gClearances = map[int][]clearance{
	0: { // радиальный зазор шарикоподшипников с цилиндрическим отверстием
		{2 /*2.5*/, 6, map[string]clearanceRange{
			"C2": {0, 7},
			"CN": {2, 13},
			"C3": {8, 23},
			"C4": {}, // nop
		}},
		{6, 10, map[string]clearanceRange{
			"C2": {0, 7},
			"CN": {2, 13},
			"C3": {8, 23},
			"C4": {14, 29},
		}},
		{10, 18, map[string]clearanceRange{
			"C2": {0, 9},
			"CN": {3, 18},
			"C3": {11, 25},
			"C4": {18, 33},
		}},
		{18, 24, map[string]clearanceRange{
			"C2": {0, 10},
			"CN": {5, 20},
			"C3": {13, 28},
			"C4": {20, 36},
		}},
		{24, 30, map[string]clearanceRange{
			"C2": {1, 11},
			"CN": {5, 20},
			"C3": {13, 28},
			"C4": {23, 41},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {1, 11},
			"CN": {6, 20},
			"C3": {15, 33},
			"C4": {28, 46},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {1, 11},
			"CN": {6, 23},
			"C3": {18, 36},
			"C4": {30, 51},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {1, 15},
			"CN": {8, 28},
			"C3": {23, 43},
			"C4": {38, 61},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {1, 15},
			"CN": {10, 30},
			"C3": {25, 51},
			"C4": {46, 71},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {1, 18},
			"CN": {12, 36},
			"C3": {30, 58},
			"C4": {53, 84},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {2, 20},
			"CN": {15, 41},
			"C3": {36, 66},
			"C4": {61, 97},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {2, 23},
			"CN": {18, 48},
			"C3": {41, 81},
			"C4": {71, 114},
		}},
		{140, 160, map[string]clearanceRange{
			"C2": {2, 23},
			"CN": {18, 53},
			"C3": {46, 91},
			"C4": {81, 130},
		}},
		{160, 180, map[string]clearanceRange{
			"C2": {2, 25},
			"CN": {20, 61},
			"C3": {53, 102},
			"C4": {91, 147},
		}},
		{180, 200, map[string]clearanceRange{
			"C2": {2, 30},
			"CN": {25, 71},
			"C3": {63, 117},
			"C4": {107, 163},
		}},
		{200, 225, map[string]clearanceRange{
			"C2": {4, 32},
			"CN": {28, 82},
			"C3": {73, 132},
			"C4": {120, 187},
		}},
		{225, 250, map[string]clearanceRange{
			"C2": {4, 36},
			"CN": {31, 92},
			"C3": {87, 152},
			"C4": {140, 217},
		}},
		{250, 280, map[string]clearanceRange{
			"C2": {4, 39},
			"CN": {36, 97},
			"C3": {97, 162},
			"C4": {152, 237},
		}},
		{280, 315, map[string]clearanceRange{
			"C2": {8, 45},
			"CN": {42, 110},
			"C3": {110, 180},
			"C4": {175, 260},
		}},
		{315, 355, map[string]clearanceRange{
			"C2": {8, 50},
			"CN": {50, 120},
			"C3": {120, 200},
			"C4": {200, 290},
		}},
		{355, 400, map[string]clearanceRange{
			"C2": {8, 60},
			"CN": {60, 140},
			"C3": {140, 230},
			"C4": {230, 330},
		}},
		{400, 450, map[string]clearanceRange{
			"C2": {10, 70},
			"CN": {70, 160},
			"C3": {160, 260},
			"C4": {260, 370},
		}},
		{450, 500, map[string]clearanceRange{
			"C2": {10, 80},
			"CN": {80, 180},
			"C3": {180, 290},
			"C4": {290, 410},
		}},
		{500, 560, map[string]clearanceRange{
			"C2": {20, 90},
			"CN": {90, 200},
			"C3": {200, 320},
			"C4": {320, 460},
		}},
		{560, 630, map[string]clearanceRange{
			"C2": {20, 100},
			"CN": {100, 220},
			"C3": {220, 350},
			"C4": {350, 510},
		}},
		{630, 710, map[string]clearanceRange{
			"C2": {30, 120},
			"CN": {120, 250},
			"C3": {250, 390},
			"C4": {390, 560},
		}},
		{710, 800, map[string]clearanceRange{
			"C2": {30, 130},
			"CN": {130, 280},
			"C3": {280, 440},
			"C4": {440, 620},
		}},
		{800, 900, map[string]clearanceRange{
			"C2": {30, 150},
			"CN": {150, 310},
			"C3": {310, 490},
			"C4": {490, 690},
		}},
		{900, 1000, map[string]clearanceRange{
			"C2": {40, 160},
			"CN": {160, 340},
			"C3": {340, 540},
			"C4": {540, 760},
		}},
		{1000, 1120, map[string]clearanceRange{
			"C2": {40, 170},
			"CN": {170, 370},
			"C3": {370, 590},
			"C4": {590, 840},
		}},
		{1120, 1250, map[string]clearanceRange{
			"C2": {40, 180},
			"CN": {180, 400},
			"C3": {400, 640},
			"C4": {640, 910},
		}},
		{1250, 1400, map[string]clearanceRange{
			"C2": {60, 210},
			"CN": {210, 440},
			"C3": {440, 700},
			"C4": {700, 1000},
		}},
		{1400, 1600, map[string]clearanceRange{
			"C2": {60, 230},
			"CN": {230, 480},
			"C3": {480, 770},
			"C4": {770, 1100},
		}},
	},
	1: { // радиальный зазор сферических шарикоподшипников с цилиндрическим отверстием
		{0, 6, map[string]clearanceRange{
			"C2": {1, 8},
			"CN": {5, 15},
			"C3": {10, 20},
			"C4": {15, 25},
		}},
		{6, 10, map[string]clearanceRange{
			"C2": {2, 9},
			"CN": {6, 17},
			"C3": {12, 25},
			"C4": {19, 33},
		}},
		{10, 14, map[string]clearanceRange{
			"C2": {2, 10},
			"CN": {6, 19},
			"C3": {13, 26},
			"C4": {21, 35},
		}},
		{14, 18, map[string]clearanceRange{
			"C2": {3, 12},
			"CN": {8, 21},
			"C3": {15, 28},
			"C4": {23, 37},
		}},
		{18, 24, map[string]clearanceRange{
			"C2": {4, 14},
			"CN": {10, 23},
			"C3": {17, 30},
			"C4": {25, 39},
		}},
		{24, 30, map[string]clearanceRange{
			"C2": {5, 16},
			"CN": {11, 24},
			"C3": {19, 35},
			"C4": {29, 46},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {6, 18},
			"CN": {13, 29},
			"C3": {23, 40},
			"C4": {34, 53},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {6, 19},
			"CN": {14, 31},
			"C3": {25, 44},
			"C4": {37, 57},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {7, 21},
			"CN": {16, 36},
			"C3": {30, 50},
			"C4": {45, 69},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {8, 24},
			"CN": {18, 40},
			"C3": {35, 60},
			"C4": {54, 83},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {9, 27},
			"CN": {22, 48},
			"C3": {42, 70},
			"C4": {64, 96},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {10, 31},
			"CN": {25, 56},
			"C3": {50, 83},
			"C4": {75, 114},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {10, 38},
			"CN": {30, 68},
			"C3": {60, 100},
			"C4": {90, 135},
		}},
		{140, 160, map[string]clearanceRange{
			"C2": {15, 44},
			"CN": {35, 80},
			"C3": {70, 120},
			"C4": {110, 161},
		}},
	},
	2: { // радиальный зазор сферических шарикоподшипников с коническим отверстием
		{0, 6, map[string]clearanceRange{}},   // nop
		{6, 10, map[string]clearanceRange{}},  // nop
		{10, 14, map[string]clearanceRange{}}, // nop
		{14, 18, map[string]clearanceRange{}}, // nop
		{18, 24, map[string]clearanceRange{
			"C2": {7, 17},
			"CN": {13, 26},
			"C3": {20, 33},
			"C4": {28, 42},
		}},
		{24, 30, map[string]clearanceRange{
			"C2": {9, 20},
			"CN": {15, 28},
			"C3": {23, 39},
			"C4": {33, 50},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {12, 24},
			"CN": {19, 35},
			"C3": {29, 46},
			"C4": {40, 59},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {14, 27},
			"CN": {22, 39},
			"C3": {33, 52},
			"C4": {45, 65},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {18, 32},
			"CN": {27, 47},
			"C3": {41, 61},
			"C4": {56, 80},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {23, 39},
			"CN": {35, 57},
			"C3": {50, 75},
			"C4": {69, 98},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {29, 47},
			"CN": {42, 68},
			"C3": {62, 90},
			"C4": {84, 116},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {35, 56},
			"CN": {50, 81},
			"C3": {75, 108},
			"C4": {100, 139},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {40, 68},
			"CN": {60, 98},
			"C3": {90, 130},
			"C4": {120, 165},
		}},
		{140, 160, map[string]clearanceRange{
			"C2": {45, 74},
			"CN": {65, 110},
			"C3": {100, 150},
			"C4": {140, 191},
		}},
	},
	3: { // осевой зазор двухрядных радиально-упорных шарикоподшипников серий 32, 32B, 33B
		{6, 10, map[string]clearanceRange{
			"C2": {1, 11},
			"CN": {5, 21},
			"C3": {12, 28},
			"C4": {25, 45},
		}},
		{10, 18, map[string]clearanceRange{
			"C2": {1, 12},
			"CN": {6, 23},
			"C3": {13, 31},
			"C4": {27, 47},
		}},
		{18, 24, map[string]clearanceRange{
			"C2": {2, 14},
			"CN": {7, 25},
			"C3": {16, 34},
			"C4": {28, 48},
		}},
		{24, 30, map[string]clearanceRange{
			"C2": {2, 15},
			"CN": {8, 27},
			"C3": {18, 37},
			"C4": {30, 50},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {2, 16},
			"CN": {9, 29},
			"C3": {21, 40},
			"C4": {33, 54},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {2, 18},
			"CN": {11, 33},
			"C3": {23, 44},
			"C4": {36, 58},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {3, 22},
			"CN": {13, 36},
			"C3": {26, 48},
			"C4": {40, 63},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {3, 24},
			"CN": {15, 40},
			"C3": {30, 54},
			"C4": {46, 71},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {3, 26},
			"CN": {18, 46},
			"C3": {35, 63},
			"C4": {55, 83},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {4, 30},
			"CN": {22, 53},
			"C3": {42, 73},
			"C4": {65, 96},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {4, 34},
			"CN": {25, 59},
			"C3": {48, 82},
			"C4": {74, 108},
		}},
	},
	4: { // осевой зазор двухрядных радиально-упорных шарикоподшипников серии 33DA
		{6, 10, map[string]clearanceRange{
			"C2": {5, 22},
			"CN": {11, 28},
			"C3": {20, 37},
		}},
		{10, 18, map[string]clearanceRange{
			"C2": {6, 24},
			"CN": {13, 31},
			"C3": {23, 41},
		}},
		{18, 24, map[string]clearanceRange{
			"C2": {7, 25},
			"CN": {14, 32},
			"C3": {24, 42},
		}},
		{24, 30, map[string]clearanceRange{
			"C2": {8, 27},
			"CN": {16, 35},
			"C3": {27, 46},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {9, 29},
			"CN": {18, 38},
			"C3": {30, 50},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {11, 33},
			"CN": {22, 44},
			"C3": {36, 58},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {13, 36},
			"CN": {25, 48},
			"C3": {40, 63},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {15, 40},
			"CN": {29, 54},
			"C3": {46, 71},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {18, 46},
			"CN": {35, 63},
			"C3": {55, 83},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {22, 53},
			"CN": {42, 73},
			"C3": {65, 96},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {25, 59},
			"CN": {48, 82},
			"C3": {74, 108},
		}},
	},
	5: { // осевой зазор подшипников с четырёхточечным контактом
		{0, 18, map[string]clearanceRange{
			"C2": {20, 60},
			"CN": {50, 90},
			"C3": {80, 120},
		}},
		{18, 40, map[string]clearanceRange{
			"C2": {30, 70},
			"CN": {60, 110},
			"C3": {100, 150},
		}},
		{40, 60, map[string]clearanceRange{
			"C2": {40, 90},
			"CN": {80, 130},
			"C3": {120, 170},
		}},
		{60, 80, map[string]clearanceRange{
			"C2": {50, 100},
			"CN": {90, 140},
			"C3": {130, 180},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {60, 120},
			"CN": {100, 160},
			"C3": {140, 200},
		}},
		{100, 140, map[string]clearanceRange{
			"C2": {70, 140},
			"CN": {120, 180},
			"C3": {160, 220},
		}},
		{140, 180, map[string]clearanceRange{
			"C2": {80, 160},
			"CN": {140, 200},
			"C3": {180, 240},
		}},
		{180, 220, map[string]clearanceRange{
			"C2": {100, 180},
			"CN": {160, 220},
			"C3": {200, 260},
		}},
		{220, 260, map[string]clearanceRange{
			"C2": {120, 200},
			"CN": {180, 240},
			"C3": {220, 300},
		}},
		{260, 300, map[string]clearanceRange{
			"C2": {140, 220},
			"CN": {200, 280},
			"C3": {260, 340},
		}},
		{300, 355, map[string]clearanceRange{
			"C2": {160, 240},
			"CN": {220, 300},
			"C3": {280, 360},
		}},
		{355, 400, map[string]clearanceRange{
			"C2": {180, 270},
			"CN": {250, 330},
			"C3": {310, 390},
		}},
		{400, 450, map[string]clearanceRange{
			"C2": {200, 290},
			"CN": {270, 360},
			"C3": {340, 430},
		}},
		{450, 500, map[string]clearanceRange{
			"C2": {220, 310},
			"CN": {290, 390},
			"C3": {370, 470},
		}},
		{500, 560, map[string]clearanceRange{
			"C2": {240, 330},
			"CN": {310, 420},
			"C3": {400, 510},
		}},
		{560, 630, map[string]clearanceRange{
			"C2": {260, 360},
			"CN": {340, 450},
			"C3": {430, 550},
		}},
		{630, 710, map[string]clearanceRange{
			"C2": {280, 390},
			"CN": {370, 490},
			"C3": {470, 590},
		}},
		{710, 800, map[string]clearanceRange{
			"C2": {300, 420},
			"CN": {400, 540},
			"C3": {520, 660},
		}},
		{800, 900, map[string]clearanceRange{
			"C2": {330, 460},
			"CN": {440, 590},
			"C3": {570, 730},
		}},
		{900, 1000, map[string]clearanceRange{
			"C2": {360, 500},
			"CN": {480, 630},
			"C3": {620, 780},
		}},
	},
	6: { // радиальный зазор однорядных и двухрядных роликоподшипников с цилиндрическими роликами с цилиндрическим отверстием
		{0, 24, map[string]clearanceRange{
			"C1NA": {5, 15},
			"C2":   {0, 25},
			"CN":   {20, 45},
			"C3":   {35, 60},
			"C4":   {50, 75},
		}},
		{24, 30, map[string]clearanceRange{
			"C1NA": {5, 15},
			"C2":   {0, 25},
			"CN":   {20, 45},
			"C3":   {35, 60},
			"C4":   {50, 75},
		}},
		{30, 40, map[string]clearanceRange{
			"C1NA": {5, 15},
			"C2":   {5, 30},
			"CN":   {25, 50},
			"C3":   {45, 70},
			"C4":   {60, 85},
		}},
		{40, 50, map[string]clearanceRange{
			"C1NA": {5, 18},
			"C2":   {5, 35},
			"CN":   {30, 60},
			"C3":   {50, 80},
			"C4":   {70, 100},
		}},
		{50, 65, map[string]clearanceRange{
			"C1NA": {5, 20},
			"C2":   {10, 40},
			"CN":   {40, 70},
			"C3":   {60, 90},
			"C4":   {80, 110},
		}},
		{65, 80, map[string]clearanceRange{
			"C1NA": {10, 25},
			"C2":   {10, 45},
			"CN":   {40, 75},
			"C3":   {65, 100},
			"C4":   {90, 125},
		}},
		{80, 100, map[string]clearanceRange{
			"C1NA": {10, 30},
			"C2":   {15, 50},
			"CN":   {50, 85},
			"C3":   {75, 110},
			"C4":   {105, 140},
		}},
		{100, 120, map[string]clearanceRange{
			"C1NA": {10, 30},
			"C2":   {15, 55},
			"CN":   {50, 90},
			"C3":   {75, 125},
			"C4":   {125, 165},
		}},
		{120, 140, map[string]clearanceRange{
			"C1NA": {10, 35},
			"C2":   {15, 60},
			"CN":   {60, 105},
			"C3":   {100, 145},
			"C4":   {145, 190},
		}},
		{140, 160, map[string]clearanceRange{
			"C1NA": {10, 35},
			"C2":   {20, 70},
			"CN":   {70, 120},
			"C3":   {115, 165},
			"C4":   {165, 215},
		}},
		{160, 180, map[string]clearanceRange{
			"C1NA": {10, 40},
			"C2":   {25, 75},
			"CN":   {75, 125},
			"C3":   {120, 170},
			"C4":   {170, 220},
		}},
		{180, 200, map[string]clearanceRange{
			"C1NA": {15, 45},
			"C2":   {35, 90},
			"CN":   {90, 145},
			"C3":   {140, 195},
			"C4":   {195, 250},
		}},
		{200, 225, map[string]clearanceRange{
			"C1NA": {15, 50},
			"C2":   {45, 105},
			"CN":   {105, 165},
			"C3":   {160, 220},
			"C4":   {220, 280},
		}},
		{225, 250, map[string]clearanceRange{
			"C1NA": {15, 50},
			"C2":   {45, 110},
			"CN":   {110, 175},
			"C3":   {170, 235},
			"C4":   {235, 300},
		}},
		{250, 280, map[string]clearanceRange{
			"C1NA": {20, 55},
			"C2":   {55, 125},
			"CN":   {125, 195},
			"C3":   {190, 260},
			"C4":   {260, 330},
		}},
		{280, 315, map[string]clearanceRange{
			"C1NA": {20, 60},
			"C2":   {55, 130},
			"CN":   {130, 205},
			"C3":   {200, 275},
			"C4":   {275, 350},
		}},
		{315, 355, map[string]clearanceRange{
			"C1NA": {20, 65},
			"C2":   {65, 145},
			"CN":   {145, 225},
			"C3":   {225, 305},
			"C4":   {305, 385},
		}},
		{355, 400, map[string]clearanceRange{
			"C1NA": {25, 75},
			"C2":   {100, 190},
			"CN":   {190, 280},
			"C3":   {280, 370},
			"C4":   {370, 460},
		}},
		{400, 450, map[string]clearanceRange{
			"C1NA": {25, 85},
			"C2":   {110, 210},
			"CN":   {210, 310},
			"C3":   {310, 410},
			"C4":   {410, 510},
		}},
		{450, 500, map[string]clearanceRange{
			"C1NA": {25, 95},
			"C2":   {110, 220},
			"CN":   {220, 330},
			"C3":   {330, 440},
			"C4":   {440, 550},
		}},
		{500, 560, map[string]clearanceRange{
			"C1NA": {25, 100},
			"C2":   {120, 240},
			"CN":   {240, 360},
			"C3":   {360, 480},
			"C4":   {480, 600},
		}},
		{560, 630, map[string]clearanceRange{
			"C1NA": {30, 110},
			"C2":   {140, 260},
			"CN":   {260, 380},
			"C3":   {380, 500},
			"C4":   {500, 620},
		}},
		{630, 710, map[string]clearanceRange{
			"C1NA": {30, 130},
			"C2":   {145, 285},
			"CN":   {285, 425},
			"C3":   {425, 565},
			"C4":   {565, 705},
		}},
		{710, 800, map[string]clearanceRange{
			"C1NA": {35, 140},
			"C2":   {150, 310},
			"CN":   {310, 470},
			"C3":   {470, 630},
			"C4":   {630, 790},
		}},
		{800, 900, map[string]clearanceRange{
			"C1NA": {35, 160},
			"C2":   {180, 350},
			"CN":   {350, 520},
			"C3":   {520, 690},
			"C4":   {690, 860},
		}},
		{900, 1000, map[string]clearanceRange{
			"C1NA": {35, 180},
			"C2":   {200, 390},
			"CN":   {390, 580},
			"C3":   {580, 770},
			"C4":   {770, 960},
		}},
		{1000, 1120, map[string]clearanceRange{
			"C1NA": {50, 200},
			"C2":   {220, 430},
			"CN":   {430, 640},
			"C3":   {640, 850},
			"C4":   {850, 1060},
		}},
		{1120, 1250, map[string]clearanceRange{
			"C1NA": {60, 220},
			"C2":   {230, 470},
			"CN":   {470, 710},
			"C3":   {710, 950},
			"C4":   {950, 1190},
		}},
		{1250, 1400, map[string]clearanceRange{
			"C1NA": {60, 240},
			"C2":   {270, 530},
			"CN":   {530, 790},
			"C3":   {790, 1050},
			"C4":   {1050, 1310},
		}},
		{1400, 1600, map[string]clearanceRange{
			"C1NA": {70, 270},
			"C2":   {330, 610},
			"CN":   {610, 890},
			"C3":   {890, 1170},
			"C4":   {1170, 1450},
		}},
		{1600, 1800, map[string]clearanceRange{
			"C1NA": {80, 300},
			"C2":   {380, 700},
			"CN":   {700, 1020},
			"C3":   {1020, 1340},
			"C4":   {1340, 1660},
		}},
		{1800, 2000, map[string]clearanceRange{
			"C1NA": {100, 320},
			"C2":   {400, 760},
			"CN":   {760, 1120},
			"C3":   {1120, 1480},
			"C4":   {1480, 1840},
		}},
	},
	7: { // радиальный зазор однорядных и двухрядных роликоподшипников с цилиндрическими роликами с коническим отверстием
		{0, 24, map[string]clearanceRange{
			"C1NA": {10, 20},
			"C2":   {15, 40},
			"CN":   {30, 55},
			"C3":   {40, 65},
			"C4":   {50, 75},
		}},
		{24, 30, map[string]clearanceRange{
			"C1NA": {15, 25},
			"C2":   {20, 45},
			"CN":   {35, 60},
			"C3":   {45, 70},
			"C4":   {55, 80},
		}},
		{30, 40, map[string]clearanceRange{
			"C1NA": {15, 25},
			"C2":   {20, 45},
			"CN":   {40, 65},
			"C3":   {55, 80},
			"C4":   {70, 95},
		}},
		{40, 50, map[string]clearanceRange{
			"C1NA": {17, 30},
			"C2":   {25, 55},
			"CN":   {45, 75},
			"C3":   {60, 90},
			"C4":   {75, 105},
		}},
		{50, 65, map[string]clearanceRange{
			"C1NA": {20, 35},
			"C2":   {30, 60},
			"CN":   {50, 80},
			"C3":   {70, 100},
			"C4":   {90, 120},
		}},
		{65, 80, map[string]clearanceRange{
			"C1NA": {25, 40},
			"C2":   {35, 70},
			"CN":   {60, 95},
			"C3":   {85, 120},
			"C4":   {110, 145},
		}},
		{80, 100, map[string]clearanceRange{
			"C1NA": {35, 55},
			"C2":   {40, 75},
			"CN":   {70, 105},
			"C3":   {95, 130},
			"C4":   {120, 155},
		}},
		{100, 120, map[string]clearanceRange{
			"C1NA": {40, 60},
			"C2":   {50, 90},
			"CN":   {90, 130},
			"C3":   {115, 155},
			"C4":   {140, 180},
		}},
		{120, 140, map[string]clearanceRange{
			"C1NA": {45, 70},
			"C2":   {55, 100},
			"CN":   {100, 145},
			"C3":   {130, 175},
			"C4":   {160, 205},
		}},
		{140, 160, map[string]clearanceRange{
			"C1NA": {50, 75},
			"C2":   {60, 110},
			"CN":   {110, 160},
			"C3":   {145, 195},
			"C4":   {180, 230},
		}},
		{160, 180, map[string]clearanceRange{
			"C1NA": {55, 85},
			"C2":   {75, 125},
			"CN":   {125, 175},
			"C3":   {160, 210},
			"C4":   {195, 245},
		}},
		{180, 200, map[string]clearanceRange{
			"C1NA": {60, 90},
			"C2":   {85, 140},
			"CN":   {140, 195},
			"C3":   {180, 235},
			"C4":   {220, 275},
		}},
		{200, 225, map[string]clearanceRange{
			"C1NA": {60, 95},
			"C2":   {95, 155},
			"CN":   {155, 215},
			"C3":   {200, 260},
			"C4":   {245, 305},
		}},
		{225, 250, map[string]clearanceRange{
			"C1NA": {65, 100},
			"C2":   {105, 170},
			"CN":   {170, 235},
			"C3":   {220, 285},
			"C4":   {270, 335},
		}},
		{250, 280, map[string]clearanceRange{
			"C1NA": {75, 110},
			"C2":   {115, 185},
			"CN":   {185, 255},
			"C3":   {240, 310},
			"C4":   {295, 365},
		}},
		{280, 315, map[string]clearanceRange{
			"C1NA": {80, 120},
			"C2":   {130, 205},
			"CN":   {205, 280},
			"C3":   {265, 340},
			"C4":   {325, 400},
		}},
		{315, 355, map[string]clearanceRange{
			"C1NA": {90, 135},
			"C2":   {145, 225},
			"CN":   {225, 305},
			"C3":   {290, 370},
			"C4":   {355, 435},
		}},
		{355, 400, map[string]clearanceRange{
			"C1NA": {100, 150},
			"C2":   {165, 255},
			"CN":   {255, 345},
			"C3":   {330, 420},
			"C4":   {405, 495},
		}},
		{400, 450, map[string]clearanceRange{
			"C1NA": {110, 170},
			"C2":   {185, 285},
			"CN":   {285, 385},
			"C3":   {370, 470},
			"C4":   {455, 555},
		}},
		{450, 500, map[string]clearanceRange{
			"C1NA": {120, 190},
			"C2":   {205, 315},
			"CN":   {315, 425},
			"C3":   {410, 520},
			"C4":   {505, 615},
		}},
		{500, 560, map[string]clearanceRange{
			"C1NA": {130, 210},
			"C2":   {230, 350},
			"CN":   {350, 470},
			"C3":   {455, 575},
			"C4":   {560, 680},
		}},
		{560, 630, map[string]clearanceRange{
			"C1NA": {140, 230},
			"C2":   {260, 380},
			"CN":   {380, 500},
			"C3":   {500, 620},
			"C4":   {620, 740},
		}},
		{630, 710, map[string]clearanceRange{
			"C1NA": {160, 260},
			"C2":   {295, 435},
			"CN":   {435, 575},
			"C3":   {565, 705},
			"C4":   {695, 835},
		}},
		{710, 800, map[string]clearanceRange{
			"C1NA": {170, 290},
			"C2":   {325, 485},
			"CN":   {485, 645},
			"C3":   {630, 790},
			"C4":   {775, 935},
		}},
		{800, 900, map[string]clearanceRange{
			"C1NA": {190, 330},
			"C2":   {370, 540},
			"CN":   {540, 710},
			"C3":   {700, 870},
			"C4":   {860, 1030},
		}},
		{900, 1000, map[string]clearanceRange{
			"C1NA": {210, 360},
			"C2":   {410, 600},
			"CN":   {600, 790},
			"C3":   {780, 970},
			"C4":   {960, 1150},
		}},
		{1000, 1120, map[string]clearanceRange{
			"C1NA": {230, 400},
			"C2":   {455, 665},
			"CN":   {665, 875},
			"C3":   {865, 1075},
			"C4":   {1065, 1275},
		}},
		{1120, 1250, map[string]clearanceRange{
			"C1NA": {250, 440},
			"C2":   {490, 730},
			"CN":   {730, 970},
			"C3":   {960, 1200},
			"C4":   {1200, 1440},
		}},
		{1250, 1400, map[string]clearanceRange{
			"C1NA": {270, 460},
			"C2":   {550, 810},
			"CN":   {810, 1070},
			"C3":   {1070, 1330},
			"C4":   {1330, 1590},
		}},
		{1400, 1600, map[string]clearanceRange{
			"C1NA": {300, 500},
			"C2":   {640, 920},
			"CN":   {920, 1200},
			"C3":   {1200, 1480},
			"C4":   {1480, 1760},
		}},
		{1600, 1800, map[string]clearanceRange{
			"C1NA": {320, 530},
			"C2":   {700, 1020},
			"CN":   {1020, 1340},
			"C3":   {1340, 1660},
			"C4":   {1660, 1980},
		}},
		{1800, 2000, map[string]clearanceRange{
			"C1NA": {340, 560},
			"C2":   {760, 1120},
			"CN":   {1120, 1480},
			"C3":   {1480, 1840},
			"C4":   {1840, 2200},
		}},
	},
	8: { // радиальный зазор сферических роликоподшипников с цилиндрическим отверстием
		{18, 24, map[string]clearanceRange{
			"C2": {10, 20},
			"CN": {20, 35},
			"C3": {35, 45},
			"C4": {45, 60},
		}},
		{24, 30, map[string]clearanceRange{
			"C2": {15, 25},
			"CN": {25, 40},
			"C3": {40, 55},
			"C4": {55, 75},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {15, 30},
			"CN": {30, 45},
			"C3": {45, 60},
			"C4": {60, 80},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {20, 35},
			"CN": {35, 55},
			"C3": {55, 75},
			"C4": {75, 100},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {20, 40},
			"CN": {40, 65},
			"C3": {65, 90},
			"C4": {90, 120},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {30, 50},
			"CN": {50, 80},
			"C3": {80, 110},
			"C4": {110, 145},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {35, 60},
			"CN": {60, 100},
			"C3": {100, 135},
			"C4": {135, 180},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {40, 75},
			"CN": {75, 120},
			"C3": {120, 160},
			"C4": {160, 210},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {50, 95},
			"CN": {95, 145},
			"C3": {145, 190},
			"C4": {190, 240},
		}},
		{140, 160, map[string]clearanceRange{
			"C2": {60, 110},
			"CN": {110, 170},
			"C3": {170, 220},
			"C4": {220, 280},
		}},
		{160, 180, map[string]clearanceRange{
			"C2": {65, 120},
			"CN": {120, 180},
			"C3": {180, 240},
			"C4": {240, 310},
		}},
		{180, 200, map[string]clearanceRange{
			"C2": {70, 130},
			"CN": {130, 200},
			"C3": {200, 260},
			"C4": {260, 340},
		}},
		{200, 225, map[string]clearanceRange{
			"C2": {80, 140},
			"CN": {140, 220},
			"C3": {220, 290},
			"C4": {290, 380},
		}},
		{225, 250, map[string]clearanceRange{
			"C2": {90, 150},
			"CN": {150, 240},
			"C3": {240, 320},
			"C4": {320, 420},
		}},
		{250, 280, map[string]clearanceRange{
			"C2": {100, 170},
			"CN": {170, 260},
			"C3": {260, 350},
			"C4": {350, 460},
		}},
		{280, 315, map[string]clearanceRange{
			"C2": {110, 190},
			"CN": {190, 280},
			"C3": {280, 370},
			"C4": {370, 500},
		}},
		{315, 355, map[string]clearanceRange{
			"C2": {120, 200},
			"CN": {200, 310},
			"C3": {310, 410},
			"C4": {410, 550},
		}},
		{355, 400, map[string]clearanceRange{
			"C2": {130, 220},
			"CN": {220, 340},
			"C3": {340, 450},
			"C4": {450, 600},
		}},
		{400, 450, map[string]clearanceRange{
			"C2": {140, 240},
			"CN": {240, 370},
			"C3": {370, 500},
			"C4": {500, 660},
		}},
		{450, 500, map[string]clearanceRange{
			"C2": {140, 260},
			"CN": {260, 410},
			"C3": {410, 550},
			"C4": {550, 720},
		}},
		{500, 560, map[string]clearanceRange{
			"C2": {150, 280},
			"CN": {280, 440},
			"C3": {440, 600},
			"C4": {600, 780},
		}},
		{560, 630, map[string]clearanceRange{
			"C2": {170, 310},
			"CN": {310, 480},
			"C3": {480, 650},
			"C4": {650, 850},
		}},
		{630, 710, map[string]clearanceRange{
			"C2": {190, 350},
			"CN": {350, 530},
			"C3": {530, 700},
			"C4": {700, 920},
		}},
		{710, 800, map[string]clearanceRange{
			"C2": {210, 390},
			"CN": {390, 580},
			"C3": {580, 770},
			"C4": {770, 1010},
		}},
		{800, 900, map[string]clearanceRange{
			"C2": {230, 430},
			"CN": {430, 650},
			"C3": {650, 860},
			"C4": {860, 1120},
		}},
		{900, 1000, map[string]clearanceRange{
			"C2": {260, 480},
			"CN": {480, 710},
			"C3": {710, 930},
			"C4": {930, 1220},
		}},
		{1000, 1120, map[string]clearanceRange{
			"C2": {290, 530},
			"CN": {530, 770},
			"C3": {770, 1050},
			"C4": {1050, 1430},
		}},
		{1120, 1250, map[string]clearanceRange{
			"C2": {320, 580},
			"CN": {580, 840},
			"C3": {840, 1140},
			"C4": {1140, 1560},
		}},
		{1250, 1400, map[string]clearanceRange{
			"C2": {350, 630},
			"CN": {630, 910},
			"C3": {910, 1240},
			"C4": {1240, 1700},
		}},
		{1400, 1600, map[string]clearanceRange{
			"C2": {380, 700},
			"CN": {700, 1020},
			"C3": {1020, 1390},
			"C4": {1390, 1890},
		}},
	},
	9: { // радиальный зазор сферических роликоподшипников с коническим отверстием
		{18, 24, map[string]clearanceRange{
			"C2": {15, 25},
			"CN": {25, 35},
			"C3": {35, 45},
			"C4": {45, 60},
		}},
		{24, 30, map[string]clearanceRange{
			"C2": {20, 30},
			"CN": {30, 40},
			"C3": {40, 55},
			"C4": {55, 75},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {25, 35},
			"CN": {35, 50},
			"C3": {50, 65},
			"C4": {65, 85},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {30, 45},
			"CN": {45, 60},
			"C3": {60, 80},
			"C4": {80, 100},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {40, 55},
			"CN": {55, 75},
			"C3": {75, 95},
			"C4": {95, 120},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {50, 70},
			"CN": {70, 95},
			"C3": {95, 120},
			"C4": {120, 150},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {55, 80},
			"CN": {80, 110},
			"C3": {110, 140},
			"C4": {140, 180},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {65, 100},
			"CN": {100, 135},
			"C3": {135, 170},
			"C4": {170, 220},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {80, 120},
			"CN": {120, 160},
			"C3": {160, 200},
			"C4": {200, 260},
		}},
		{140, 160, map[string]clearanceRange{
			"C2": {90, 130},
			"CN": {130, 180},
			"C3": {180, 230},
			"C4": {230, 300},
		}},
		{160, 180, map[string]clearanceRange{
			"C2": {100, 140},
			"CN": {140, 200},
			"C3": {200, 260},
			"C4": {260, 340},
		}},
		{180, 200, map[string]clearanceRange{
			"C2": {110, 160},
			"CN": {160, 220},
			"C3": {220, 290},
			"C4": {290, 370},
		}},
		{200, 225, map[string]clearanceRange{
			"C2": {120, 180},
			"CN": {180, 250},
			"C3": {250, 320},
			"C4": {320, 410},
		}},
		{225, 250, map[string]clearanceRange{
			"C2": {140, 200},
			"CN": {200, 270},
			"C3": {270, 350},
			"C4": {350, 450},
		}},
		{250, 280, map[string]clearanceRange{
			"C2": {150, 220},
			"CN": {220, 300},
			"C3": {300, 390},
			"C4": {390, 490},
		}},
		{280, 315, map[string]clearanceRange{
			"C2": {170, 240},
			"CN": {240, 330},
			"C3": {330, 430},
			"C4": {430, 540},
		}},
		{315, 355, map[string]clearanceRange{
			"C2": {190, 270},
			"CN": {270, 360},
			"C3": {360, 470},
			"C4": {470, 590},
		}},
		{355, 400, map[string]clearanceRange{
			"C2": {210, 300},
			"CN": {300, 400},
			"C3": {400, 520},
			"C4": {520, 650},
		}},
		{400, 450, map[string]clearanceRange{
			"C2": {230, 330},
			"CN": {330, 440},
			"C3": {440, 570},
			"C4": {570, 720},
		}},
		{450, 500, map[string]clearanceRange{
			"C2": {260, 370},
			"CN": {370, 490},
			"C3": {490, 630},
			"C4": {630, 790},
		}},
		{500, 560, map[string]clearanceRange{
			"C2": {290, 410},
			"CN": {410, 540},
			"C3": {540, 680},
			"C4": {680, 870},
		}},
		{560, 630, map[string]clearanceRange{
			"C2": {320, 460},
			"CN": {460, 600},
			"C3": {600, 760},
			"C4": {760, 980},
		}},
		{630, 710, map[string]clearanceRange{
			"C2": {350, 510},
			"CN": {510, 670},
			"C3": {670, 850},
			"C4": {850, 1090},
		}},
		{710, 800, map[string]clearanceRange{
			"C2": {390, 570},
			"CN": {570, 750},
			"C3": {750, 960},
			"C4": {960, 1220},
		}},
		{800, 900, map[string]clearanceRange{
			"C2": {440, 640},
			"CN": {640, 840},
			"C3": {840, 1070},
			"C4": {1070, 1370},
		}},
		{900, 1000, map[string]clearanceRange{
			"C2": {490, 710},
			"CN": {710, 930},
			"C3": {930, 1190},
			"C4": {1190, 1520},
		}},
		{1000, 1120, map[string]clearanceRange{
			"C2": {540, 780},
			"CN": {780, 1020},
			"C3": {1020, 1300},
			"C4": {1300, 1650},
		}},
		{1120, 1250, map[string]clearanceRange{
			"C2": {600, 860},
			"CN": {860, 1120},
			"C3": {1120, 1420},
			"C4": {1420, 1800},
		}},
		{1250, 1400, map[string]clearanceRange{
			"C2": {660, 940},
			"CN": {940, 1220},
			"C3": {1220, 1550},
			"C4": {1550, 1960},
		}},
		{1400, 1600, map[string]clearanceRange{
			"C2": {740, 1060},
			"CN": {1060, 1380},
			"C3": {1380, 1750},
			"C4": {1750, 2200},
		}},
	},
	10: { // радиальный зазор сферических подшипников с бочкообразными роликами с цилиндрическим отверстием
		{0, 30, map[string]clearanceRange{
			"C2": {2, 9},
			"CN": {9, 17},
			"C3": {17, 28},
			"C4": {28, 40},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {3, 10},
			"CN": {10, 20},
			"C3": {20, 30},
			"C4": {30, 45},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {3, 13},
			"CN": {13, 23},
			"C3": {23, 35},
			"C4": {35, 50},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {4, 15},
			"CN": {15, 27},
			"C3": {27, 40},
			"C4": {40, 55},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {5, 20},
			"CN": {20, 35},
			"C3": {35, 55},
			"C4": {55, 75},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {7, 25},
			"CN": {25, 45},
			"C3": {45, 65},
			"C4": {65, 90},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {10, 30},
			"CN": {30, 50},
			"C3": {50, 70},
			"C4": {70, 95},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {15, 35},
			"CN": {35, 55},
			"C3": {55, 80},
			"C4": {80, 110},
		}},
		{140, 160, map[string]clearanceRange{
			"C2": {20, 40},
			"CN": {40, 65},
			"C3": {65, 95},
			"C4": {95, 125},
		}},
		{160, 180, map[string]clearanceRange{
			"C2": {25, 45},
			"CN": {45, 70},
			"C3": {70, 100},
			"C4": {100, 130},
		}},
		{180, 225, map[string]clearanceRange{
			"C2": {30, 50},
			"CN": {50, 75},
			"C3": {75, 105},
			"C4": {105, 135},
		}},
		{225, 250, map[string]clearanceRange{
			"C2": {35, 55},
			"CN": {55, 80},
			"C3": {80, 110},
			"C4": {110, 140},
		}},
		{250, 280, map[string]clearanceRange{
			"C2": {40, 60},
			"CN": {60, 85},
			"C3": {85, 115},
			"C4": {115, 145},
		}},
		{280, 315, map[string]clearanceRange{
			"C2": {40, 70},
			"CN": {70, 100},
			"C3": {100, 135},
			"C4": {135, 170},
		}},
		{315, 355, map[string]clearanceRange{
			"C2": {45, 75},
			"CN": {75, 105},
			"C3": {105, 140},
			"C4": {140, 175},
		}},
	},
	11: { // радиальный зазор сферических подшипников с бочкообразными роликами с коническим отверстием
		{0, 30, map[string]clearanceRange{
			"C2": {9, 17},
			"CN": {17, 28},
			"C3": {28, 40},
			"C4": {40, 55},
		}},
		{30, 40, map[string]clearanceRange{
			"C2": {10, 20},
			"CN": {20, 30},
			"C3": {30, 45},
			"C4": {45, 60},
		}},
		{40, 50, map[string]clearanceRange{
			"C2": {13, 23},
			"CN": {23, 35},
			"C3": {35, 50},
			"C4": {50, 65},
		}},
		{50, 65, map[string]clearanceRange{
			"C2": {15, 27},
			"CN": {27, 40},
			"C3": {40, 55},
			"C4": {55, 75},
		}},
		{65, 80, map[string]clearanceRange{
			"C2": {20, 35},
			"CN": {35, 55},
			"C3": {55, 75},
			"C4": {75, 95},
		}},
		{80, 100, map[string]clearanceRange{
			"C2": {25, 45},
			"CN": {45, 65},
			"C3": {65, 90},
			"C4": {90, 120},
		}},
		{100, 120, map[string]clearanceRange{
			"C2": {30, 50},
			"CN": {50, 70},
			"C3": {70, 95},
			"C4": {95, 125},
		}},
		{120, 140, map[string]clearanceRange{
			"C2": {35, 55},
			"CN": {55, 80},
			"C3": {80, 110},
			"C4": {110, 140},
		}},
		{140, 160, map[string]clearanceRange{
			"C2": {40, 65},
			"CN": {65, 95},
			"C3": {95, 125},
			"C4": {125, 155},
		}},
		{160, 180, map[string]clearanceRange{
			"C2": {45, 70},
			"CN": {70, 100},
			"C3": {100, 130},
			"C4": {130, 160},
		}},
		{180, 225, map[string]clearanceRange{
			"C2": {50, 75},
			"CN": {75, 105},
			"C3": {105, 135},
			"C4": {135, 165},
		}},
		{225, 250, map[string]clearanceRange{
			"C2": {55, 80},
			"CN": {80, 110},
			"C3": {110, 140},
			"C4": {140, 170},
		}},
		{250, 280, map[string]clearanceRange{
			"C2": {60, 85},
			"CN": {85, 115},
			"C3": {115, 145},
			"C4": {145, 175},
		}},
		{280, 315, map[string]clearanceRange{
			"C2": {70, 100},
			"CN": {100, 135},
			"C3": {135, 170},
			"C4": {170, 205},
		}},
		{315, 355, map[string]clearanceRange{
			"C2": {75, 105},
			"CN": {105, 140},
			"C3": {140, 175},
			"C4": {175, 210},
		}},
	},
}
