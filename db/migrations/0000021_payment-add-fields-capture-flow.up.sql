alter table payment add column payment_flow varchar(255) not null default 'autocapture';
alter table payment add column authorized_at timestamp;
alter table payment add column captured_at timestamp;
alter table payment add column released_at timestamp;