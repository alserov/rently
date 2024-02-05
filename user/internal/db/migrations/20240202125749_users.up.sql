CREATE TABLE IF NOT EXISTS users
(
    uuid            VARCHAR(40) NOT NULL PRIMARY KEY,
    username        TEXT        NOT NULL,
    password        VARCHAR(100) NOT NULL,
    role            VARCHAR(5)  NOT NULL,
    email           TEXT        NOT NULL,
    passport_number VARCHAR(40) NOT NULL,
    payment_source  VARCHAR(40) NOT NULL,
    phone_number    VARCHAR(16) NOT NULL
);
