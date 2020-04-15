-- NOT WORKING
UPDATE users
SET (bio, will, grade) = (
    coalesce('new bio v2', bio),
    coalesce('', will),
    coalesce(grade, 0)
);