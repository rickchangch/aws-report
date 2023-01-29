/*
Copyright © 2022 Rick Chang <medo972283@gmail.com>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/rickchangch/aws-report/models"
	"github.com/rickchangch/aws-report/utils"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// weeklyCost represents the weekly command
var weeklyCost = &cobra.Command{
	Use:     "weekly",
	Aliases: []string{"w"},
	Short:   "Generate weekly analytics report through fetching the data from company's PostgreSQL DB.",
	Long: `Generate weekly analytics report through fetching the data from company's PostgreSQL DB,
which will sync data from AWS Costexplorer periodically, instead of requesting AWS Costexplorer SDK directly.

In this way, unnecessary costs caused by multiple accesses per day can be avoided.`,
	RunE: generateWeeklyReport,
}

var (
	dbhost    string
	dbport    string
	dbname    string
	user      string
	password  string
	startDate string
	endDate   string
	abridge   bool
)

func init() {
	rootCmd.AddCommand(weeklyCost)

	// DB info
	weeklyCost.Flags().StringVarP(&dbhost, "dbhost", "d", "localhost", "Host address of Postgre DB.")
	weeklyCost.Flags().StringVarP(&dbport, "dbport", "p", "5432", "Port.")
	weeklyCost.Flags().StringVarP(&dbname, "dbname", "n", "awsdeacon", "DB name.")
	weeklyCost.Flags().StringVarP(&user, "user", "u", "root", "User account.")
	weeklyCost.Flags().StringVarP(&password, "password", "w", "root", "User password.")

	// Filters
	weeklyCost.Flags().StringVarP(&startDate, "start-date", "s", "2023-01-01", "Start date, YYYY-MM-DD.")
	weeklyCost.Flags().StringVarP(&endDate, "end-date", "e", "2023-01-07", "End date, YYYY-MM-DD.")
	weeklyCost.Flags().BoolVarP(&abridge, "abridge", "a", false, "Whether filter out 0 cost services.")
}

// Generate weekly report of Costexplorer.
func generateWeeklyReport(cmd *cobra.Command, args []string) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbhost,
		user,
		password,
		dbname,
		dbport,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return utils.ErrDBConnectionFailed
	}

	startDate, err := utils.DateUtil.TransferToTimeFormat(startDate)
	if err != nil {
		return utils.ErrInvalidDateFormat
	}
	endDate, err := utils.DateUtil.TransferToTimeFormat(endDate)
	if err != nil {
		return utils.ErrInvalidDateFormat
	}

	// Fetch weekly data separately
	weeks := utils.DateUtil.SliceByWeek(startDate, endDate)
	weekCosts, err := models.CostModel.GetByMultiDateRange(db, weeks)
	if err != nil {
		return fmt.Errorf(utils.ErrDBQueryFail.Error()+": ", err)
	}

	// Aggregate cost by service
	headers := []string{}
	defRow := []float64{}
	for i := 0; i < len(weeks); i++ {
		defRow = append(defRow, 0.0)
		// Collect week days as headers
		headers = append(
			headers,
			fmt.Sprintf(
				"%s-%s",
				weeks[i][0].Format(utils.DATE_LAYOUT),
				weeks[i][1].Format(utils.DATE_LAYOUT),
			),
		)
	}
	fmt.Println(strings.Join(headers, ", "))

	serviceCosts := map[string][]float64{}
	for i, weekCost := range weekCosts {
		for _, row := range weekCost {
			if _, ok := serviceCosts[row.Service]; !ok {
				serviceCosts[row.Service] = append([]float64{}, defRow...)
			}
			serviceCosts[row.Service][i] += row.Value
		}
	}

	// Convert to CSV string format for printing
	result := [][]string{}
	for serviceName, costs := range serviceCosts {
		row := []string{}
		row = append(row, serviceName)
		count := 0.0
		for _, v := range costs {
			count += v
			row = append(row, fmt.Sprintf("%.4f", v))
		}

		// skip 0 cost
		if abridge && count == 0 {
			continue
		}

		result = append(result, row)
	}
	for _, row := range result {
		fmt.Println(strings.Join(row, ", "))
	}

	return nil
}
