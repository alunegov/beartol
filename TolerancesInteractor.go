package beartol

type TolerancesInteractor interface {
	GetInnerDiameterTolerance(rb RollingBearing) (int, int, error)
	GetOuterDiameterTolerance(rb RollingBearing) (int, int, error)
	GetClearance(rb RollingBearing) (int, int, error)
}
