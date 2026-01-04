-- =======================================
-- MEMBERS
-- =======================================
INSERT INTO members (first_name, last_name, date_of_birth, email)
SELECT
    first_name, last_name, date_of_birth, email
FROM (VALUES
('Luca','Rossi',DATE '1985-03-12','luca.rossi@example.com'),
('Giulia','Bianchi',DATE '1990-07-22','giulia.bianchi@example.com'),
('Marco','Verdi',DATE '1978-11-05','marco.verdi@example.com'),
('Francesca','Neri',DATE '1995-02-18','francesca.neri@example.com'),
('Alessandro','Ferrari',DATE '1982-05-09','alessandro.ferrari@example.com'),
('Martina','Romano',DATE '1988-12-01','martina.romano@example.com'),
('Davide','Conti',DATE '1991-04-15','davide.conti@example.com'),
('Sara','Greco',DATE '1992-10-20','sara.greco@example.com'),
('Simone','Galli',DATE '1983-07-11','simone.galli@example.com'),
('Elena','Marini',DATE '1994-03-03','elena.marini@example.com'),
('Paolo','Rinaldi',DATE '1979-08-25','paolo.rinaldi@example.com'),
('Chiara','Lombardi',DATE '1993-11-17','chiara.lombardi@example.com'),
('Matteo','Moretti',DATE '1987-02-05','matteo.moretti@example.com'),
('Valentina','Costa',DATE '1990-09-12','valentina.costa@example.com'),
('Antonio','Giordano',DATE '1980-01-30','antonio.giordano@example.com'),
('Federica','Corsi',DATE '1989-06-06','federica.corsi@example.com'),
('Giorgio','Parisi',DATE '1985-04-19','giorgio.parisi@example.com'),
('Laura','Pellegrini',DATE '1991-07-23','laura.pellegrini@example.com'),
('Stefano','Fabbri',DATE '1982-12-14','stefano.fabbri@example.com'),
('Monica','Bellini',DATE '1995-05-02','monica.bellini@example.com'),
('Riccardo','Barbieri',DATE '1984-03-25','riccardo.barbieri@example.com'),
('Ilaria','Martini',DATE '1990-10-30','ilaria.martini@example.com'),
('Gabriele','Riva',DATE '1986-06-18','gabriele.riva@example.com'),
('Francesca','De Luca',DATE '1992-09-05','francesca.deluca@example.com'),
('Lorenzo','Sartori',DATE '1983-08-14','lorenzo.sartori@example.com'),
('Giovanna','Bruni',DATE '1991-11-29','giovanna.bruni@example.com'),
('Enrico','Mancini',DATE '1980-12-10','enrico.mancini@example.com'),
('Elisa','Vitale',DATE '1994-01-22','elisa.vitale@example.com'),
('Daniele','Bianco',DATE '1987-05-19','daniele.bianco@example.com'),
('Roberta','Serra',DATE '1993-02-16','roberta.serra@example.com'),
('Nicola','Pagani',DATE '1985-09-03','nicola.pagani@example.com'),
('Veronica','Fontana',DATE '1989-03-28','veronica.fontana@example.com'),
('Fabio','Capri',DATE '1982-07-07','fabio.capri@example.com'),
('Silvia','Landi',DATE '1990-12-21','silvia.landi@example.com'),
('Michele','Donati',DATE '1983-10-01','michele.donati@example.com'),
('Valeria','Gatti',DATE '1991-05-11','valeria.gatti@example.com'),
('Alberto','Grassi',DATE '1986-08-23','alberto.grassi@example.com'),
('Simona','Pace',DATE '1988-04-04','simona.pace@example.com'),
('Claudio','Ruggeri',DATE '1980-11-15','claudio.ruggeri@example.com'),
('Giada','Serafini',DATE '1995-07-12','giada.serafini@example.com'),
('Emanuele','Fiorentino',DATE '1984-01-29','emanuele.fiorentino@example.com'),
('Martina','Riva',DATE '1989-09-09','martina.riva@example.com'),
('Alessio','Ferraro',DATE '1987-03-05','alessio.ferraro@example.com'),
('Elena','Moro',DATE '1992-06-20','elena.moro@example.com'),
('Roberto','Grillo',DATE '1981-12-18','roberto.grillo@example.com'),
('Paola','Costa',DATE '1985-02-14','paola.costa@example.com'),
('Filippo','Bertoni',DATE '1983-08-30','filippo.bertoni@example.com'),
('Anna','Biagi',DATE '1990-11-02','anna.biagi@example.com'),
('Lorenzo','Cattaneo',DATE '1984-07-25','lorenzo.cattaneo@example.com'),
('Giulia','Martino',DATE '1991-01-08','giulia.martino@example.com'),
('Riccardo','Rossi',DATE '1982-10-11','riccardo.rossi@example.com'),
('Sara','Ferri',DATE '1986-05-14','sara.ferri@example.com')
) AS t(first_name,last_name,date_of_birth,email);

-- =======================================
-- PHONE NUMBERS
-- =======================================
INSERT INTO phone_numbers (member_id, number, description)
SELECT id, '333' || (1000000 + id), 'mobile' FROM members;

-- =======================================
-- ADDRESSES
-- =======================================
INSERT INTO addresses (member_id, country, city, street, street_number)
SELECT id, 'Italy', CASE WHEN id%4=1 THEN 'Rome' WHEN id%4=2 THEN 'Milan' WHEN id%4=3 THEN 'Florence' ELSE 'Naples' END,
       'Street ' || id, id::text
FROM members;

-- =======================================
-- FACILITIES CATALOG
-- =======================================
INSERT INTO facilities_catalog (name, description, suggested_price) VALUES
('Tennis Court','Outdoor clay court',25.00),
('Swimming Pool','Indoor 25m pool',15.00),
('Boat Dock','Dock for small boats',50.00);

-- =======================================
-- FACILITIES
-- =======================================
INSERT INTO facilities (facility_type_id, identifier)
SELECT 1,'TC-'||id FROM generate_series(1,10) AS id
UNION ALL
SELECT 2,'SP-'||id FROM generate_series(1,5) AS id
UNION ALL
SELECT 3,'BD-'||id FROM generate_series(1,3) AS id;

-- =======================================
-- RENTED FACILITIES
-- =======================================
INSERT INTO rented_facilities (facility_id, member_id, rented_at, expires_at)
SELECT
  (id%18)+1,
  id,
  now() - ((id%30) || ' days')::interval,
  now() + (((id%10)+1) || ' days')::interval
FROM members;

-- =======================================
-- MEMBERSHIP STATUSES
-- =======================================
INSERT INTO membership_statuses (status) VALUES
('active'),('expired'),('suspended');

-- =======================================
-- MEMBERSHIPS
-- =======================================
INSERT INTO memberships (member_id, number, created_at)
SELECT id, 10000+id, now() - ((id%365) || ' days')::interval
FROM members;

-- =======================================
-- MEMBERSHIP PERIODS
-- =======================================
INSERT INTO membership_periods (membership_id, valid_from, expires_at, status_id)
SELECT id,
       now() - ((id%400) || ' days')::interval,
       now() + ((id%100) || ' days')::interval,
       CASE WHEN id%5=0 THEN 2 -- expired
            WHEN id%7=0 THEN 3 -- suspended
            ELSE 1 END
FROM memberships;

-- =======================================
-- PAYMENTS
-- =======================================
INSERT INTO payments (membership_period_id, amount, payment_method)
SELECT id,
       300.00,
       CASE WHEN id%4=0 THEN 'cash' WHEN id%4=1 THEN 'credit card' WHEN id%4=2 THEN 'paypal' ELSE 'bank transfer' END
FROM membership_periods
WHERE id%6 != 0; -- leave some unpaid

-- =======================================
-- BOATS
-- =======================================
INSERT INTO boats (rented_facility_id, name, length_meters, width_meters)
SELECT id,'Boat-'||id, 8.0+(id%5), 2.5+(id%2) FROM rented_facilities WHERE facility_id%3=0;

-- =======================================
-- INSURANCES
-- =======================================
INSERT INTO insurances (boat_id, provider, number, expires_at)
SELECT
    id,
    'Generali',
    'INS-' || id,
    now() + (((id % 365) + 30) * interval '1 day')
FROM boats;

-- =======================================
-- WAITING LIST
-- =======================================
INSERT INTO members_waiting (member_id, facility_type_id, queued_at)
SELECT id, ((id%3)+1), now() - ((id%15) || ' days')::interval
FROM members
WHERE id%4=0;
