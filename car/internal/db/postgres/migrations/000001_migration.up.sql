CREATE TABLE IF NOT EXISTS rents
(
    rent_uuid        varchar(40) NOT NULL,
    rent_price       float,
    car_uuid         varchar(40) NOT NULL,
    phone_number     varchar(25) NOT NULL,
    passport_number  varchar(16) NOT NULL,
    card_credentials varchar(40) NOT NULL,
    rent_start       timestamp   NOT NULL,
    rent_end         timestamp   NOT NULL
)