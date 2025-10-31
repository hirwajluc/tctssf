-- Add approval tracking columns to loans table
ALTER TABLE loans 
  ADD COLUMN treasurer_approved_by INT(11) DEFAULT NULL AFTER approved_by,
  ADD COLUMN treasurer_approved_at TIMESTAMP NULL DEFAULT NULL AFTER treasurer_approved_by,
  ADD COLUMN vice_president_approved_by INT(11) DEFAULT NULL AFTER treasurer_approved_at,
  ADD COLUMN vice_president_approved_at TIMESTAMP NULL DEFAULT NULL AFTER vice_president_approved_by,
  ADD COLUMN president_approved_by INT(11) DEFAULT NULL AFTER vice_president_approved_at,
  ADD COLUMN president_approved_at TIMESTAMP NULL DEFAULT NULL AFTER president_approved_by,
  ADD COLUMN rejected_by INT(11) DEFAULT NULL AFTER president_approved_at,
  ADD COLUMN rejected_at TIMESTAMP NULL DEFAULT NULL AFTER rejected_by,
  ADD COLUMN rejection_reason TEXT DEFAULT NULL AFTER rejected_at;

-- Show updated structure
DESCRIBE loans;
