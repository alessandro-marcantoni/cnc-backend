-- Remove Armadietti Normali slots 62 to 80
DELETE FROM facilities
WHERE facility_type_id = 6
  AND identifier IN ('62','63','64','65','66','67','68','69','70',
                     '71','72','73','74','75','76','77','78','79','80');

-- Remove Deposito SUP Chiuso slots 11 to 20
DELETE FROM facilities
WHERE facility_type_id = 8
  AND identifier IN ('11','12','13','14','15','16','17','18','19','20');
