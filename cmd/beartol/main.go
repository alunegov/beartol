package main

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"

	"github.com/alunegov/beartol"
)

func main() {
	var questions = []*survey.Question{
		{
			Name: "type",
			Prompt: &survey.Select{
				Message:  "Type",
				Options:  beartol.RollingBearingTypeNames,
				Default:  0,
				PageSize: len(beartol.RollingBearingTypeNames),
			},
		},
		{
			Name: "typeSpecial",
			Prompt: &survey.Select{
				Message:  "TypeSpecial",
				Options:  beartol.RollingBearingTypeSpecialNames,
				Default:  0,
				PageSize: len(beartol.RollingBearingTypeSpecialNames),
			},
		},
		{
			Name:     "innerDiameter",
			Prompt:   &survey.Input{Message: "Inner diameter d"},
			Validate: survey.Required,
		},
		{
			Name:     "outerDiameter",
			Prompt:   &survey.Input{Message: "Outer diameter D"},
			Validate: survey.Required,
		},
		{
			Name: "classTochn",
			Prompt: &survey.Select{
				Message:  "ClassTochn",
				Options:  beartol.ClassTochnNames,
				Default:  "PN",
				PageSize: len(beartol.ClassTochnNames),
			},
		},
		{
			Name: "clearanceGroup",
			Prompt: &survey.Select{
				Message:  "ClearanceGroup",
				Options:  beartol.ClearanceGroupNames,
				Default:  "CN",
				PageSize: len(beartol.ClearanceGroupNames),
			},
		},
	}

	rb := surveyRollingBearing{}

	if err := survey.Ask(questions, &rb); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := beartol.Validate(rb.RollingBearing); err != nil {
		fmt.Println(err.Error())
		return
	}

	var tolerancesInteractor beartol.TolerancesInteractor

	tolerancesInteractor = &beartol.FagTolerancesInteractor{}

	innerDiameterToleranceA, innerDiameterToleranceB, innerErr := tolerancesInteractor.GetInnerDiameterTolerance(rb.RollingBearing)
	outerDiameterToleranceA, outerDiameterToleranceB, outerErr := tolerancesInteractor.GetOuterDiameterTolerance(rb.RollingBearing)
	clearanceMin, clearanceMax, clearanceErr := tolerancesInteractor.GetClearance(rb.RollingBearing)

	fmt.Printf("inner: %d %d µm, err: %v\n", innerDiameterToleranceA, innerDiameterToleranceB, innerErr)
	fmt.Printf("outer: %d %d µm, err: %v\n", outerDiameterToleranceA, outerDiameterToleranceB, outerErr)
	fmt.Printf("clearance: %d-%d µm, err: %v\n", clearanceMin, clearanceMax, clearanceErr)
}

type surveyRollingBearing struct {
	beartol.RollingBearing
}

func (thiz *surveyRollingBearing) WriteAnswer(name string, value interface{}) error {
	var err error = nil
	switch name {
	case "type":
		thiz.Type = beartol.RollingBearingType(value.(survey.OptionAnswer).Index)
	case "typeSpecial":
		thiz.TypeSpecial = beartol.RollingBearingTypeSpecial(value.(survey.OptionAnswer).Index)
	case "innerDiameter":
		thiz.InnerDiameter, err = strconv.Atoi(value.(string))
	case "outerDiameter":
		thiz.OuterDiameter, err = strconv.Atoi(value.(string))
	case "classTochn":
		thiz.ClassTochn = value.(survey.OptionAnswer).Value
	case "clearanceGroup":
		thiz.ClearanceGroup = value.(survey.OptionAnswer).Value
	default:
		err = fmt.Errorf("unsupp surveyRollingBearing field %s", name)
	}
	return err
}
