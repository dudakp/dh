select b.*
from BOOK b
         join AUTHOR a on a.id = b.author_id
where a.surname like '{{.Arg}}%';