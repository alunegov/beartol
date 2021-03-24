package beartol

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

type RollingBearingTypeSpecial int

const (
	RollingBearingTypeSpecial_None RollingBearingTypeSpecial = iota
	RollingBearingTypeSpecial_72B_73B
	RollingBearingTypeSpecial_33DA
	RollingBearingTypeSpecial_NN30ASK
)

type RollingBearing struct {
	Type           RollingBearingType
	TypeSpecial    RollingBearingTypeSpecial
	InnerDiameter  int
	OuterDiameter  int
	ClassTochn     string
	ClearanceGroup string
}
