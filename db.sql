DROP VIEW IF EXISTS tmp_project_id_dappradar_intersect_revain;
	CREATE VIEW tmp_project_id_dappradar_intersect_revain
	AS
	(SELECT 
		distinct dappradar.id
	FROM 
			(SELECT * FROM project WHERE src = 'dappradar') AS dappradar
		INNER JOIN
			(SELECT * FROM project WHERE src = 'revain') AS revan
			ON revan.name = dappradar.id
		LEFT JOIN
			(SELECT * FROM chain_list) AS chainList
			ON chainlist.chainname = dappradar.chainname
	WHERE 
		(	(dappradar.category = 'defi' AND revan.category = 'Crypto Exchanges')
		OR
			(dappradar.category = 'exchanges' AND revan.category = 'Crypto Exchanges')
		OR
			(dappradar.category = 'games' AND revan.category = 'Blockchain Games')
		OR
			(dappradar.category = 'marketplaces' AND revan.category = 'NFT Marketplaces')
		)
	ORDER BY dappradar.id
);