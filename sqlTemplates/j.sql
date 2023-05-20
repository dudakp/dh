{{- /*gotype: dh/pkg/executor.simplePredicate*/ -}}

select * from BOOK book join AUTHOR author on author.ID = book.AUTHOR_ID where {{.column}} in {{.arg}}
