create table lead
(
    id                        bigserial primary key,
    site                      varchar                  not null,
    name                      varchar                  not null default '',
    re_stage_constructor      varchar                  not null,
    re_region                 varchar                  not null,
    re_type                   varchar                  not null,
    re_purpose_of_acquisition varchar                  not null,
    re_count_of_rooms         varchar                  not null,
    purchase_budget           varchar                  not null,
    phone                     varchar                  not null,
    email                     varchar                  not null default '',
    communication_method      varchar                  not null,
    description               varchar                  not null default '',
    status                    varchar                  not null default 'new',

    created_at                timestamp with time zone not null default now(),
    updated_at                timestamp with time zone not null default now()
);

alter table lead
    add name varchar default '' not null;

create table admin
(
    id         bigserial                not null primary key,
    login      varchar                  not null,
    password   varchar                  not null,
    role       varchar                  not null default 'admin',
    status     varchar                  not null default 'enabled',
    token      varchar unique           not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create table client
(
    id         bigserial                not null primary key,
    first_name varchar                  not null,
    last_name  varchar                  not null,
    phone      varchar                  not null,
    login      varchar                  not null unique,
    password   varchar                  not null,
    token      varchar unique           not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create table comment
(
    id          bigserial                not null primary key,
    lead_id     bigint                   not null,
    admin_id    bigint                   not null,
    admin_login varchar                  not null,
    text        varchar                  not null,
    created_at  timestamp with time zone not null default now()
);

alter table comment
    add admin_login varchar default '' not null;

create table estate
(
    id                bigserial                not null primary key,
    status            varchar                  not null default 'new',
    luxury            boolean                  not null default false,
    images            text[] not null default {''},
    created_at        timestamp with time zone not null default now(),
    updated_at        timestamp with time zone not null default now(),
    --description
    name              varchar                  not null,
    price             int                      not null default 0,
    country           int                      not null,
    city              int                      not null,
    address           varchar                  not null,
    beds              int                      not null,
    baths             int                      not null,
    area_in_meter     int                      not null,
    property_type     int                      not null,
    year_built        int                      not null,
    description       varchar                  not null,
    latitude          varchar                  not null,
    longitude         varchar                  not null,
    --interior
    appliances        varchar                  not null,
    interior_features varchar                  not null,
    kitchen_features  varchar                  not null,
    total_bedrooms    int                      not null,
    full_bathrooms    int                      not null,
    half_bathrooms    int                      not null,
    floor_description varchar                  not null,
    fireplace         varchar                  not null,
    cooling           int                      not null default 0, -- bool 0 - null, 1 - false, 2 - true
    heating           int                      not null default 0, -- bool 0 - null, 1 - false, 2 - true
    --exterior
    lot_size_in_acres int                      not null,
    exterior_features varchar                  not null,
    arch_style        varchar                  not null,
    roof              varchar                  not null,
    sewer             varchar                  not null,
    --other
    area_name         varchar                  not null,
    garage            int                      not null,
    parking           varchar                  not null,
    view              varchar                  not null,
    pool              int                      not null default 0, -- bool 0 - null, 1 - false, 2 - true
    pool_description  varchar                  not null,
    water_source      varchar                  not null,
    utilities         varchar                  not null
);

create table text
(
    key   varchar not null primary key,
    value varchar not null
);

create table favorite
(
    id        bigserial                not null primary key,
    user_id   bigint                   not null,
    estate_id bigint                   not null,
    create_at timestamp with time zone not null default now()
);

create table availability
(
    id           bigserial                not null primary key,
    lending_id   integer                  not null,
    price_aed    integer                  not null,
    price_usd    integer                  not null,
    unique_id    varchar                  not null,
    bedroom      integer                  not null,
    parking      integer                  not null,
    area         varchar                  not null,
    plot         varchar                  not null,
    special_gift varchar                  not null,
    created_at   timestamp with time zone not null default now(),
    updated_at   timestamp with time zone not null default now()
);

create table lending
(
    id                     bigserial                not null primary key,
    name                   varchar                  not null,
    main_description       varchar                  not null,
    full_name              varchar                  not null,
    slogan                 varchar                  not null,
    address                varchar                  not null,
    starting_price_aed     integer                  not null,
    starting_price_usd     integer                  not null,
    property_type          varchar                  not null,
    furnishing             varchar                  not null,
    features_and_amenities integer[] not null,
    title                  varchar                  not null,
    description            varchar                  not null,
    video                  varchar                  not null,
    file_plan              varchar                  not null default '',
    title_plan             varchar                  not null,
    background_image       varchar                  not null default '',
    background_for_mobile  varchar                  not null default '',
    main_logo              varchar                  not null default '',
    partner_logo           varchar                  not null default '',
    our_logo               varchar                  not null default '',
    images                 text[] not null,
    latitude               varchar                  not null,
    longitude              varchar                  not null,
    created_at             timestamp with time zone not null default now(),
    updated_at             timestamp with time zone not null default now()
);

create table feature_or_amenity
(
    id   bigserial not null primary key,
    name varchar   not null,
    icon varchar   not null
);

create table service_keys
(
    id         serial  not null primary key,
    name       varchar not null,
    key        varchar not null unique,
    lending_id integer not null
);
