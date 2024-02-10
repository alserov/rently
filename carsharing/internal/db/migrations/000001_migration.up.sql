CREATE TABLE IF NOT EXISTS rents
(
    uuid            varchar(40) NOT NULL,
    car_uuid        varchar(40) NOT NULL,
    user_uuid       varchar(40) DEFAULT 'NOT AUTHORIZED',
    phone_number    varchar(25) NOT NULL,
    passport_number varchar(16) NOT NULL,
    rent_start      timestamp   NOT NULL,
    rent_end        timestamp   NOT NULL
);

CREATE TABLE IF NOT EXISTS charges
(
    uuid          text                              NOT NULL,
    rent_uuid     varchar(40)                       NOT NULL,
    charge_amount float CHECK ( charge_amount > 0 ) NOT NULL,
    status        varchar(10)
);

CREATE TABLE IF NOT EXISTS cars
(
    uuid          text                              NOT NULL,
    brand         varchar(30)                       NOT NULL,
    type          varchar(20)                       NOT NULL,
    max_speed     int                               NOT NULL,
    seats         int CHECK ( seats > 1 )           NOT NULL,
    category      varchar(10)                       NOT NULL,
    price_per_day float CHECK ( price_per_day > 0 ) NOT NULL,
    image_uuid    text                              NOT NULL
);

CREATE TABLE IF NOT EXISTS images
(
    uuid     text NOT NULL,
    car_uuid text NOT NULL
);

CREATE INDEX idx_cars_uuid ON cars (uuid);
CREATE INDEX idx_cars_brand ON cars (brand);
CREATE INDEX idx_cars_category ON cars (category);