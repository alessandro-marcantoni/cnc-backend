DELETE FROM facility_pricing_rules;

INSERT INTO facility_pricing_rules (facility_type_id, required_facility_type_id, special_price, description)
VALUES (4, 3, 50.00, 'Prezzo speciale per Rastrelliera Tavole Chiusa con Box'),
       (7, 3, 50.00, 'Prezzo speciale per Deposito Surf Lato Officina con Box'),
       (9, 3, 50.00, 'Prezzo speciale per Canoe Rastrelliere Esterne con Box'),
       (1, 3, 50.00, 'Prezzo speciale per Rastrelliera Tavole Aperta con Box'),
       (8, 3, 50.00, 'Prezzo speciale per Deposito SUP Chiuso con Box'),
       (4, 5, 50.00, 'Prezzo speciale per Rastrelliera Tavole Chiusa con Armadietti Grandi Legno'),
       (7, 5, 50.00, 'Prezzo speciale per Deposito Surf Lato Officina con Armadietti Grandi Legno'),
       (9, 5, 50.00, 'Prezzo speciale per Canoe Rastrelliere Esterne con Armadietti Grandi Legno'),
       (1, 5, 50.00, 'Prezzo speciale per Rastrelliera Tavole Aperta con Armadietti Grandi Legno'),
       (8, 5, 50.00, 'Prezzo speciale per Deposito SUP Chiuso con Armadietti Grandi Legno');
