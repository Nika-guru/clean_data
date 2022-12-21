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
	userId bigint null, --| -1 (updated)
	role int8 null, --| 0: user
	password varchar,
     image varchar null,
	accountType varchar null, --'normal' --facebook, twitter, github, ...
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
drop index if exists idx_sub_category_id;
create index idx_sub_category_id on sub_category using hash(categoryId); --get list




drop table if exists product_category;
create table product_category(
	productId bigint not null,
	categoryId bigint not null,
	subCategoryId bigint null,
	createdDate timestamp not null,
	updatedDate timestamp not null
);
drop index if exists idx_product_category_product_id;
create index idx_product_category_product_id on product_category using btree(productId); --get list


drop table if exists product_contact;
create table product_contact(
	productId bigint not null,
	url varchar null,
	type varchar null,
	createdDate timestamp not null,
	updatedDate timestamp not null
);
drop index if exists idx_product_contact_product_id;
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
drop index if exists idx_product_statistic_product_id;
create index idx_product_statistic_product_id on product_statistic using hash(productId); --get one
drop index if exists idx_product_statistic_tvl;
create index idx_product_statistic_tvl on product_statistic using btree(tvl); --sorting
drop index if exists idx_product_statistic_price;
create index idx_product_statistic_price on product_statistic using btree(price); --soring
drop index if exists idx_product_statistic_holder;
create index idx_product_statistic_holder on product_statistic using btree(holder); --soring
drop index if exists idx_product_statistic_marketcap;
create index idx_product_statistic_marketcap on product_statistic using btree(marketcap); --soring
drop index if exists idx_product_statistic_volume;
create index idx_product_statistic_volume on product_statistic using btree(volume); --soring
drop index if exists idx_product_statistic_totalUsed;
create index idx_product_statistic_totalUsed on product_statistic using btree(totalUsed); --soring


drop table if exists blockchain;
create table blockchain(
	id bigserial not null,
	blockchainId varchar null,
	blockchainName varchar null,
	info jsonb null,
	createdDate timestamp not null,
	updatedDate timestamp not null
);
 
drop table if exists product_blockchain;
create table product_blockchain (
	product_id bigint not null,
	blockchain_id bigint not null,
	createdDate timestamp not null,
	updatedDate timestamp not null
 );
drop index if exists idx_product_blockchain_product_id;
create index idx_product_blockchain_product_id on product_blockchain using btree(product_id); --get list
drop index if exists idx_product_blockchain_blockchain_id;
create index idx_product_blockchain_blockchain_id on product_blockchain using btree(blockchain_id); --get list
 

create table tag(
	id bigserial not null,
	tagName varchar null,
	createdDate timestamp not null,
	updatedDate timestamp not null
 )

create table product_tag(
	product_id bigint not null,
	tag_id bigint not null,
	createdDate timestamp not null,
	updatedDate timestamp not null
)
drop index if exists idx_product_tag_product_id;
create index idx_product_tag_product_id on product_tag using btree(product_id); --get list
drop index if exists idx_product_tag_tag_id;
create index idx_product_tag_tag_id on product_tag using btree(tag_id); --get list


------------------------------remove all data in all table
truncate table account, account_reaction, category, product, product_category, product_contact, product_raw_category, product_statistic, reply, review, sub_category, blockchain restart identity;


drop table if exists product_raw_category;
create table product_raw_category(
	productId bigint not null,
	prodcutCategories varchar
);


-----------------------------raw data
drop table if exists fund_raising;
create table fund_raising(
	projectCode varchar,
	projectName varchar,
	projectLogo varchar,
	
	investorCode varchar,
	investorName varchar,
	investorLogo varchar,

	fundStageCode varchar,
	fundStageName varchar,
	fundAmount float8,
	fundDate timestamp,

	createdDate timestamp,
	updatedDate timestamp,

	description text,
	announcementUrl varchar,
	valulation float8,

	srcFund varchar,
	srcInvestor varchar,
	extradata json
)

drop table if exists investor;
create table investor(
	investorCode varchar,
	investorName varchar,
	investorImage varchar,
	categoryName varchar,
	yearFounded int8,
	location varchar,
	socials json,
	description text,
	src varchar,
	createddate timestamp NOT NULL,
    updateddate timestamp NOT NULL,
	extradata json
);
drop index if exists idx_investor_investor_code;
create index idx_investor_investor_code on investor using btree(investorCode);
drop index if exists idx_investor_investor_name;
create index idx_investor_investor_name on investor using btree(investorName);



drop table if exists project
create table project (
	id varchar, --code
	name varchar,
	category varchar,
	subCategory varchar,
	social json,
	image varchar,
	description varchar,
	chainid varchar,
    chainname varchar,
	extradata json,
    createddate timestamp NOT NULL,
    updateddate timestamp NOT NULL
) ;
create index idx_project_productid on project using hash(id);


drop table if exists coin;
create table coin (
    coinId varchar,
	type varchar, --coin ^ token
	
	address varchar, --lower (NULL: 'NULL')
	chainid varchar,
    chainName varchar,
	
	symbol varchar, --upper
	name varchar,
	tag varchar, 
	totalsupply varchar,
	maxsupply varchar,
	marketcap varchar,
	volumetrading varchar,	

	image varchar,
	decimals int8, -- (NULL: 0)
	src varchar,
	detail json
);
drop index if exists idx_coin_address_address;
create index idx_coin_address_address on coin using btree(address);
drop index if exists idx_coin_address_symbol;
create index idx_coin_address_symbol on coin using btree(symbol);
drop index if exists idx_coin_address_name;
create index idx_coin_address_name on coin using btree(name);
drop index if exists idx_coin_address_chainid;
create index idx_coin_address_chainid on coin using btree(chainid);


