--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.5
-- Dumped by pg_dump version 9.6.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: arp_table_entries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE arp_table_entries (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    mac_addr macaddr NOT NULL,
    ip_addr inet NOT NULL,
    table_id uuid NOT NULL
);


ALTER TABLE arp_table_entries OWNER TO postgres;

--
-- Name: arp_tables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE arp_tables (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    switch_id uuid NOT NULL
);


ALTER TABLE arp_tables OWNER TO postgres;

--
-- Name: communities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE communities (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    value text NOT NULL
);


ALTER TABLE communities OWNER TO postgres;

--
-- Name: mac_addr_table_entries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE mac_addr_table_entries (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    mac_addr macaddr NOT NULL,
    port_id uuid NOT NULL,
    vlan_id uuid NOT NULL,
    table_id uuid NOT NULL
);


ALTER TABLE mac_addr_table_entries OWNER TO postgres;

--
-- Name: mac_addr_tables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE mac_addr_tables (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    switch_id uuid NOT NULL
);


ALTER TABLE mac_addr_tables OWNER TO postgres;

--
-- Name: ports; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE ports (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    name text NOT NULL,
    number integer NOT NULL,
    switch_id uuid NOT NULL
);


ALTER TABLE ports OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE schema_migration (
    version character varying(255) NOT NULL
);


ALTER TABLE schema_migration OWNER TO postgres;

--
-- Name: switch_models; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE switch_models (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    name text NOT NULL,
    vendor text NOT NULL
);


ALTER TABLE switch_models OWNER TO postgres;

--
-- Name: switches; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE switches (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    name text NOT NULL,
    location text NOT NULL,
    ip_addr text NOT NULL,
    model_id uuid NOT NULL,
    community_id uuid NOT NULL
);


ALTER TABLE switches OWNER TO postgres;

--
-- Name: vlans; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE vlans (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    id uuid NOT NULL,
    name text NOT NULL,
    number integer NOT NULL,
    deleted boolean NOT NULL,
    switch_id uuid NOT NULL
);


ALTER TABLE vlans OWNER TO postgres;

--
-- Name: arp_table_entries arp_table_entries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY arp_table_entries
    ADD CONSTRAINT arp_table_entries_pkey PRIMARY KEY (id);


--
-- Name: arp_tables arp_tables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY arp_tables
    ADD CONSTRAINT arp_tables_pkey PRIMARY KEY (id);


--
-- Name: communities communities_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY communities
    ADD CONSTRAINT communities_pkey PRIMARY KEY (id);


--
-- Name: mac_addr_table_entries mac_addr_table_entries_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY mac_addr_table_entries
    ADD CONSTRAINT mac_addr_table_entries_pkey PRIMARY KEY (id);


--
-- Name: mac_addr_tables mac_addr_tables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY mac_addr_tables
    ADD CONSTRAINT mac_addr_tables_pkey PRIMARY KEY (id);


--
-- Name: ports ports_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY ports
    ADD CONSTRAINT ports_pkey PRIMARY KEY (id);


--
-- Name: switch_models switch_models_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY switch_models
    ADD CONSTRAINT switch_models_pkey PRIMARY KEY (id);


--
-- Name: switches switches_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY switches
    ADD CONSTRAINT switches_pkey PRIMARY KEY (id);


--
-- Name: vlans vlans_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY vlans
    ADD CONSTRAINT vlans_pkey PRIMARY KEY (id);


--
-- Name: version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX version_idx ON schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

