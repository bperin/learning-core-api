-- name: ListSubjects :many
SELECT * FROM subjects
ORDER BY name ASC;

-- name: ListSubSubjectsBySubjectID :many
SELECT * FROM sub_subjects
WHERE subject_id = $1
ORDER BY name ASC;
