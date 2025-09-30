-- Simplify offense cost system to just amount + unit
ALTER TABLE offense_types 
  DROP COLUMN cost_type,
  DROP COLUMN cost_action,
  ADD COLUMN cost_unit VARCHAR(100);

-- Set default unit for existing monetary offenses
UPDATE offense_types 
SET cost_unit = 'dollars' 
WHERE cost_amount IS NOT NULL;