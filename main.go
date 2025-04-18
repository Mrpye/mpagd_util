package main

import (
	"fmt"

	"github.com/Mrpye/mpagd_util/cmd"
	"github.com/Mrpye/mpagd_util/mpagd"
)

func ReadBlankProjectImportWriteCompare() {
	filePath := "blank.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	apjFile.ReadAPJ()
	//apjFile.Display()

	agdFilePath := "deleteme.agd"
	opt := mpagd.CreateImportOptions()
	opt.SetOwOptionsTrue()
	opt.SetIgnoreOptions(true, true, false, true, true, true, true, true, true)
	apjFile.ImportAGD(agdFilePath, opt)

	// // Example comparison
	outputFilePath := "output/output.apj" // Replace with the other file path
	apjFile.WriteAPJ(outputFilePath)

	otherAPJFile := mpagd.NewAPJFile(outputFilePath)
	otherAPJFile.ReadAPJ()

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		fmt.Println("Differences found:")
		for key, diff := range differences {
			fmt.Printf("Key: %s\n", key)
			switch key {
			case "Objects":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "Sprites":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "SpriteInfo":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "ULAPalette":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			default:
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			}
		}
	} else {
		fmt.Println("No differences found.")
	}

}

func ReadBlankProjectCreateBlankWriteCompare() {
	filePath := "blank.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	apjFile.ReadAPJ()
	//apjFile.Display()

	//agdFilePath := "deleteme.agd"
	//apjFile.ImportAGD(agdFilePath)

	// // Example comparison
	outputFilePath := "output/output.apj" // Replace with the other file path

	otherAPJFile := mpagd.NewAPJFile(outputFilePath)
	otherAPJFile.CreateBlank()

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		fmt.Println("Differences found:")
		for key, diff := range differences {
			fmt.Printf("Key: %s\n", key)
			switch key {
			case "Objects":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "Sprites":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "SpriteInfo":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "ULAPalette":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			default:
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			}
		}
	} else {
		fmt.Println("No differences found.")
	}

}
func ReadBlankProjectWriteCompare() {
	filePath := "blank.apj" // Replace with the actual file path
	apjFile := mpagd.NewAPJFile(filePath)
	apjFile.ReadAPJ()
	//apjFile.Display()

	// // Example comparison
	outputFilePath := "output/output.apj" // Replace with the other file path

	otherAPJFile := mpagd.NewAPJFile(outputFilePath)
	if err := otherAPJFile.ReadAPJ(); err != nil {
		fmt.Println("Error reading other APJ file:", err)
		return
	}
	apjFile.SaveAsYAML("output/apj.yaml")
	otherAPJFile.SaveAsYAML("output/other_apj.yaml")

	differences := apjFile.CompareData(otherAPJFile)
	if len(differences) > 0 {
		fmt.Println("Differences found:")
		for key, diff := range differences {
			fmt.Printf("Key: %s\n", key)
			switch key {
			case "Objects":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "Sprites":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "SpriteInfo":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			case "ULAPalette":
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			default:
				fmt.Printf("  Self: %+v\n  Other: %+v\n", diff.(map[string]interface{})["self"], diff.(map[string]interface{})["other"])
			}
		}
	} else {
		fmt.Println("No differences found.")
	}

}

func main() {
	//ReadProjectWriteCompare()
	//ReadProjectImportWriteCompare()
	// ReadBlankProjectWriteCompare()
	// ReadBlankProjectCreateBlankWriteCompare()
	// CreateBlankProjectWriteCompare()
	// ReadBlankProjectImportWriteCompare()
	//Test3()
	//Print()
	//ImportWriteCompare()
	//ReadProjectWriteCompare()
	cmd.Execute()

}
