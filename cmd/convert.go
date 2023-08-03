/*
Copyright © 2023 KAI CHU CHUNG
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"

	"github.com/cage1016/alfred-xlsx2csv/alfred"
)

var (
	av = aw.NewArgVars()
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert Xlsx to CSV",
	Run:   runConvertCmd,
}

func ErrorHandle(err error) {
	av.Var("error", err.Error())
	if err := av.Send(); err != nil {
		wf.Fatalf("failed to send args to Alfred: %v", err)
	}
}

func runConvertCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		ErrorHandle(fmt.Errorf("No file specified"))
		logrus.Errorf("No file specified")
		return
	}
	fp := args[0]

	dir := filepath.Dir(fp)
	fn := strings.TrimSuffix(path.Base(fp), filepath.Ext(path.Base(fp)))

	xlsx, err := excelize.OpenFile(fp)
	if err != nil {
		ErrorHandle(errors.Wrap(err, "Error opening the Excel file"))
		logrus.Errorf("Error opening the Excel file: %v", err)
		return
	}

	index := 0
	i, _ := cmd.Flags().GetInt("index")

	if i != -1 {
		index = i
		av.Var("action", "")
	} else {
		sheetMap := xlsx.GetSheetMap()
		if len(sheetMap) > 1 {
			av.Var("file", fp)
			av.Var("action", "choose_sheet")
			if err := av.Send(); err != nil {
				wf.Fatalf("failed to send args to Alfred: %v", err)
			}
			return
		}
	}

	rows, err := xlsx.GetRows(xlsx.GetSheetName(index))
	if err != nil {
		ErrorHandle(errors.Wrap(err, "Error getting rows from Excel"))
		logrus.Errorf("Error getting rows from Excel: %v", err)
		return
	}

	// get sheet name from file with extension
	cfp := filepath.Join(dir, fn+".csv")
	csvFile, err := os.Create(cfp)
	if err != nil {
		ErrorHandle(errors.Wrap(err, "Error creating CSV file"))
		logrus.Errorf("Error creating CSV file: %v", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// set delimiter
	writer.Comma = rune(alfred.GetDelimiter(wf)[0])

	for _, row := range rows {
		writer.Write(row)
	}

	av.Var("file", fp)
	av.Var("csvfile", cfp)
	if err := av.Send(); err != nil {
		wf.Fatalf("failed to send args to Alfred: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.PersistentFlags().IntP("index", "i", -1, "Index of sheet")
	convertCmd.PersistentFlags().StringP("delimiter", "d", ",", "Delimiter to use between fields")
}
