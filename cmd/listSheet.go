/*
Copyright © 2023 KAI CHU CHUNG
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

// listSheetCmd represents the listSheet command
var listSheetCmd = &cobra.Command{
	Use:   "listSheet",
	Short: "List all sheets in the Excel file",
	Run:   runListSheetCmd,
}

func runListSheetCmd(cmd *cobra.Command, args []string) {
	file, _ := cmd.Flags().GetString("file")

	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		wf.Fatalf("Error opening the Excel file: %v", errors.Wrap(err, "Error opening the Excel file"))
	}

	sheetMap := xlsx.GetSheetMap()
	wf.Configure(aw.SuppressUIDs(true))
	for i, sheet := range sheetMap {
		wf.NewItem(fmt.Sprintf("%d - %s", i, sheet)).
			Subtitle(fmt.Sprintf("⇧, ↩ Convert sheet '%s' to CSV", sheet)).
			Arg(strconv.Itoa(i - 1)).
			UID(strconv.Itoa(i - 1)).
			Quicklook(file).
			Valid(true)
	}

	if len(args) >= 2 {
		wf.Filter(strings.Join(args[1:], " "))
	}

	if wf.IsEmpty() {
		wf.NewItem("No matching items").
			Subtitle("Try a different query?").
			Valid(true)
	}
	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(listSheetCmd)
	listSheetCmd.Flags().StringP("file", "f", "", "Excel file to convert")
}
