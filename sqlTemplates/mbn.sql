select movie.name as movie_name, director, b.TITLE based_on, a.SURNAME as book_author
from movie
         left join book b on b.id = based_on
         left join author a on a.id = b.author_id
where movie.name like '{{.Arg}}%';
