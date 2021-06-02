package externalmetrics

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type labelSelector struct {
	r labels.Requirement
}

func (l labelSelector) String() string {
	switch l.r.Operator() {
	case selection.Equals:
		return fmt.Sprintf(`%s='%s'`, l.r.Key(), l.r.Values().List()[0])
	case selection.NotEquals:
		return fmt.Sprintf(`%s!='%s'`, l.r.Key(), l.r.Values().List()[0])
	case selection.In:
		values := strings.Join(l.r.Values().List(), "|")
		return fmt.Sprintf(`%s=~'^(%s)$'`, l.r.Key(), values)
	case selection.NotIn:
		values := strings.Join(l.r.Values().List(), "|")
		return fmt.Sprintf(`%s!~'^(%s)$'`, l.r.Key(), values)
	}

	return ""
}
