package memory

import (
	"fmt"
	"io"

	"github.com/Cloud-Foundations/Dominator/lib/format"
)

func (p *prober) writeHtml(writer io.Writer) {
	fmt.Fprintln(writer, `<style>
                          table, th, td {
                          border-collapse: collapse;
                          }
                          </style>`)
	//fmt.Fprintln(writer, `<table border="1" style="width:100%">`)
	fmt.Fprintln(writer, `<table border="1">`)
	fmt.Fprintln(writer, "  <tr>")
	fmt.Fprintln(writer, "    <th>Memory Size</th>")
	if p.haveAvailable {
		fmt.Fprintf(writer, "    <th>Utilisation: %.1f%%</th>\n",
			float64(p.total-p.available)/float64(p.total)*100)
	} else {
		fmt.Fprintf(writer, "    <th>Used: %.1f%%</th>\n",
			float64(p.total-p.free)/float64(p.total)*100)
	}
	fmt.Fprintln(writer, "  </tr>")
	fmt.Fprintf(writer, "  <tr>\n")
	fmt.Fprintf(writer, "    <td><center>%s</td>\n",
		format.FormatBytes(p.total))
	fmt.Fprint(writer, "    <td>")
	if p.haveAvailable {
		p.writeHtmlBarAvailable(writer)
	} else {
		p.writeHtmlBar(writer)
	}
	fmt.Fprintln(writer, "</td>")
	fmt.Fprintln(writer, "  </tr>")
	fmt.Fprintln(writer, "</table>")
	fmt.Fprintln(writer, "</body>")
}

func (p *prober) writeHtmlBar(writer io.Writer) {
	usedBytes := p.total - p.free
	barColour := "grey"
	leftBarWidth := float64(usedBytes) / float64(p.total)
	rightBarWidth := float64(p.free) / float64(p.total)
	if p.free < p.total/100 {
		barColour = "orange"
	}
	fmt.Fprint(writer, `<table border="0" style="width:200px"><tr>`)
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%;background-color:%s\">&nbsp;</td>",
		leftBarWidth*100, barColour)
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%;background-color:%s\">&nbsp;</td>",
		rightBarWidth*100, "white")
	fmt.Fprint(writer, "</tr></table>")
}

func (p *prober) writeHtmlBarAvailable(writer io.Writer) {
	usedBytes := p.total - p.available
	barColour := "grey"
	leftBarWidth := float64(usedBytes) / float64(p.total)
	middleBarWidth := float64(p.available-p.free) / float64(p.total)
	rightBarWidth := float64(p.free) / float64(p.total)
	if p.available < p.total/100 {
		barColour = "orange"
	}
	fmt.Fprint(writer, `<table border="0" style="width:200px"><tr>`)
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%;background-color:%s\">&nbsp;</td>",
		leftBarWidth*100, "blue")
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%;background-color:%s\">&nbsp;</td>",
		middleBarWidth*100, barColour)
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%;background-color:%s\">&nbsp;</td>",
		rightBarWidth*100, "white")
	fmt.Fprint(writer, "</tr></table>")
}
