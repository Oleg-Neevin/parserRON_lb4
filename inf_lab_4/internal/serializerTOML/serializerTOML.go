package serializertoml

import (
	"fmt"
	"io"
	"strings"

	s "gitlab.se.ifmo.ru/s503298/inf_lab_4/pkg"
)

func escapeTomlString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func WriteTOML(w io.Writer, sch s.Schedule) error {
	for _, day := range sch.Days {
		if _, err := fmt.Fprintln(w, "[[days]]"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "name = \"%s\"\n", escapeTomlString(day.Name)); err != nil {
			return err
		}
		for _, les := range day.Lessons {
			if _, err := fmt.Fprintln(w, "\n[[days.lessons]]"); err != nil {
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
				if _, err := fmt.Fprintf(w, "%s = \"%s\"\n", kv.k, escapeTomlString(kv.v)); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprintln(w, ""); err != nil {
			return err
		}
	}
	return nil
}
