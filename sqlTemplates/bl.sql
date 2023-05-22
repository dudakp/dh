{{- /*gotype: dh/pkg/executor.simplePredicate*/ -}}

select *
from BOOK
where {{.column}} like ({{.arg}})