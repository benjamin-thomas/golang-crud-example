--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: touch(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION touch() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        NEW.updated_on := now();
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.touch() OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE addresses (
    id integer NOT NULL,
    city_id integer NOT NULL,
    name character varying(100) NOT NULL,
    created_on timestamp without time zone DEFAULT now() NOT NULL,
    updated_on timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE addresses OWNER TO postgres;

--
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE addresses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE addresses_id_seq OWNER TO postgres;

--
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE addresses_id_seq OWNED BY addresses.id;


--
-- Name: cities; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE cities (
    id integer NOT NULL,
    country_id integer NOT NULL,
    name character varying(100) NOT NULL,
    created_on timestamp without time zone DEFAULT now() NOT NULL,
    updated_on timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE cities OWNER TO postgres;

--
-- Name: cities_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE cities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE cities_id_seq OWNER TO postgres;

--
-- Name: cities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE cities_id_seq OWNED BY cities.id;


--
-- Name: countries; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE countries (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    created_on timestamp without time zone DEFAULT now() NOT NULL,
    updated_on timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE countries OWNER TO postgres;

--
-- Name: countries_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE countries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE countries_id_seq OWNER TO postgres;

--
-- Name: countries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE countries_id_seq OWNED BY countries.id;


--
-- Name: migrations; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE migrations (
    version integer NOT NULL
);


ALTER TABLE migrations OWNER TO postgres;

--
-- Name: zip_codes; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE zip_codes (
    id integer NOT NULL,
    city_id integer NOT NULL,
    code character varying(50),
    created_on timestamp without time zone DEFAULT now() NOT NULL,
    updated_on timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE zip_codes OWNER TO postgres;

--
-- Name: zip_codes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE zip_codes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE zip_codes_id_seq OWNER TO postgres;

--
-- Name: zip_codes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE zip_codes_id_seq OWNED BY zip_codes.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY addresses ALTER COLUMN id SET DEFAULT nextval('addresses_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY cities ALTER COLUMN id SET DEFAULT nextval('cities_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY countries ALTER COLUMN id SET DEFAULT nextval('countries_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY zip_codes ALTER COLUMN id SET DEFAULT nextval('zip_codes_id_seq'::regclass);


--
-- Name: addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: cities_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY cities
    ADD CONSTRAINT cities_pkey PRIMARY KEY (id);


--
-- Name: countries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (id);


--
-- Name: zip_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY zip_codes
    ADD CONSTRAINT zip_codes_pkey PRIMARY KEY (id);


--
-- Name: addresses_name_uidx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX addresses_name_uidx ON addresses USING btree (name);


--
-- Name: cities_name_uidx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX cities_name_uidx ON cities USING btree (name);


--
-- Name: countries_name_uidx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX countries_name_uidx ON countries USING btree (name);


--
-- Name: migrations_version_uidx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX migrations_version_uidx ON migrations USING btree (version);


--
-- Name: zip_codes_city_id_code_uidx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX zip_codes_city_id_code_uidx ON zip_codes USING btree (city_id, code);


--
-- Name: addresses_touch_trg; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER addresses_touch_trg BEFORE UPDATE ON addresses FOR EACH ROW EXECUTE PROCEDURE touch();


--
-- Name: cities_touch_trg; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER cities_touch_trg BEFORE UPDATE ON cities FOR EACH ROW EXECUTE PROCEDURE touch();


--
-- Name: countries_touch_trg; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER countries_touch_trg BEFORE UPDATE ON countries FOR EACH ROW EXECUTE PROCEDURE touch();


--
-- Name: zip_codes_touch_trg; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER zip_codes_touch_trg BEFORE UPDATE ON zip_codes FOR EACH ROW EXECUTE PROCEDURE touch();


--
-- Name: addresses_city_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY addresses
    ADD CONSTRAINT addresses_city_id_fkey FOREIGN KEY (city_id) REFERENCES cities(id);


--
-- Name: cities_country_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY cities
    ADD CONSTRAINT cities_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries(id);


--
-- Name: zip_codes_city_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY zip_codes
    ADD CONSTRAINT zip_codes_city_id_fkey FOREIGN KEY (city_id) REFERENCES cities(id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

