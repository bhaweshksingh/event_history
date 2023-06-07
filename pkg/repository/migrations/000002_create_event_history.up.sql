create table if not exists event_history (
    key varchar(100),
    value varchar(100),
    user_id varchar(100),
    action varchar(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
