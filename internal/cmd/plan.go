package cmd

import (
	"fmt"
	"os"

	"github.com/beauknowssoftware/makehcl/internal/plan"
	"github.com/jessevdk/go-flags"
	"github.com/olekukonko/tablewriter"
)

type PlanCommand struct {
	Filename flags.Filename `short:"f" long:"filename"`
	All      bool           `short:"a" long:"all"`
	Reason   bool           `short:"r" long:"reason"`
	Goal     Goal           `positional-args:"yes"`
}

func reasonDescription(r plan.Reason) string {
	switch r.ReasonType {
	case plan.DoesNotExist:
		return fmt.Sprintf("%v does not exist", r.Target)
	case plan.DependencyPlanned:
		return fmt.Sprintf("dependency %v is planned", r.Target)
	case plan.OlderThanDependency:
		return fmt.Sprintf("dependency %v has changed", r.Target)
	default:
		return fmt.Sprintf("unknown (%v)", r.Target)
	}
}

func (c *PlanCommand) Execute(_ []string) error {
	var o plan.DoOptions
	o.Filename = string(c.Filename)
	o.IgnoreLastModified = c.All
	o.Goal = c.Goal.strings()

	p, err := plan.Do(o)
	if err != nil {
		return err
	}

	if c.Reason && !c.All && len(p) > 0 {
		printReasonTable(p)
	} else {
		for _, e := range p {
			fmt.Println(e.Target)
		}
	}

	return nil
}

func printReasonTable(p plan.Plan) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Target", "Reason"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	for _, e := range p {
		table.Append([]string{e.Target, reasonDescription(e.Reason)})
	}

	table.Render()
}
