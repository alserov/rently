CREATE TABLE IF NOT EXISTS rents
(
    rent_uuid       varchar(40) NOT NULL,
    car_uuid        varchar(40) NOT NULL,
    user_uuid       varchar(40),
    phone_number    varchar(25) NOT NULL,
    passport_number varchar(16) NOT NULL,
    rent_start      timestamp   NOT NULL,
    rent_end        timestamp   NOT NULL
);

CREATE TABLE IF NOT EXISTS charges
(
    rent_uuid     varchar(40)                       NOT NULL,
    charge_uuid   text                              NOT NULL,
    charge_amount float CHECK ( charge_amount > 0 ) NOT NULL
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

CREATE INDEX idx_cars_uuid_pagination ON cars (uuid);
CREATE INDEX idx_cars_brand_pagination ON cars (brand);
CREATE INDEX idx_cars_category_pagination ON cars (category);