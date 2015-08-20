SELECT a.id AS addr_id
     , a.name AS addr_name
     , a.line1
     , a.line2
     , a.line3
     , c.id AS city_id
     , c.name AS city_name
     , z.id AS zip_code_id
     , z.code AS zip_code
     , cn.id AS country_id
     , cn.name AS country_name
  FROM addresses AS a
 INNER
  JOIN cities AS c
    ON c.id  = a.city_id
 INNER
  JOIN countries AS cn
    ON cn.id = c.country_id

    LEFT JOIN zip_codes AS z
    ON z.city_id = c.id
