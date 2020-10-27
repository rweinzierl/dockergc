package lib

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func createTabWriter() *tabwriter.Writer {
	w := new(tabwriter.Writer)
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 4, 0, 3, ' ', 0)
	return w
}

func addTabRow(w *tabwriter.Writer, cells ...interface{}) {
	format := strings.Repeat("%v\t", len(cells)-1) + "%v\n"
	fmt.Fprintf(w, format, cells...)
}
