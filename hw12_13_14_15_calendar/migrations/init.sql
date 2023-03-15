CREATE SEQUENCE app_event_id_seq INCREMENT BY 1 MINVALUE 1 START 1;
CREATE TABLE app_event
(
    id                    INT                                                      NOT NULL,
    title                 VARCHAR(128)                                             NOT NULL,
    begin_date            TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    end_date              TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT NULL,
    description           TEXT                           DEFAULT NULL,
    owner_id              INT                                                      NOT NULL,
    notification_sent     BOOLEAN                        DEFAULT FALSE             NOT NULL,
    notification_received BOOLEAN                        DEFAULT FALSE             NOT NULL,
    PRIMARY KEY (id)
);
CREATE INDEX IDX_13EE8992166D1F9C ON app_event (id);

