package query

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type LabelSelector labels.Requirement

func (l LabelSelector) String() string {
	r := labels.Requirement(l)

	switch r.Operator() {
	case selection.Equals:
		return fmt.Sprintf(`%s="%s"`, r.Key(), r.Values().List()[0])
	case selection.NotEquals:
		return fmt.Sprintf(`%s!="%s"`, r.Key(), r.Values().List()[0])
	case selection.In:
		values := strings.Join(r.Values().List(), "|")
		return fmt.Sprintf(`%s=~"^(%s)$"`, r.Key(), values)
	case selection.NotIn:
		values := strings.Join(r.Values().List(), "|")
		return fmt.Sprintf(`%s!~"^(%s)$"`, r.Key(), values)
	}

	return ""
}

type LabelSelectors map[string]LabelSelector

func NewLabelSelectors(requirements []labels.Requirement) LabelSelectors {
	selectors := make(map[string]LabelSelector)
	for _, r := range requirements {
		selectors[r.Key()] = LabelSelector(r)
	}

	return selectors
}

func (ls LabelSelectors) String() string {
	strs := make([]string, 0)
	for _, l := range ls {
		strs = append(strs, l.String())
	}
	return strings.Join(strs, ", ")
}
