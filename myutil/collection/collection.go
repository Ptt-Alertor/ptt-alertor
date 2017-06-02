package collection

import "reflect"

func Index(c interface{}, elem interface{}) int {
	arrV := reflect.ValueOf(c)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if reflect.DeepEqual(arrV.Index(i).Interface(), elem) {
				return i
			}
		}
	}
	return -1
}

func In(c interface{}, elem interface{}) bool {
	if Index(c, elem) == -1 {
		return false
	}
	return true
}
