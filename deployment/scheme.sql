
CREATE TYPE taskstatus AS enum ('done', 'runned', 'waiting', 'failed');

CREATE TABLE task
(
  id          SERIAL NOT NULL CONSTRAINT task_pkey PRIMARY KEY,
  created_at  TIMESTAMP DEFAULT Now() NOT NULL,
  started_at  TIMESTAMP,
  finished_at TIMESTAMP,
  status      TASKSTATUS DEFAULT 'waiting' :: taskstatus NOT NULL
);

CREATE TABLE data
(
  id     SERIAL NOT NULL CONSTRAINT data_pkey PRIMARY KEY,
  taskid INTEGER NOT NULL CONSTRAINT data_task_id_fk REFERENCES task,
  data   JSON
);

CREATE UNIQUE INDEX task_id_uindex
  ON task (id);

CREATE UNIQUE INDEX data_id_uindex
  ON data (id);