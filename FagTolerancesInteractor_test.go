package beartol

import "testing"

func TestGetInnerDiameterTolerance(t *testing.T) {
	sut := &FagTolerancesInteractor{}

	rb := RollingBearing{RollingBearingType_BallRadial, RollingBearingTypeSpecial_None, 20, 40, "PN", "CN"}

	tolA, tolB, err := sut.GetInnerDiameterTolerance(rb)

	if err != nil {
		t.Fatalf("err")
	}
	if tolA != 0 {
		t.Errorf("tolA")
	}
	if tolB != -10 {
		t.Errorf("tolB")
	}
}

func TestGetOuterDiameterTolerance(t *testing.T) {
	sut := &FagTolerancesInteractor{}

	rb := RollingBearing{RollingBearingType_BallRadial, RollingBearingTypeSpecial_None, 20, 40, "PN", "CN"}

	tolA, tolB, err := sut.GetOuterDiameterTolerance(rb)

	if err != nil {
		t.Fatalf("err")
	}
	if tolA != 0 {
		t.Errorf("tolA")
	}
	if tolB != -11 {
		t.Errorf("tolB")
	}
}

func TestGetClearance(t *testing.T) {
	sut := &FagTolerancesInteractor{}

	rb := RollingBearing{RollingBearingType_BallRadial, RollingBearingTypeSpecial_None, 20, 40, "PN", "CN"}

	clearanceMin, clearanceMax, err := sut.GetClearance(rb)

	if err != nil {
		t.Fatalf("err")
	}
	if clearanceMin != 5 {
		t.Errorf("clearanceMin")
	}
	if clearanceMax != 20 {
		t.Errorf("clearanceMax")
	}
}

func TestDiameterTolerancesData(t *testing.T) {
	for bearingType, classesTochn := range gDiameterTolerances {
		for classTochn, tolerances := range classesTochn {
			// Inner
			// нижняя граница диаметра меньше верхней
			for i := 0; i < len(tolerances.Inner); i++ {
				if tolerances.Inner[i].DiameterLow >= tolerances.Inner[i].DiameterHigh {
					t.Errorf("Diameter on %d, %s, [%d]", bearingType, classTochn, i)
				}
			}

			// "увеличение" диаметра
			for i := 1; i < len(tolerances.Inner); i++ {
				// нижняя граница диаметра увеличивается
				if tolerances.Inner[i-1].DiameterLow >= tolerances.Inner[i].DiameterLow {
					t.Errorf("DiameterLow on %d, %s, [%d]", bearingType, classTochn, i)
				}
				// верхняя граница диаметра увеличивается
				if tolerances.Inner[i-1].DiameterHigh >= tolerances.Inner[i].DiameterHigh {
					t.Errorf("DiameterHigh on %d, %s, [%d]", bearingType, classTochn, i)
				}
				// отклонениеA увеличивается
				if tolerances.Inner[i-1].ToleranceA < tolerances.Inner[i].ToleranceA {
					t.Errorf("ToleranceA on %d, %s, [%d]", bearingType, classTochn, i)
				}
				// отклонениеB увеличивается
				if tolerances.Inner[i-1].ToleranceB < tolerances.Inner[i].ToleranceB {
					t.Errorf("ToleranceB on %d, %s, [%d]", bearingType, classTochn, i)
				}
			}

			// Outer
			for i := 0; i < len(tolerances.Outer); i++ {
				if tolerances.Outer[i].DiameterLow >= tolerances.Outer[i].DiameterHigh {
					t.Errorf("Diameter on %d, %s, [%d]", bearingType, classTochn, i)
				}
			}

			for i := 1; i < len(tolerances.Outer); i++ {
				if tolerances.Outer[i-1].DiameterLow >= tolerances.Outer[i].DiameterLow {
					t.Errorf("DiameterLow on %d, %s, [%d]", bearingType, classTochn, i)
				}
				if tolerances.Outer[i-1].DiameterHigh >= tolerances.Outer[i].DiameterHigh {
					t.Errorf("DiameterHigh on %d, %s, [%d]", bearingType, classTochn, i)
				}
				if tolerances.Outer[i-1].ToleranceA < tolerances.Outer[i].ToleranceA {
					t.Errorf("ToleranceA on %d, %s, [%d]", bearingType, classTochn, i)
				}
				if tolerances.Outer[i-1].ToleranceB < tolerances.Outer[i].ToleranceB {
					t.Errorf("ToleranceB on %d, %s, [%d]", bearingType, classTochn, i)
				}
			}
		}
	}
}

func TestClearancesData(t *testing.T) {
	allClearanceGroups := []string{"C1NA", "C2", "CN", "C3", "C4"}
	commonClearanceGroups := []string{"C2", "CN", "C3", "C4"}

	for bearingType, clearances := range gClearances {
		for i := 0; i < len(clearances); i++ {
			// нижняя граница диаметра меньше верхней
			if clearances[i].DiameterLow >= clearances[i].DiameterHigh {
				t.Errorf("Diameter on %d, [%d] %d-%d mm", bearingType, i, clearances[i].DiameterLow, clearances[i].DiameterHigh)
			}
			// нижняя граница зазора меньше верхней
			for clearanceGroup, clearanceRange := range clearances[i].Clearance {
				// skipping nops
				if (clearanceRange.Min == 0) && (clearanceRange.Max == 0) {
					continue
				}
				if clearanceRange.Min >= clearanceRange.Max {
					t.Errorf("clearanceRange on %d, [%d] %d-%d mm, %s", bearingType, i, clearances[i].DiameterLow, clearances[i].DiameterHigh, clearanceGroup)
				}
			}
			// зазор увеличивается с "увеличением" группы зазора
			for j := 1; j < len(commonClearanceGroups); j++ {
				// skipping nops
				if (clearances[i].Clearance[commonClearanceGroups[j]].Min == 0) && (clearances[i].Clearance[commonClearanceGroups[j]].Max == 0) {
					continue
				}
				if clearances[i].Clearance[commonClearanceGroups[j-1]].Min >= clearances[i].Clearance[commonClearanceGroups[j]].Min {
					t.Errorf("clearanceMinA on %d, [%d] %d-%d mm, %s", bearingType, i, clearances[i].DiameterLow, clearances[i].DiameterHigh, commonClearanceGroups[j])
				}
				if clearances[i].Clearance[commonClearanceGroups[j-1]].Max >= clearances[i].Clearance[commonClearanceGroups[j]].Max {
					t.Errorf("clearanceMaxA on %d, [%d] %d-%d mm, %s", bearingType, i, clearances[i].DiameterLow, clearances[i].DiameterHigh, commonClearanceGroups[j])
				}
			}
		}

		// "увеличение" диаметра
		for i := 1; i < len(clearances); i++ {
			// нижняя граница диаметра увеличивается
			if clearances[i-1].DiameterLow >= clearances[i].DiameterLow {
				t.Errorf("DiameterLow on %d, [%d] %d mm", bearingType, i, clearances[i].DiameterLow)
			}
			// верхняя граница диаметра увеличивается
			if clearances[i-1].DiameterHigh >= clearances[i].DiameterHigh {
				t.Errorf("DiameterHigh on %d, [%d] %d mm", bearingType, i, clearances[i].DiameterHigh)
			}
			// зазор увеличивается
			for j := 0; j < len(allClearanceGroups); j++ {
				// skipping nops
				if (clearances[i].Clearance[allClearanceGroups[j]].Min == 0) && (clearances[i].Clearance[allClearanceGroups[j]].Max == 0) {
					continue
				}
				if clearances[i-1].Clearance[allClearanceGroups[j]].Min > clearances[i].Clearance[allClearanceGroups[j]].Min {
					t.Errorf("clearanceMinB on %d, [%d] %d-%d mm, %s", bearingType, i, clearances[i].DiameterLow, clearances[i].DiameterHigh, allClearanceGroups[j])
				}
				if clearances[i-1].Clearance[allClearanceGroups[j]].Max > clearances[i].Clearance[allClearanceGroups[j]].Max {
					t.Errorf("clearanceMaxB on %d, [%d] %d-%d mm, %s", bearingType, i, clearances[i].DiameterLow, clearances[i].DiameterHigh, allClearanceGroups[j])
				}
			}
		}
	}
}
