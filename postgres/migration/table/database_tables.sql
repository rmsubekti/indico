create table if not exists "user" (
    id serial primary key,
    "name" varchar(100) not null,
    email varchar(120) unique not null,
    pass varchar(120) not null,
    "role" varchar(120) not null
);

create table if not exists product (
    id serial primary key,
    "name" varchar(100) not null,
    sku varchar(120) not null,
    qty int not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    deleted_at timestamptz    
);

create table if not exists warehouse (
    id serial primary key,
    "name" varchar(100) not null,
    "address" varchar(120) not null,
    capacity int not null
);

create table if not exists stock (
    id serial primary key,
    warehouse_id int not null,
    product_id int not null,
    quantity int not null,
    constraint fk_warehouse_stock foreign key (warehouse_id) 
    references warehouse(id),
    constraint fk_stock_product foreign key (product_id) 
    references product(id)
);

create table if not exists "order" (
    id serial primary key,
    from_warehouse_id int references warehouse(id),
    to_warehouse_id int references warehouse(id),
    "type" varchar(8) not null,
    "status" varchar(20) not null,
    "note" varchar(200) not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    deleted_at timestamptz
);

create table if not exists order_detail (
    id serial primary key,
    order_id int not null,
    product_id int not null,
    quantity int not null,
    constraint fk_order_detail foreign key (order_id) 
    references "order"(id),
    constraint fk_product_order foreign key (product_id) 
    references product(id)
);
