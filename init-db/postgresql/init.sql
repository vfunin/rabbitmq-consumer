CREATE TABLE public.messages (
  id bigserial PRIMARY KEY,
  correlation_id text,
  created_at timestamptz default CURRENT_TIMESTAMP,
  message text
);


CREATE INDEX messages_correlation_id_index
    ON messages (correlation_id);
