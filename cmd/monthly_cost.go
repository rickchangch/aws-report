/*
Copyright © 2022 Rick Chang <medo972283@gmail.com>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/rickchangch/aws-report/utils"
	"github.com/spf13/cobra"
)

// monthlyCost represents the monthly command
var monthlyCost = &cobra.Command{
	Use:     "monthly",
	Aliases: []string{"m"},
	Short:   "Generate monthly analytics report via csv exported from AWS Costexplorer console.",
	Long: `Generate monthly analytics report via csv exported from AWS Costexplorer console.

To avoid unnecessary cost, download the csv report from AWS Costexplorer console,
then pass into the file via the  flag '--file|-f'
`,
	RunE: generateMonthlyReport,
}

var (
	filePath string
	order    string
	decimal  string
)

func init() {
	rootCmd.AddCommand(monthlyCost)

	monthlyCost.Flags().StringVarP(&filePath, "filePath", "f", "", "Local file path of csv file exported from AWS Costexplorer console.")
	monthlyCost.MarkFlagRequired("filePath")

	monthlyCost.Flags().StringVarP(&order, "order", "o", "amount", "Specify the column which is used to sort the result in DESC order.")
	monthlyCost.Flags().StringVarP(&decimal, "decimal", "d", "4", "Decide the digit numbers behind the decimal point.")
}

// Generate monthly report of Costexplorer.
func generateMonthlyReport(cmd *cobra.Command, args []string) error {

	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		return utils.ErrFileNotExist
	}

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	// Read from CSV
	var totalFieldName string
	idx, months := 0, []string{}
	serviceNames, serviceCosts := []string{}, map[string][]float64{}
	for {
		idx++

		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch idx {
		case 1:
			serviceNames = row
			totalFieldName = row[len(row)-1]
		case 2:
			continue // Total amount, is useless for analytics report
		default:
			// Build month header
			date, _ := utils.DateUtil.TransferToTimeFormat(row[0])
			months = append(months, fmt.Sprintf(
				"%s (%d)",
				date.Format(utils.MONTH_LAYOUT),
				utils.DateUtil.DaysInMonth(date.Month()),
			))

			// Collect costs info
			for i := 1; i < len(row); i++ {
				value, _ := strconv.ParseFloat(row[i], 64)
				serviceCosts[serviceNames[i]] = append(serviceCosts[serviceNames[i]], value)
			}
		}
	}

	// Build CSV rows
	result := [][]string{}
	var totalCostRow []string
	for serviceName, costs := range serviceCosts {
		row := []string{}
		format := fmt.Sprintf("%%.%sf", decimal)

		row = append(row, serviceName)

		for _, v := range costs {
			row = append(row, fmt.Sprintf(format, v))
		}

		// Increase Amount
		row = append(row, fmt.Sprintf(format, costs[1]-costs[0]))

		// Increase Rate
		number := costs[1] / costs[0] * 100
		if math.IsNaN(number) {
			number = 0.0
		}
		row = append(row, fmt.Sprintf("%.2f%%", number))

		// Reserve the total cost row then add it after sorting finish
		if serviceName == totalFieldName {
			totalCostRow = row
		} else {
			result = append(result, row)
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		var prevStr, nextStr string
		if order == "rate" {
			prevStr = strings.Replace(result[i][4], "%", "", 1)
			nextStr = strings.Replace(result[j][4], "%", "", 1)
		} else {
			prevStr = result[i][3]
			nextStr = result[j][3]
		}

		prev, _ := strconv.ParseFloat(prevStr, 64)
		next, _ := strconv.ParseFloat(nextStr, 64)

		return prev > next
	})

	// Pin header and total cost row
	header := []string{"Service"}
	header = append(header, months...)
	header = append(header, []string{"Increase Amount", "%"}...)
	result = append([][]string{header, totalCostRow, {"-"}}, result...)

	for _, row := range result {
		fmt.Println(strings.Join(row, ", "))
	}

	return nil
}
