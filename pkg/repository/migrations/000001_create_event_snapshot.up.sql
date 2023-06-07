create table if not exists event_snapshot (
    key varchar(100),
    value varchar(100),
    user_id varchar(100),
    PRIMARY KEY (key, user_id)
    );
