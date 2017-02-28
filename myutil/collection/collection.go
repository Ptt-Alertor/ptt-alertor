package collection

import "reflect"

func Index(c []interface{}, elem interface{}) int {
	for i, v := range c {
		if reflect.DeepEqual(v, elem) {
			return i
		}
	}
	return -1
}

func In(c []interface{}, elem interface{}) bool {
	if Index(c, elem) == -1 {
		return false
	}
	return true
}
