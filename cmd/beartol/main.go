package main

import (
	"fmt"

	"github.com/alunegov/beartol"
)

func main() {
	var tolerancesInteractor beartol.TolerancesInteractor

	tolerancesInteractor = &beartol.FagTolerancesInteractor{}

	rb := beartol.RollingBearing{
		Type:           beartol.RollingBearingType_BallRadial,
		TypeSpecial:    beartol.RollingBearingTypeSpecial_None,
		InnerDiameter:  20,
		OuterDiameter:  40,
		ClassTochn:     "PN",
		ClearanceGroup: "CN",
	}

	innerDiameterToleranceA, innerDiameterToleranceB, innerErr := tolerancesInteractor.GetInnerDiameterTolerance(rb)
	outerDiameterToleranceA, outerDiameterToleranceB, outerErr := tolerancesInteractor.GetOuterDiameterTolerance(rb)
	clearanceMin, clearanceMax, clearanceErr := tolerancesInteractor.GetClearance(rb)

	fmt.Printf("inner: %d %d мкм %v\n", innerDiameterToleranceA, innerDiameterToleranceB, innerErr)
	fmt.Printf("outer: %d %d мкм %v\n", outerDiameterToleranceA, outerDiameterToleranceB, outerErr)
	fmt.Printf("clearance: %d-%d мкм %v\n", clearanceMin, clearanceMax, clearanceErr)
}
