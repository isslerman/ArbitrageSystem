--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orderHistory (
    id VARCHAR(255) PRIMARY KEY,
    spread DOUBLE PRECISION,
    aexcid VARCHAR(255), 
    aprice DOUBLE PRECISION, 
    apricevet DOUBLE PRECISION, 
    avolume DOUBLE PRECISION, 
    bexcid VARCHAR(255), 
    bprice DOUBLE PRECISION, 
    bpricevet DOUBLE PRECISION, 
    bvolume DOUBLE PRECISION, 
    created_at integer,
    updated_at integer,
);