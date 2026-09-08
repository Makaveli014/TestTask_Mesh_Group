-- SQL migration to create the segmentation table
CREATE TABLE IF NOT EXISTS segmentation (
    id SERIAL PRIMARY KEY,
    address_sap_id VARCHAR(255) NOT NULL UNIQUE,
    adr_segment VARCHAR(16),
    segment_id BIGINT
);

CREATE INDEX IF NOT EXISTS idx_segmentation_address_sap_id ON segmentation(address_sap_id);
