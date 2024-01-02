CREATE TABLE IF NOT EXISTS rents
(
    rent_uuid       varchar(40) NOT NULL,
    rent_price      float       NOT NULL,
    car_uuid        varchar(40) NOT NULL,
    phone_number    varchar(25) NOT NULL,
    passport_number varchar(16) NOT NULL,
    charge_id       varchar(40) NOT NULL,
    rent_start      timestamp   NOT NULL,
    rent_end        timestamp   NOT NULL
);


CREATE TABLE IF NOT EXISTS cars
(
    brand         varchar(30)                     NOT NULL,
    type          varchar(20)                     NOT NULL,
    max_speed     int                             NOT NULL,
    seats         int CHECK ( seats > 1 )         NOT NULL,
    category      varchar(10)                     NOT NULL,
    price_per_day int CHECK ( price_per_day > 0 ) NOT NULL,
    uuid          text                            NOT NULL
);