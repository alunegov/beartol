package beartol

import "fmt"

// RollingBearingType обозначает тип подшипника качения
type RollingBearingType int

const (
	// радиальные шарикоподшипники 160, 161, 60, S60, 618, 62, S62, 622, 623, 63, S63, 64
	RollingBearingType_BallRadial RollingBearingType = iota
	// радиально-упорные шарикоподшипники 72B, 73B, 32B, 33B, 32, 33, 33DA
	RollingBearingType_BallRadialUpor
	// шпиндельные подшипники B70, B719, B72, HCS70, HCS719, HSS70, HSS719
	RollingBearingType_Shpindel
	// подшипники с четырёхточечным контактом QJ2, QJ3
	RollingBearingType_4ptContact
	// сферические шарикоподшипники с цилиндрическим отверстием 12, 13, 22, 23, 112
	RollingBearingType_BallSphere_Cylynder
	// сферические шарикоподшипники с коническим отверстием 12K, 13K, 22K, 23K
	RollingBearingType_BallSphere_Kone
	// цилиндрические роликоподшипники NU10, 19, 2, 22, 23, 3, NJ2, 22, 23, 3, NUP2, 22, 23, 3, N2, 3, NN30ASK
	RollingBearingType_RollerCylynder
	// конические роликоподшипники 302, 303, 313, 320, 322, 323, 329, 330, 331, 332, T......
	RollingBearingType_RollerKone
	// сферические роликоподшипники с цилиндрическим отверстием 213, 222, 223, 230, 231, 232, 233, 239, 240, 241
	RollingBearingType_RollerSphere_Cylynder
	// сферические роликоподшипники с коническим отверстием 213K, 222K, 223K, 230K, 231K, 232K, 239K, 240K30, 241K30
	RollingBearingType_RollerSphere_Kone
	// сферические роликоподшипники с бочкообразными роликами с цилиндрическим отверстием 202, 203
	RollingBearingType_RollerSphere_Bochka_Cylynder
	// сферические роликоподшипники с бочкообразными роликами с коническим отверстием 202K, 203K
	RollingBearingType_RollerSphere_Bochka_Kone
	// упорные шарикоподшипники 511, 512, 513, 514, 532, 533, 522, 523, 542, 543
	RollingBearingType_BallUpor
	// радиально-упорные шарикоподшипники 7602, 7603, 2344, 2347
	RollingBearingType_BallRadialUpor2
	// упорные цилиндрические роликоподшипники 811, 812
	RollingBearingType_RollerCylynderUpor
	// упорные сферические роликоподшипники 292E, 293E, 294E
	RollingBearingType_RollerSphereUpor
)

var RollingBearingTypeNames = []string{
	"радиальные шарикоподшипники 160, 161, 60, S60, 618, 62, S62, 622, 623, 63, S63, 64",
	"радиально-упорные шарикоподшипники 72B, 73B, 32B, 33B, 32, 33, 33DA",
	"шпиндельные подшипники B70, B719, B72, HCS70, HCS719, HSS70, HSS719",
	"подшипники с четырёхточечным контактом QJ2, QJ3",
	"сферические шарикоподшипники с цилиндрическим отверстием 12, 13, 22, 23, 112",
	"сферические шарикоподшипники с коническим отверстием 12K, 13K, 22K, 23K",
	"цилиндрические роликоподшипники NU10, 19, 2, 22, 23, 3, NJ2, 22, 23, 3, NUP2, 22, 23, 3, N2, 3, NN30ASK",
	"конические роликоподшипники 302, 303, 313, 320, 322, 323, 329, 330, 331, 332, T......",
	"сферические роликоподшипники с цилиндрическим отверстием 213, 222, 223, 230, 231, 232, 233, 239, 240, 241",
	"сферические роликоподшипники с коническим отверстием 213K, 222K, 223K, 230K, 231K, 232K, 239K, 240K30, 241K30",
	"сферические роликоподшипники с бочкообразными роликами с цилиндрическим отверстием 202, 203",
	"сферические роликоподшипники с бочкообразными роликами с коническим отверстием 202K, 203K",
	"упорные шарикоподшипники 511, 512, 513, 514, 532, 533, 522, 523, 542, 543",
	"радиально-упорные шарикоподшипники 7602, 7603, 2344, 2347",
	"упорные цилиндрические роликоподшипники 811, 812",
	"упорные сферические роликоподшипники 292E, 293E, 294E",
}

// RollingBearingTypeSpecial обозначает особенности типа подшипника качения
type RollingBearingTypeSpecial int

const (
	// RollingBearingTypeSpecial_None обозначает подшипник без особенностей
	RollingBearingTypeSpecial_None RollingBearingTypeSpecial = iota
	// RollingBearingTypeSpecial_72B_73B обозначает подшипник 72B_73B
	RollingBearingTypeSpecial_72B_73B
	// RollingBearingTypeSpecial_33DA обозначает подшипник 33DA
	RollingBearingTypeSpecial_33DA
	// RollingBearingTypeSpecial_NN30ASK обозначает подшипник NN30ASK
	RollingBearingTypeSpecial_NN30ASK
)

var RollingBearingTypeSpecialNames = []string{
	"подшипник без особенностей",
	"подшипник 72B_73B",
	"подшипник 33DA",
	"подшипник NN30ASK",
}

var ClassTochnNames = []string{"PN", "P6", "P6X", "P5", "P4", "P4S", "SP", "UP"}

var ClearanceGroupNames = []string{"C1NA", "C2", "CN", "C3", "C4"}

// RollingBearing описывает подшипник качения
type RollingBearing struct {
	// Тип подшипника
	Type RollingBearingType
	// Особенности типа подшипника
	TypeSpecial RollingBearingTypeSpecial
	// Внутренний диаметр, мм
	InnerDiameter int
	// Наружный диаметр, мм
	OuterDiameter int
	// Класс точности: PN, P6, P6X, P5, P4, P4S, SP, UP
	ClassTochn string
	// Группа зазора: C1NA, C2, CN, C3, C4
	ClearanceGroup string
}

func Validate(rb RollingBearing) error {
	switch rb.TypeSpecial {
	case RollingBearingTypeSpecial_72B_73B:
		fallthrough
	case RollingBearingTypeSpecial_33DA:
		if rb.Type != RollingBearingType_BallRadialUpor {
			return fmt.Errorf("for special types 72B_73B or 33DA allowed type is BallRadialUpor")
		}
	case RollingBearingTypeSpecial_NN30ASK:
		if rb.Type != RollingBearingType_RollerCylynder {
			return fmt.Errorf("for special type NN30ASK allowed type is RollerCylynder")
		}
	}

	if rb.InnerDiameter >= rb.OuterDiameter {
		return fmt.Errorf("inner diameter %d is greater (or equal) whan outer %d", rb.InnerDiameter, rb.OuterDiameter)
	}

	return nil
}
