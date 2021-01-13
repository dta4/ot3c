package data

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

//TerminationPlanItems is the plan for controlled termination
var terminationPlanItems []TerminationItem
var ActiveTerminationPlan TerminationPlan

type TerminationPlan struct {
	CreatedAt time.Time
	Items     []TerminationItem
}

//TerminationItem is an item on the Termination list
type TerminationItem struct {
	TVR   VirtualResource
	TDate time.Time
	//TSave is the amount of money saved by this termination
	TSave float64
	//TASAP determines if this termination is a ASAP Termination
	TASAP bool
}

//terminationPlanFileItem
type terminationFileItem struct {
	VRID  string
	TDate time.Time
	TSave float64
	TASAP bool
}

//terminationPlanFile is the root struct of the Termination Plan file
type terminationPlanFile struct {
	CreatedAt time.Time
	Items     []terminationFileItem
}

//CurrentTerminationPlanToString renders the current TerminationPlan to string
func CurrentTerminationPlanToString() string {
	a := strings.Builder{}

	var totalSaved float64 = 0
	for _, ti := range ActiveTerminationPlan.Items {
		a.WriteString(fmt.Sprintf("ID: %v\n", ti.TVR.GetID()))
		a.WriteString(fmt.Sprintf("Prio: %v\n", GetPrio(ti.TVR)))
		a.WriteString(fmt.Sprintf("Termination Date: %v\n", ti.TDate.String()))
		a.WriteString(fmt.Sprintf("Costs saved: %v €\n", ti.TSave))
		a.WriteString(fmt.Sprintf("Termination ASAP: %v\n", ti.TASAP))
		a.WriteString(fmt.Sprintln("-----------------------------------------------"))
		totalSaved = totalSaved + ti.TSave
	}
	a.WriteString(fmt.Sprintf("Savings from Termination Plan: %v €", totalSaved))
	return a.String()
}

//RenderTerminationPlanToYAMLFile encodes the TerminationPlan to a String writable to a YAML file
func RenderTerminationPlanToYAMLFile() (string, error) {
	var items []terminationFileItem = make([]terminationFileItem, 0)

	for _, tvr := range ActiveTerminationPlan.Items {
		tfi := terminationFileItem{
			VRID:  tvr.TVR.GetID(),
			TDate: tvr.TDate,
			TSave: tvr.TSave,
			TASAP: tvr.TASAP,
		}
		items = append(items, tfi)
	}
	var filePlan terminationPlanFile = terminationPlanFile{
		Items:     items,
		CreatedAt: ActiveTerminationPlan.CreatedAt,
	}
	buffer, err := yaml.Marshal(filePlan)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func AddToTerminationPlan(t TerminationItem) {
	terminationPlanItems = append(terminationPlanItems, t)
}
func ClearTerminationPlan() {
	terminationPlanItems = make([]TerminationItem, 0)
}
func SetTerminationItems(t []TerminationItem) {
	terminationPlanItems = t
}

//UpdateTerminationPlan creates new TerminationPlan based upon current Termination Items
func UpdateTerminationPlan() {
	ActiveTerminationPlan = TerminationPlan{
		CreatedAt: time.Now(),
		Items:     terminationPlanItems,
	}
}

func ParseTerminationPlanFromString(str string) error {
	var tpf terminationPlanFile
	err := yaml.Unmarshal([]byte(str), &tpf)
	if err != nil {
		return err
	}
	ClearTerminationPlan()
	for _, key := range tpf.Items {
		vr := FindVRByID(key.VRID)
		if vr != nil {
			item := TerminationItem{
				TVR:   vr,
				TASAP: key.TASAP,
				TSave: key.TSave,
				TDate: key.TDate,
			}
			AddToTerminationPlan(item)
		}
	}
	UpdateTerminationPlan()
	ActiveTerminationPlan.CreatedAt = tpf.CreatedAt
	return nil
}
