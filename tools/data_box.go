package tools

import (
	"fmt"
	"reflect"
)

///////////////////////////////////////////////////////////////////////////////////////
/*
// 为了高效、便捷的为Slice中的对象结构中的成员进行赋值操作（待赋值的数据可能在数据库中），从而实现该方法。
// KeyField在slice中可以有重复
//
// 如有结构：

type T struct {
	KeyId   int
	TName   string
	V       V
	PV      *V
	VSlice  []V
	PVSlice []*V
}

生成了对应的Slice：objs := []T{T{...},T{...},T{...},....}
需要填充其中成员的V 或 PV 或 VSlice 或 PVSlice域的值，使用如下结构的对象：
type V struct {
	Id    int
	VName string
}

使用DataBox后的方式是：

err := NewDataBox(objs).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
	retMap := make(map[int]*V)

	// 一次性从数据库中得到所有的Mobjs
	KeyIds := keywords.([]int)
	if len(KeyIds) > 0 {
		VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
		for i, v := range VSlice {
			retMap[v.Id] = &VSlice[i]
		}
	}

	return retMap
}).SaveToField("V")

*/
///////////////////////////////////////////////////////////////////////////////////////

type DataBox struct {
	// 主数据
	sliceData interface{}

	// sliceData的reflect.Value
	sliceDataValue reflect.Value

	// ptr type of sliceData中单个对象的数据类型
	elemPtrType reflect.Type

	// 主数据Map，key(keyField)指向相同keyField的elem的引用对象组成的slice
	key2ElemSlice reflect.Value

	// JoinByMap的到的Map，主key(keyField)指向返回的数据对象
	retMap reflect.Value

	// 内部错误传递，不用关心该成员
	err error
}

// 用户需要实现的函数：通过key列表来获取对象，用于JoinByMap
// key列表时data里面通过keyfield抽取的key列表
// 第一个返回的map的key为传入的key；map里面的value可以是查到的单个对象，也可以是一个可以对应的对象slice
// 与FetchObjsByKeysRetSliceFunc相比，性能更高，且子结构可以不带父对象的ID，缺点是需要多写点代码
type FetchObjsByKeysRetMapFunc func(keyList interface{}) (key2ObjsMap interface{}) // map[key][]obj  or  map[key]obj  or map[key][]*obj  or  map[key]*obj

// 功能：生成NewDataBox对象，以便向data(必须是slice)中的成员(slice中的成员)添加指定数据，数据由FetchObjsByKeysRetMapFunc或者
// FetchObjsByKeysRetSliceFunc返回slice中的对象和map中的value可以是结构或者是结构的指针
// keyFieldName指明后续获取数据的key，后续会收集所有的key，传到FetchObjs...函数，以便集中化获取数据，提升效率
// 通过SaveToField函数指明的fieldname，得到数据填充的位置
// 其它特点：
// 1. 可以回填对象或者对象Slice
// 2. 可以多重填充
// 3. 被填充对象为指针，填充对象为结构对象 这种情况会导致异常
func NewDataBox(data interface{}) *DataBox {
	databox := &DataBox{}
	databox.sliceData = data
	if data == nil {
		databox.err = fmt.Errorf("Empty Input Data")
		return databox
	}

	databox.sliceDataValue = reflect.ValueOf(data)
	if !databox.sliceDataValue.IsValid() || databox.sliceDataValue.IsNil() {
		databox.err = fmt.Errorf("Wrong Input Data Format")
		return databox
	}

	if databox.sliceDataValue.Kind() == reflect.Ptr {
		databox.sliceDataValue = databox.sliceDataValue.Elem()
		databox.sliceData = databox.sliceDataValue.Interface()
	}

	var elemType reflect.Type
	if databox.sliceDataValue.Kind() == reflect.Slice { // 如果是slice类型, 取出实体类型
		elemType = databox.sliceDataValue.Type().Elem()
	} else {
		databox.err = fmt.Errorf("Wrong Input Data's Type. Only Should Be Slice Or Map")
		return databox
	}

	// 确认真实类型
	if elemType.Kind() == reflect.Ptr {
		databox.elemPtrType = elemType
	} else {
		databox.elemPtrType = reflect.PtrTo(elemType)
	}

	if databox.elemPtrType.Elem().Kind() != reflect.Struct {
		databox.err = fmt.Errorf("Wrong Input Data's Element Type. Only Should Be Struct")
		return databox
	}

	return databox
}

func (d *DataBox) KeyField(keyFieldName string) *DataBox {
	if d.err != nil {
		return d
	}

	structField, ok := d.elemPtrType.Elem().FieldByName(keyFieldName)
	if !ok {
		d.err = fmt.Errorf("keyFieldName not in element")
		return d
	}

	d.key2ElemSlice = reflect.MakeMap(reflect.MapOf(structField.Type, reflect.SliceOf(d.elemPtrType)))

	//将原始数据, 整理进中间的比较字典
	for i := 0; i < d.sliceDataValue.Len(); i++ {
		item := d.sliceDataValue.Index(i)
		if item.Kind() == reflect.Ptr {
			d.storeKVIntoMap(item.Elem().FieldByIndex(structField.Index), item, &d.key2ElemSlice)
		} else {
			d.storeKVIntoMap(item.FieldByIndex(structField.Index), item.Addr(), &d.key2ElemSlice)
		}
	}

	return d
}

func (d *DataBox) JoinByMap(f FetchObjsByKeysRetMapFunc) *DataBox {
	if d.err != nil {
		return d
	}

	keyList := reflect.MakeSlice(reflect.SliceOf(d.key2ElemSlice.Type().Key()), 0, len(d.key2ElemSlice.MapKeys()))
	for _, v := range d.key2ElemSlice.MapKeys() {
		keyList = reflect.Append(keyList, v)
	}

	// 获取关联数据
	retMap := f(keyList.Interface())
	if retMap == nil {
		return nil
	}

	d.retMap = reflect.ValueOf(retMap)
	return d
}

func (d *DataBox) SaveToField(fieldName string) error {
	if d.err != nil {
		return d.err
	}

	destField, exist := d.elemPtrType.Elem().FieldByName(fieldName)
	if !exist {
		return fmt.Errorf("%s not exsit in struct", fieldName)
	}

	// destField 与 用户函数返回的map的value类型必须相同 (兼容指针情况)
	var retTypeIsPtr bool = false
	retType := d.retMap.Type().Elem()
	if retType.Kind() == reflect.Ptr {
		retTypeIsPtr = true
		retType = d.retMap.Type().Elem().Elem()
	}

	var destTypeIsPtr bool = false
	destType := destField.Type
	if destType.Kind() == reflect.Ptr {
		destTypeIsPtr = true
		destType = destField.Type.Elem()
	}

	if retType.Kind() != destType.Kind() {
		// 如果不相等，那么返回的可能是slice，那做特殊处理
		if retType.Kind() == reflect.Slice {
			if retType.Elem().Kind() == reflect.Ptr {
				if retType.Elem().Elem().Kind() == destType.Kind() {
					return d.saveToFieldWithSlice(&destField, destTypeIsPtr)
				}
			} else if retType.Kind() == destType.Kind() {
				return d.saveToFieldWithSlice(&destField, destTypeIsPtr)
			}
		}

		if retType.Kind() == reflect.Ptr && retType.Elem().Kind() == reflect.Slice {
			if retType.Elem().Elem().Kind() == reflect.Ptr {
				if retType.Elem().Elem().Elem().Kind() == destType.Kind() {
					return d.saveToFieldWithSlice(&destField, destTypeIsPtr)
				}
			} else if retType.Elem().Kind() == destType.Kind() {
				return d.saveToFieldWithSlice(&destField, destTypeIsPtr)
			}
		}

		return fmt.Errorf("Type Error When SaveToField: retType: %s != desttype: %s", retType.String(), destType.String())
	}

	// 填充slice的时候，可能需要一个一个处理
	var isNeedConvertWhenElemTypeNotEquelInSlice = false
	if destType.Kind() == reflect.Slice && retType.Kind() == reflect.Slice {
		if destType.Elem().Kind() != retType.Elem().Kind() {
			isNeedConvertWhenElemTypeNotEquelInSlice = true
		}
	}

	// saveStructField可以是slice或者是struct或者是一个普通类型
	iter := d.retMap.MapRange()
	for iter.Next() {
		saveElemSlice := d.key2ElemSlice.MapIndex(iter.Key())
		if !saveElemSlice.IsValid() {
			continue
		}

		retObjValue := iter.Value()
		if isNeedConvertWhenElemTypeNotEquelInSlice {
			retObjValue = convertSlice(iter.Value(), destType)
		}

		// 把value放入saveElem slice中对象的fieldName处
		for i := 0; i < saveElemSlice.Len(); i++ {
			dest := saveElemSlice.Index(i).Elem().FieldByIndex(destField.Index)
			if destTypeIsPtr {
				if retTypeIsPtr {
					dest.Set(retObjValue)
				} else {
					dest.Set(retObjValue.Addr())
				}
			} else {
				if retTypeIsPtr {
					dest.Set(retObjValue.Elem())
				} else {
					dest.Set(retObjValue)
				}
			}
		}
	}

	return nil
}

// 把keyField指向的key作为map的key，obj in objSlice作为value，组成一个map，value为Ptr类型
//func Slice2MapByField(objSlice interface{}, keyField string) map[interface{}]interface{} {
//	retMap := map[interface{}]interface{}{}
//
//	typeOfObjSlice := reflect.TypeOf(objSlice)
//	if typeOfObjSlice.Kind() != reflect.Slice {
//		panic("Slice2MapByField's objSlice must be slice")
//	}
//
//	var ok bool
//	var structField reflect.StructField
//	if typeOfObjSlice.Elem().Kind() == reflect.Ptr {
//		structField, ok = typeOfObjSlice.Elem().Elem().FieldByName(keyField)
//	} else {
//		structField, ok = typeOfObjSlice.Elem().FieldByName(keyField)
//	}
//
//	if !ok {
//		return retMap
//	}
//
//	valueOfObjSlice := reflect.ValueOf(objSlice)
//	for i := 0; i < valueOfObjSlice.Len(); i++ {
//		if typeOfObjSlice.Elem().Kind() == reflect.Ptr {
//			key := valueOfObjSlice.Index(i).Elem().Field(structField.Index[0]).Interface()
//			retMap[key] = valueOfObjSlice.Index(i).Interface()
//		} else {
//			key := valueOfObjSlice.Index(i).Field(structField.Index[0]).Interface()
//			retMap[key] = valueOfObjSlice.Index(i).Addr().Interface()
//		}
//	}
//
//	return retMap
//}

///////////////////////////////////////////////////////////////////////////////////////
//
// 以下为内部函数
//
///////////////////////////////////////////////////////////////////////////////////////

func (d *DataBox) storeKVIntoMap(key reflect.Value, value reflect.Value, m *reflect.Value) {
	fieldDataInMapSlice := m.MapIndex(key)
	if !fieldDataInMapSlice.IsValid() {
		fieldDataInMapSlice = reflect.MakeSlice(m.Type().Elem(), 0, 1)
	}

	fieldDataInMapSlice = reflect.Append(fieldDataInMapSlice, value)
	m.SetMapIndex(key, fieldDataInMapSlice)
}

// 当使用FetchObjsByKeysRetSliceFunc函数时，返回的成员数不定，因此临时map把它们处理成slice，这里需要把slice填到Struct中去
func (d *DataBox) saveToFieldWithSlice(saveStructField *reflect.StructField, destTypeIsPtr bool) error {
	if d.err != nil {
		return d.err
	}

	var retTypeIsPtr bool = false
	rettype := d.retMap.Type().Elem().Elem() // map的成员为slice，再返回slice的成员
	if rettype.Kind() == reflect.Ptr {
		retTypeIsPtr = true
	}

	// saveStructField可以是slice或者是struct或者是一个普通类型
	iter := d.retMap.MapRange()
	for iter.Next() {
		saveElemSlice := d.key2ElemSlice.MapIndex(iter.Key())
		if !saveElemSlice.IsValid() {
			continue
		}

		if iter.Value().Len() <= 0 {
			continue
		}

		// 把value放入saveElem slice中对象的fieldName处
		for i := 0; i < saveElemSlice.Len(); i++ {
			dest := saveElemSlice.Index(i).Elem().FieldByIndex(saveStructField.Index)
			if destTypeIsPtr {
				if retTypeIsPtr {
					dest.Set(iter.Value().Index(0))
				} else {
					dest.Set(iter.Value().Index(0).Addr())
				}
			} else {
				if retTypeIsPtr {
					dest.Set(iter.Value().Index(0).Elem())
				} else {
					dest.Set(iter.Value().Index(0))
				}
			}
		}
	}

	return nil
}

// convert Slice from []Obj to []*Obj  or  from []*Obj to []Obj
func convertSlice(value reflect.Value, destSliceType reflect.Type) reflect.Value {
	if value.Type() == destSliceType {
		return value
	}

	destSlice := reflect.MakeSlice(destSliceType, 0, value.Len())

	var ptr2Struct bool = false
	if value.Type().Elem().Kind() == reflect.Ptr && destSliceType.Elem().Kind() == reflect.Struct {
		ptr2Struct = true
	}

	var struct2Prt bool = false
	if value.Type().Elem().Kind() == reflect.Struct && destSliceType.Elem().Kind() == reflect.Ptr {
		struct2Prt = true
	}

	for i := 0; i < value.Len(); i++ {
		if ptr2Struct {
			destSlice = reflect.Append(destSlice, value.Index(i).Elem())
		}

		if struct2Prt {
			destSlice = reflect.Append(destSlice, value.Index(i).Addr())
		}
	}

	return destSlice
}


