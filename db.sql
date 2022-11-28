-----------------------review service---------------------
drop table if exists review;
create table review (
	id bigserial not null,
	accountId bigint null,
	content varchar null,
	star bigint null,
	productId bigint,
	createddate timestamp NOT NULL,
    updateddate timestamp NOT NULL
);
drop index if exists idx_review_productInfoId;
create index idx_review_productInfoId on review using btree(productId); --get list
drop index if exists idx_review_star;
create index idx_review_star on review using btree(star); --sorting

drop table if exists account_reaction;
create table account_reaction(
	commentId bigint not null,
	type int8 not null,
	accountId bigint not null,
	username varchar  null,
	reactionType bigint not null, --like, dislike, trash
	productId bigint,
	createddate timestamp NOT NULL,
    updateddate timestamp NOT NULL
);
drop index if exists idx_account_reaction_reactionType;
create index idx_account_reaction_reactionType on account_reaction using hash(reactionType); --get one
drop index if exists idx_account_reaction_commentId;
create index idx_account_reaction_replyId on account_reaction using btree(commentId); --get list

drop table if exists reply;
create table reply (
	id  bigserial not null,
	reviewId bigint null,
	accountid bigint null,
	productId bigint,
	content varchar null,
	createddate timestamp NOT NULL,
    updateddate timestamp NOT NULL
);
drop index if exists idx_reply_reviewId;
create index idx_reply_reviewId on reply using btree(reviewId); --get list

drop table if exists account;
create table account(
	id bigserial not null,
	userId bigint null,
	role int8 null,
	password varchar,
     image varchar null,
	accountType varchar null, --facebook, twitter, github, ...
	email varchar null,
	username varchar null,
	createddate timestamp NOT NULL,
    updateddate timestamp NOT NULL
);

----------------------------product info service------------------------------------
drop table if exists product;
create table product(
	productId bigserial not null,
	productName varchar not null,
	productImage varchar null,
	productDescription text null,
	productDetail jsonb null,
	crawlSource varchar null,
	createdDate timestamp not null,
	updatedDate timestamp not null
);
drop index if exists idx_product_product_id;
create index idx_product_product_id on product using hash(productId); --get one
-------------

drop table if exists category;
create table category(
	id bigserial not null,
	name varchar not null
);

drop table if exists sub_category;
create table sub_category(
	id bigserial not null,
	name varchar not null,
	categoryId bigint not null
);
create index idx_sub_category_id on sub_category using hash(categoryId); --get list




drop table if exists product_category;
create table product_category(
	productId bigint not null,
	categoryId bigint not null,
	subCategoryId bigint null,
	createdDate timestamp not null,
	updatedDate timestamp not null
);
create index idx_product_category_product_id on product_category using btree(productId); --get list


drop table if exists product_contact;
create table product_contact(
	productId bigint not null,
	url varchar null,
	type varchar null,
	createdDate timestamp not null,
	updatedDate timestamp not null
);
create index idx_product_contact_product_id on product_contact using btree(productId); --get list

drop table if exists product_statistic;
create table product_statistic(
	productId bigint not null, 
	
	totalReviews bigint null,
	averageStar int8 null,
	
	price varchar null,
	holder bigint null,
	marketcap varchar null,
	volume varchar null,
	tvl varchar null,
	totalUsed varchar  null,
	source varchar null,
	createdDate timestamp not null,
	updatedDate timestamp not null
);
create index idx_product_statistic_product_id on product_statistic using hash(productId); --get one
create index idx_product_statistic_tvl on product_statistic using btree(tvl); --sorting
create index idx_product_statistic_price on product_statistic using btree(price); --soring
create index idx_product_statistic_holder on product_statistic using btree(holder); --soring
create index idx_product_statistic_marketcap on product_statistic using btree(marketcap); --soring
create index idx_product_statistic_volume on product_statistic using btree(volume); --soring
create index idx_product_statistic_totalUsed on product_statistic using btree(totalUsed); --soring


create table product_raw_category(
	productId bigint not null,
	prodcutCategories varchar
);


------------------------------remove all data in all table
truncate table account, account_reaction, category, product, product_category, product_contact, product_raw_category, product_statistic, reply, review, sub_category restart identity;
