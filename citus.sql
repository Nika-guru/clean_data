SELECT create_reference_table('product_category');
SELECT create_reference_table('sub_category');
SELECT create_reference_table('category');

SELECT create_reference_table('product_blockchain');
SELECT create_reference_table('blockchain');

SELECT create_reference_table('product_tag');
SELECT create_reference_table('tag');


SELECT create_distributed_table('product_contact', 'productid');
SELECT create_distributed_table('product', 'productid');
SELECT create_distributed_table('product_statistic', 'productid');

SELECT create_distributed_table('review', 'productid');
SELECT create_distributed_table('reply', 'productid');
SELECT create_distributed_table('account_reaction', 'productid');
SELECT create_distributed_table('account', 'userid');
