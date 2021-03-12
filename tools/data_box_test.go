package tools

import (
	"fmt"
	"testing"
)

// go test -v data_box.go data_box_test.go -test.run TestJoinToVByV
// go test -v data_box.go data_box_test.go -test.run TestJoinByMapPV

type V struct {
	Id    int
	VName string
}

type T struct {
	KeyId   int
	TName   string
	V       V
	PV      *V
	VSlice  []V
	PVSlice []*V
}

func TestJoinToVByV(t *testing.T) {
	TSlice := []T{T{KeyId: 1, TName: "n1"}, T{KeyId: 2, TName: "n2"}, T{KeyId: 3, TName: "n3"}}
	err := NewDataBox(TSlice).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
		retMap := make(map[int]*V)
		KeyIds := keywords.([]int)
		// 一次性从数据库中得到所有的Mobjs
		if len(KeyIds) > 0 {
			VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
			for i, v := range VSlice {
				retMap[v.Id] = &VSlice[i]
			}
		}

		return retMap
	}).SaveToField("V")

	if err != nil {
		t.Fatal(err)
	} else {
		for _, et := range TSlice {
			fmt.Printf("%+v", et)
			if et.KeyId != et.V.Id {
				t.Fatal("et.KeyId != et.V.Id")
			}
		}
	}
}

func TestJoinToPVByV(t *testing.T) {
	TSlice := []T{{KeyId: 1, TName: "n1"}, T{KeyId: 2, TName: "n2"}, T{KeyId: 3, TName: "n3"}}
	err := NewDataBox(TSlice).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
		retMap := make(map[int]*V)
		KeyIds := keywords.([]int)
		// 一次性从数据库中得到所有的Mobjs
		if len(KeyIds) > 0 {
			VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
			for i, _ := range VSlice {
				retMap[VSlice[i].Id] = &VSlice[i]
			}
		}

		return retMap
	}).SaveToField("PV")

	if err != nil {
		t.Fatal(err)
	} else {
		for _, et := range TSlice {
			fmt.Printf("%+v, PV: %+v", et, *et.PV)
			if et.KeyId != et.PV.Id {
				t.Fatal("et.KeyId != et.PV.Id")
			}
		}
	}
}

func TestJoinToVSliceByV(t *testing.T) {
	TSlice := []T{{KeyId: 1, TName: "n1"}, T{KeyId: 2, TName: "n2"}, T{KeyId: 3, TName: "n3"}}
	err := NewDataBox(TSlice).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
		retMap := make(map[int][]V)
		KeyIds := keywords.([]int)
		// 一次性从数据库中得到所有的Mobjs
		if len(KeyIds) > 0 {
			VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
			for _, v := range VSlice {
				retMap[v.Id] = []V{v}
			}
		}

		return retMap
	}).SaveToField("VSlice")

	if err != nil {
		t.Fatal(err)
	} else {
		for _, et := range TSlice {
			fmt.Printf("%+v", et)
			if et.KeyId != et.VSlice[0].Id {
				t.Fatal("et.KeyId != et.VSlice.Id")
			}
		}
	}
}

func TestJoinToPVSliceByV(t *testing.T) {
	TSlice := []T{{KeyId: 1, TName: "n1"}, T{KeyId: 2, TName: "n2"}, T{KeyId: 3, TName: "n3"}}
	err := NewDataBox(TSlice).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
		retMap := make(map[int][]*V)
		KeyIds := keywords.([]int)
		// 一次性从数据库中得到所有的Mobjs
		if len(KeyIds) > 0 {
			VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
			for i, v := range VSlice {
				retMap[v.Id] = []*V{&VSlice[i]}
			}
		}

		return retMap
	}).SaveToField("PVSlice")

	if err != nil {
		t.Fatal(err)
	} else {
		for _, et := range TSlice {
			fmt.Printf("%+v PVSlice:%+v", et, et.PVSlice[0])
			if et.KeyId != et.PVSlice[0].Id {
				t.Fatal("et.KeyId != et.PVSlice.Id")
			}
		}
	}
}

//func TestSlice2MapByField(t *testing.T) {
//	VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}}
//	retMap := Slice2MapByField(VSlice, "Id")
//	if len(retMap) == 0 {
//		t.Fatal("fail")
//	}
//
//	for k, v := range retMap {
//		fmt.Println(k.(int), v.(*V))
//		fmt.Printf("k: %+v, v:%+v\n", k, v)
//	}
//
//	VSlice2 := []*V{&V{Id: 1, VName: "n1"}, &V{Id: 2, VName: "n2"}, &V{Id: 3, VName: "n3"}}
//	retMap = Slice2MapByField(VSlice2, "Id")
//	if len(retMap) == 0 {
//		t.Fatal("fail")
//	}
//
//	for k, v := range retMap {
//		fmt.Println(k.(int), v.(*V))
//		fmt.Printf("k: %+v, v:%+v\n", k, v)
//	}
//}

func TestJoinToVByPT(t *testing.T) {
	TSlice := []*T{&T{KeyId: 1, TName: "n1"}, &T{KeyId: 2, TName: "n2"}, &T{KeyId: 3, TName: "n3"}}
	err := NewDataBox(TSlice).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
		retMap := make(map[int]V)
		KeyIds := keywords.([]int)
		// 一次性从数据库中得到所有的Mobjs
		if len(KeyIds) > 0 {
			VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
			for _, v := range VSlice {
				retMap[v.Id] = v
			}
		}

		return retMap
	}).SaveToField("V")

	if err != nil {
		t.Fatal(err)
	} else {
		for _, et := range TSlice {
			fmt.Printf("%+v", et)
			if et.KeyId != et.V.Id {
				t.Fatal("et.KeyId != et.V.Id")
			}
		}
	}
}

func TestJoinToPVByPT(t *testing.T) {
	TSlice := []*T{&T{KeyId: 1, TName: "n1"}, &T{KeyId: 2, TName: "n2"}, &T{KeyId: 3, TName: "n3"}}
	err := NewDataBox(TSlice).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
		retMap := make(map[int]*V)
		KeyIds := keywords.([]int)
		// 一次性从数据库中得到所有的Mobjs
		if len(KeyIds) > 0 {
			VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
			for i, _ := range VSlice {
				retMap[VSlice[i].Id] = &VSlice[i]
			}
		}

		return retMap
	}).SaveToField("PV")

	if err != nil {
		t.Fatal(err)
	} else {
		for _, et := range TSlice {
			fmt.Printf("%+v, PV: %+v", et, *et.PV)
			if et.KeyId != et.PV.Id {
				t.Fatal("et.KeyId != et.PV.Id")
			}
		}
	}
}

func TestJoinToVSliceByPT(t *testing.T) {
	TSlice := []*T{&T{KeyId: 1, TName: "n1"}, &T{KeyId: 2, TName: "n2"}, &T{KeyId: 3, TName: "n3"}}
	err := NewDataBox(TSlice).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
		retMap := make(map[int][]V)
		KeyIds := keywords.([]int)
		// 一次性从数据库中得到所有的Mobjs
		if len(KeyIds) > 0 {
			VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
			for _, v := range VSlice {
				retMap[v.Id] = []V{v}
			}
		}

		return retMap
	}).SaveToField("VSlice")

	if err != nil {
		t.Fatal(err)
	} else {
		for _, et := range TSlice {
			fmt.Printf("%+v", et)
			if et.KeyId != et.VSlice[0].Id {
				t.Fatal("et.KeyId != et.VSlice.Id")
			}
		}
	}
}

func TestJoinToPVSliceByPT(t *testing.T) {
	TSlice := []*T{&T{KeyId: 1, TName: "n1"}, &T{KeyId: 2, TName: "n2"}, &T{KeyId: 3, TName: "n3"}}
	err := NewDataBox(TSlice).KeyField("KeyId").JoinByMap(func(keywords interface{}) interface{} {
		retMap := make(map[int][]*V)
		KeyIds := keywords.([]int)
		// 一次性从数据库中得到所有的Mobjs
		if len(KeyIds) > 0 {
			VSlice := []V{{Id: 1, VName: "n1"}, V{Id: 2, VName: "n2"}, V{Id: 3, VName: "n3"}} // get from db
			for i, v := range VSlice {
				retMap[v.Id] = []*V{&VSlice[i]}
			}
		}

		return retMap
	}).SaveToField("PVSlice")

	if err != nil {
		t.Fatal(err)
	} else {
		for _, et := range TSlice {
			fmt.Printf("%+v PVSlice:%+v", et, et.PVSlice[0])
			if et.KeyId != et.PVSlice[0].Id {
				t.Fatal("et.KeyId != et.PVSlice.Id")
			}
		}
	}
}