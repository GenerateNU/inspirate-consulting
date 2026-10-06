DROP TABLE IF EXISTS preferences;
CREATE TEMPORARY TABLE preferences AS
SELECT n.user_id, n.days_before_due
FROM notification_preferences AS n
WHERE n.email_enabled = TRUE;


SELECT t.deadline, t.description, t.user_id
FROM todo_items AS t
JOIN preferences AS p
ON p.user_id=t.user_id
WHERE t.completed_at IS NULL
AND DATEDIFF(day, now(), t.deadline) <= p.days_before_due;
