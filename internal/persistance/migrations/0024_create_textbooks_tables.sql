-- +goose Up
-- Create subjects table
CREATE TABLE subjects (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create sub_subjects table
CREATE TABLE sub_subjects (
    id UUID PRIMARY KEY,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(subject_id, name)
);

-- Create indexes for faster queries
CREATE INDEX idx_subjects_name ON subjects(name);
CREATE INDEX idx_sub_subjects_subject_id ON sub_subjects(subject_id);
CREATE INDEX idx_sub_subjects_name ON sub_subjects(name);

-- +goose Down
DROP INDEX IF EXISTS idx_sub_subjects_name;
DROP INDEX IF EXISTS idx_sub_subjects_subject_id;
DROP INDEX IF EXISTS idx_subjects_name;
DROP TABLE IF EXISTS sub_subjects;
DROP TABLE IF EXISTS subjects;
