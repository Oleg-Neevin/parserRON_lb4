package serializerxml

import (
	"fmt"
	"io"
	"strings"

	s "gitlab.se.ifmo.ru/s503298/inf_lab_4/pkg"
)

func escapeXML(str string) string {
	str = strings.ReplaceAll(str, "&", "&amp;")
	str = strings.ReplaceAll(str, "<", "&lt;")
	str = strings.ReplaceAll(str, ">", "&gt;")
	str = strings.ReplaceAll(str, `"`, "&quot;")
	str = strings.ReplaceAll(str, `'`, "&apos;")
	return str
}

func WriteXML(w io.Writer, sch s.Schedule) error {
	if _, err := fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "<schedule>"); err != nil {
		return err
	}

	for _, day := range sch.Days {
		if _, err := fmt.Fprintln(w, "  <day>"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "    <name>%s</name>\n", escapeXML(day.Name)); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "    <lessons>"); err != nil {
			return err
		}
		for _, les := range day.Lessons {
			if _, err := fmt.Fprintln(w, "    <lesson>"); err != nil {
				return err
			}
			fmtStrs := []struct {
				k, v string
			}{
				{"time", les.Time},
				{"subject", les.Subject},
				{"teacher", les.Teacher},
				{"room", les.Room},
				{"building", les.Building},
				{"type", les.Type},
			}
			for _, kv := range fmtStrs {
				if _, err := fmt.Fprintf(w, "        <%s>%s</%s>\n", kv.k, escapeXML(kv.v), kv.k); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w, "      </lesson>"); err != nil {
				return err
			}

		}
		if _, err := fmt.Fprintln(w, "    </lessons>\n  </day>"); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "</schedule>"); err != nil {
		return err
	}

	return nil
}
