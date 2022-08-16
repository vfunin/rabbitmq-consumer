create table messages
(
    id             int auto_increment
        primary key,
    correlation_id varchar(255)                        null,
    created_at     timestamp default CURRENT_TIMESTAMP not null,
    message        text                                null
);

create index messages_correlation_id_index
    on messages (correlation_id);
