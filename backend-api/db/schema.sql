\restrict dbmate

-- Dumped from database version 18.4 (Debian 18.4-1.pgdg13+1)
-- Dumped by pg_dump version 18.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: principal_login_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.principal_login_events (
    id character varying NOT NULL,
    principal character varying DEFAULT ''::character varying NOT NULL,
    tenant character varying DEFAULT ''::character varying NOT NULL,
    ip_address inet NOT NULL,
    user_agent character varying DEFAULT ''::character varying NOT NULL,
    method character varying NOT NULL,
    success boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: principal_sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.principal_sessions (
    id character varying NOT NULL,
    principal character varying DEFAULT ''::character varying NOT NULL,
    ip_address inet NOT NULL,
    user_agent character varying DEFAULT ''::character varying NOT NULL,
    tenant character varying DEFAULT ''::character varying NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
    expires_at timestamp with time zone DEFAULT (now() + '01:00:00'::interval) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: principal_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.principal_tokens (
    id character varying NOT NULL,
    principal character varying DEFAULT ''::character varying NOT NULL,
    kind character varying DEFAULT ''::character varying NOT NULL,
    tenant character varying DEFAULT ''::character varying NOT NULL,
    expires_at timestamp with time zone DEFAULT (now() + '00:05:00'::interval) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: principal_login_events principal_login_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.principal_login_events
    ADD CONSTRAINT principal_login_events_pkey PRIMARY KEY (id);


--
-- Name: principal_sessions principal_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.principal_sessions
    ADD CONSTRAINT principal_sessions_pkey PRIMARY KEY (id);


--
-- Name: principal_tokens principal_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.principal_tokens
    ADD CONSTRAINT principal_tokens_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- PostgreSQL database dump complete
--

\unrestrict dbmate


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20240612000001');
