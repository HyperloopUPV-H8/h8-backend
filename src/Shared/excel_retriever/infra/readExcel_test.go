package infra

import (
	"log"
	"path"
	"reflect"
	"testing"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever/domain"
	"github.com/xuri/excelize/v2"
)

func TestReadExcel(t *testing.T) {
	type testCase struct {
		file   string
		expect domain.Document
	}

	tests := map[string]testCase{
		"test1": {
			file:   "test1.xlsx",
			expect: test1,
		},
		"test2": {
			file:   "test2.xlsx",
			expect: test2,
		},
		"test3": {
			file:   "test3.xlsx",
			expect: test3,
		},
		"test4": {
			file:   "test4.xlsx",
			expect: test4,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			file := openExcelFile(path.Join("test", "resources", test.file))
			got := GetDocument(file)
			if !reflect.DeepEqual(got, test.expect) {
				t.Fatalf("expected %v, got %v", test.expect, got)
			}
		})
	}
}

func openExcelFile(name string) *excelize.File {
	file, err := excelize.OpenFile(name)
	if err != nil {
		log.Fatalln(err)
	}

	return file
}

var test1 domain.Document = domain.Document{
	Sheets: map[string]domain.Sheet{
		"Test1": {
			Name: "Test1",
			Tables: map[string]domain.Table{
				"PacketDescription": {
					Name: "PacketDescription",
					Rows: [][]string{
						{"1", "Voltage", "200", "Input", "UDP"},
						{"2", "Speed", "300", "Input", "UDP"},
						{"3", "Current", "400", "Input", "UDP"},
						{"4", "Airgap", "500", "Input", "UDP"},
						{"5", "Position", "600", "Input", "UDP"},
						{"6", "Battery", "700", "Input", "UDP"},
					},
				},

				"PacketStructure": {
					Name: "PacketStructure",
					Rows: [][]string{
						{"Voltage", "Speed", "Current", "Airgap", "Position", "Battery"},
						{"Voltage0", "Speed0", "Current0", "Airgap0", "Position0", "Battery0"},
					},
				},

				"ValueDescription": {
					Name: "ValueDescription",
					Rows: [][]string{
						{"Voltage0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Speed0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Current0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Airgap0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Battery0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
					},
				},
			},
		},
	},
}

var test2 domain.Document = domain.Document{
	Sheets: map[string]domain.Sheet{
		"Test2": {
			Name: "Test2",
			Tables: map[string]domain.Table{
				"PacketDescription": {
					Name: "PacketDescription",
					Rows: [][]string{
						{"1", "Voltage", "200", "Input", "UDP"},
						{"2", "Speed", "300", "Input", "UDP"},
						{"3", "Current", "400", "Input", "UDP"},
						{"4", "Airgap", "500", "Input", "UDP"},
						{"5", "Position", "600", "Input", "UDP"},
						{"6", "Battery", "700", "Input", "UDP"},
					},
				},

				"PacketStructure": {
					Name: "PacketStructure",
					Rows: [][]string{
						{"Voltage", "Speed", "Current", "Airgap", "Position", "Battery"},
						{"Voltage0", "Speed0", "Current0", "Airgap0", "Position0", "Battery0"},
						{"", "Speed1", "Current1", "", "Position1", ""},
						{"", "", "Current2", "", "Position2", ""},
						{"", "", "", "", "Position3", ""},
					},
				},

				"ValueDescription": {
					Name: "ValueDescription",
					Rows: [][]string{
						{"Voltage0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Speed0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Speed1", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Current0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Current1", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Current2", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Airgap0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position1", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position2", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position3", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Battery0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
					},
				},
			},
		},
	},
}

var test3 domain.Document = domain.Document{
	Sheets: map[string]domain.Sheet{
		"Test3": {
			Name: "Test3",
			Tables: map[string]domain.Table{
				"PacketDescription": {
					Name: "PacketDescription",
					Rows: [][]string{
						{"1", "Voltage", "200", "Input", "UDP"},
						{"2", "Speed", "300", "Input", "UDP"},
						{"3", "Current", "400", "Input", "UDP"},
						{"4", "Airgap", "500", "Input", "UDP"},
						{"5", "Position", "600", "Input", "UDP"},
						{"6", "Battery", "700", "Input", "UDP"},
					},
				},

				"PacketStructure": {
					Name: "PacketStructure",
					Rows: [][]string{
						{"Voltage", "Speed", "Current", "Airgap", "Position", "Battery"},
						{"Voltage0", "Speed0", "Current0", "Airgap0", "Position0", "Battery0"},
					},
				},

				"ValueDescription": {
					Name: "ValueDescription",
					Rows: [][]string{
						{"Voltage0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Speed0", "bool", "", "", "", ""},
						{"Current0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Airgap0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position0", "ENUM(a,b,c)", "", "", "", ""},
						{"Battery0", "bool", "", "", "", ""},
					},
				},
			},
		},
	},
}

var test4 domain.Document = domain.Document{
	Sheets: map[string]domain.Sheet{
		"Test4": {
			Name: "Test4",
			Tables: map[string]domain.Table{
				"PacketDescription": {
					Name: "PacketDescription",
					Rows: [][]string{
						{"1", "Voltage", "200", "Input", "UDP"},
						{"2", "Speed", "300", "Input", "UDP"},
						{"3", "Current", "400", "Input", "UDP"},
						{"4", "Airgap", "500", "Input", "UDP"},
						{"5", "Position", "600", "Input", "UDP"},
						{"6", "Battery", "700", "Input", "UDP"},
					},
				},

				"PacketStructure": {
					Name: "PacketStructure",
					Rows: [][]string{
						{"Voltage", "Speed", "Current", "Airgap", "Position", "Battery"},
						{"Voltage0", "Speed0", "Current0", "Airgap0", "Position0", "Battery0"},
					},
				},

				"ValueDescription": {
					Name: "ValueDescription",
					Rows: [][]string{
						{"Voltage0", "uint8", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Speed0", "float64", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Current0", "uint16", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Airgap0", "int32", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Position0", "foo", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
						{"Battery0", "float32", "cdeg#/100#", "deg##", "[0,10]", "[-10,20]"},
					},
				},
			},
		},
	},
}
