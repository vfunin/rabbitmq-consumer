# Go RabbitMQ Consumer

---

The application reads messages from the queue and saves them to the database.

#### Supported databases
- [x] PostgreSQL
- [x] MySQL

---

## Building the source

```bash
make build
```
or for build with flags
```bash
make build-flags
```

### Local environment

```bash
make all
```
---

## Using the app

- Create a table in your database:

PostgreSQL:
```postgresql
CREATE TABLE public.messages (
  id bigserial PRIMARY KEY,
  correlation_id text,
  created_at timestamptz default CURRENT_TIMESTAMP,
  message text
);


CREATE INDEX messages_correlation_id_index
    ON messages (correlation_id);
```

MySQL
```mysql
CREATE TABLE messages
(
    id             int auto_increment primary key,
    correlation_id varchar(255)                        null,
    created_at     timestamp default CURRENT_TIMESTAMP not null,
    message        text                                null
);

CREATE INDEX messages_correlation_id_index
    ON messages (correlation_id);


```

- Change config.yaml
- Run the app:
```bash
./rabbitmq-consumer -config config.yaml
```
---

## Development plan

- [x] Base functionality
- [ ] Add auto reconnect to DB
- [ ] Change structure
- [ ] Add Sentry support
- [ ] Add tests
- [ ] Add optional saving headers and properties
