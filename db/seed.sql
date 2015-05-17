-- The following insert statements assume a fully migrated database

SELECT 'db/seed.sql' AS filepath;

INSERT INTO countries (name) VALUES
  ('Afghanistan'),
  ('Albania'),
  ('Algeria'),
  ('Andorra'),
  ('Angola'),
  ('Antigua & Deps'),
  ('Argentina'),
  ('Armenia'),
  ('Australia'),
  ('Austria'),
  ('Azerbaijan'),
  ('Bahamas'),
  ('Bahrain'),
  ('Bangladesh'),
  ('Barbados'),
  ('Belarus'),
  ('Belgium'),
  ('Belize'),
  ('Benin'),
  ('Bhutan'),
  ('Bolivia'),
  ('Bosnia Herzegovina'),
  ('Botswana'),
  ('Brazil'),
  ('Brunei'),
  ('Bulgaria'),
  ('Burkina'),
  ('Burundi'),
  ('Cambodia'),
  ('Cameroon'),
  ('Canada'),
  ('Cape Verde'),
  ('Central African Rep'),
  ('Chad'),
  ('Chile'),
  ('China'),
  ('Colombia'),
  ('Comoros'),
  ('Congo'),
  ('Congo {Democratic Rep}'),
  ('Costa Rica'),
  ('Croatia'),
  ('Cuba'),
  ('Cyprus'),
  ('Czech Republic'),
  ('Denmark'),
  ('Djibouti'),
  ('Dominica'),
  ('Dominican Republic'),
  ('East Timor'),
  ('Ecuador'),
  ('Egypt'),
  ('El Salvador'),
  ('Equatorial Guinea'),
  ('Eritrea'),
  ('Estonia'),
  ('Ethiopia'),
  ('Fiji'),
  ('Finland'),
  ('France'),
  ('Gabon'),
  ('Gambia'),
  ('Georgia'),
  ('Germany'),
  ('Ghana'),
  ('Greece'),
  ('Grenada'),
  ('Guatemala'),
  ('Guinea'),
  ('Guinea-Bissau'),
  ('Guyana'),
  ('Haiti'),
  ('Honduras'),
  ('Hungary'),
  ('Iceland'),
  ('India'),
  ('Indonesia'),
  ('Iran'),
  ('Iraq'),
  ('Ireland {Republic}'),
  ('Israel'),
  ('Italy'),
  ('Ivory Coast'),
  ('Jamaica'),
  ('Japan'),
  ('Jordan'),
  ('Kazakhstan'),
  ('Kenya'),
  ('Kiribati'),
  ('Korea North'),
  ('Korea South'),
  ('Kosovo'),
  ('Kuwait'),
  ('Kyrgyzstan'),
  ('Laos'),
  ('Latvia'),
  ('Lebanon'),
  ('Lesotho'),
  ('Liberia'),
  ('Libya'),
  ('Liechtenstein'),
  ('Lithuania'),
  ('Luxembourg'),
  ('Macedonia'),
  ('Madagascar'),
  ('Malawi'),
  ('Malaysia'),
  ('Maldives'),
  ('Mali'),
  ('Malta'),
  ('Marshall Islands'),
  ('Mauritania'),
  ('Mauritius'),
  ('Mexico'),
  ('Micronesia'),
  ('Moldova'),
  ('Monaco'),
  ('Mongolia'),
  ('Montenegro'),
  ('Morocco'),
  ('Mozambique'),
  ('Myanmar, {Burma}'),
  ('Namibia'),
  ('Nauru'),
  ('Nepal'),
  ('Netherlands'),
  ('New Zealand'),
  ('Nicaragua'),
  ('Niger'),
  ('Nigeria'),
  ('Norway'),
  ('Oman'),
  ('Pakistan'),
  ('Palau'),
  ('Panama'),
  ('Papua New Guinea'),
  ('Paraguay'),
  ('Peru'),
  ('Philippines'),
  ('Poland'),
  ('Portugal'),
  ('Qatar'),
  ('Romania'),
  ('Russian Federation'),
  ('Rwanda'),
  ('St Kitts & Nevis'),
  ('St Lucia'),
  ('Saint Vincent & the Grenadines'),
  ('Samoa'),
  ('San Marino'),
  ('Sao Tome & Principe'),
  ('Saudi Arabia'),
  ('Senegal'),
  ('Serbia'),
  ('Seychelles'),
  ('Sierra Leone'),
  ('Singapore'),
  ('Slovakia'),
  ('Slovenia'),
  ('Solomon Islands'),
  ('Somalia'),
  ('South Africa'),
  ('South Sudan'),
  ('Spain'),
  ('Sri Lanka'),
  ('Sudan'),
  ('Suriname'),
  ('Swaziland'),
  ('Sweden'),
  ('Switzerland'),
  ('Syria'),
  ('Taiwan'),
  ('Tajikistan'),
  ('Tanzania'),
  ('Thailand'),
  ('Togo'),
  ('Tonga'),
  ('Trinidad & Tobago'),
  ('Tunisia'),
  ('Turkey'),
  ('Turkmenistan'),
  ('Tuvalu'),
  ('Uganda'),
  ('Ukraine'),
  ('United Arab Emirates'),
  ('United Kingdom'),
  ('United States'),
  ('Uruguay'),
  ('Uzbekistan'),
  ('Vanuatu'),
  ('Vatican City'),
  ('Venezuela'),
  ('Vietnam'),
  ('Yemen'),
  ('Zambia'),
  ('Zimbabwe');

DO $$
  DECLARE france_id INT;
  DECLARE uk_id INT;
BEGIN
  SELECT id FROM countries WHERE name = 'France' LIMIT 1 INTO france_id;
  SELECT id FROM countries WHERE name = 'United Kingdom' LIMIT 1 INTO uk_id;

INSERT INTO cities (country_id, name) VALUES
  (france_id, 'Paris'),
  (france_id, 'Nice'),
  (france_id, 'Marseille'),
  (france_id, 'Venelles'),
  (france_id, 'Aix-en-Provence'),
  (uk_id, 'London'),
  (uk_id, 'Oxford'),
  (uk_id, 'Brighton');

INSERT INTO country_stats (country_id, population_count) VALUES
  (france_id, 66808074),
  (uk_id, 64812393);

END $$;

DO $$
  DECLARE paris_id INT;
  DECLARE london_id INT;
BEGIN
  SELECT id FROM cities WHERE name = 'Paris' LIMIT 1 INTO paris_id;
  SELECT id FROM cities WHERE name = 'London' LIMIT 1 INTO london_id;

INSERT INTO zip_codes (city_id, code) VALUES
  (paris_id, '75017'),
  (paris_id, '75020'),
  (london_id, 'W1D 1LA');

INSERT INTO addresses (city_id, name) VALUES
  (paris_id, 'Paris addr 1'),
  (paris_id, 'Paris addr 2'),
  (london_id, 'London addr 1'),
  (london_id, 'London addr 2');
END $$;
