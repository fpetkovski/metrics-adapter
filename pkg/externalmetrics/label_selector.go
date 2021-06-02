package externalmetrics

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

type LabelSelectorList []LabelSelector

func (ls LabelSelectorList) String() string {
	strs := make([]string, len(ls))
	for i, l := range ls {
		strs[i] = l.String()
	}
	return strings.Join(strs, ", ")
}
