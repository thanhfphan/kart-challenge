-- Create outbox table for reliable event publishing
CREATE TABLE IF NOT EXISTS outbox_events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(36) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL DEFAULT 'order',
    event_data JSON NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    version INT NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    processed_at BIGINT NULL,
    
    INDEX idx_outbox_events_status (status),
    INDEX idx_outbox_events_created_at (created_at),
    INDEX idx_outbox_events_aggregate (aggregate_type, aggregate_id),
    INDEX idx_outbox_events_type (event_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
