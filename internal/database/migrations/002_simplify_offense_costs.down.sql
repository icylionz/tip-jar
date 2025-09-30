ALTER TABLE offense_types
  ADD COLUMN cost_type VARCHAR(20) DEFAULT 'monetary',
  ADD COLUMN cost_action TEXT,
  DROP COLUMN cost_unit;