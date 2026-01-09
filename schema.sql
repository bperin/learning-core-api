--
-- PostgreSQL database dump
--

\restrict 5mHDx9R1dhUHgq4vpJKaGakPr2rbul4dxm6OA6DEQVPwTyLf8cbeSNYrqRc0wL3

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1 (Homebrew)

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

--
-- Name: artifact_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.artifact_type AS ENUM (
    'INTENTS',
    'PLAN',
    'EVAL',
    'EVAL_ITEM',
    'PROMPT',
    'PROMPT_POLICY',
    'QUALITY_METRICS',
    'HINT',
    'SUMMARY',
    'OUTLINE',
    'OTHER'
);


--
-- Name: TYPE artifact_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.artifact_type IS 'Types of artifacts that can be generated in the system';


--
-- Name: review_verdict; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.review_verdict AS ENUM (
    'APPROVED',
    'REJECTED',
    'NEEDS_REVISION'
);


--
-- Name: TYPE review_verdict; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.review_verdict IS 'Possible verdicts for eval item reviews';


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;


--
-- Name: FUNCTION update_updated_at_column(); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.update_updated_at_column() IS 'Trigger function to automatically update updated_at timestamp';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: artifacts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.artifacts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    type public.artifact_type NOT NULL,
    status text DEFAULT 'READY'::text NOT NULL,
    eval_id uuid,
    eval_item_id uuid,
    attempt_id uuid,
    reviewer_id uuid,
    text text,
    output_json jsonb,
    model text,
    prompt text,
    input_hash text,
    meta jsonb,
    error text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    prompt_template_id uuid,
    schema_template_id uuid,
    model_params jsonb,
    prompt_render text
);


--
-- Name: COLUMN artifacts.reviewer_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.artifacts.reviewer_id IS 'User who reviewed or created the artifact';


--
-- Name: COLUMN artifacts.output_json; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.artifacts.output_json IS 'Structured output payload from the model';


--
-- Name: COLUMN artifacts.prompt_template_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.artifacts.prompt_template_id IS 'Prompt template used for generation';


--
-- Name: COLUMN artifacts.schema_template_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.artifacts.schema_template_id IS 'Schema template used to validate output';


--
-- Name: COLUMN artifacts.model_params; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.artifacts.model_params IS 'Model parameters (temperature, top_p, etc.)';


--
-- Name: COLUMN artifacts.prompt_render; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.artifacts.prompt_render IS 'Rendered prompt text sent to the model';


--
-- Name: curricula; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.curricula (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    subject_id uuid NOT NULL,
    parent_id uuid,
    label text NOT NULL,
    code text,
    description text,
    order_index integer DEFAULT 0 NOT NULL,
    grade_level text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: COLUMN curricula.subject_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.curricula.subject_id IS 'Subject this curriculum belongs to';


--
-- Name: COLUMN curricula.parent_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.curricula.parent_id IS 'Parent curriculum node; NULL is root';


--
-- Name: COLUMN curricula.label; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.curricula.label IS 'Human-readable curriculum label';


--
-- Name: COLUMN curricula.code; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.curricula.code IS 'Optional standardized curriculum code';


--
-- Name: COLUMN curricula.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.curricula.description IS 'Optional curriculum description';


--
-- Name: COLUMN curricula.order_index; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.curricula.order_index IS 'Sort order within parent curriculum';


--
-- Name: COLUMN curricula.grade_level; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.curricula.grade_level IS 'Grade or level alignment';


--
-- Name: COLUMN curricula.is_active; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.curricula.is_active IS 'Whether curriculum unit is active';


--
-- Name: documents; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.documents (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    filename text NOT NULL,
    title text,
    mime_type text,
    content text,
    storage_path text,
    rag_status text DEFAULT 'PENDING'::text NOT NULL,
    user_id uuid NOT NULL,
    subject_id uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    curricular text,
    subjects text[] DEFAULT '{}'::text[],
    curriculum_id uuid
);


--
-- Name: COLUMN documents.curricular; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.documents.curricular IS 'Curricular classification or framework (e.g., Common Core, IB)';


--
-- Name: COLUMN documents.subjects; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.documents.subjects IS 'List of academic subjects associated with this document';


--
-- Name: COLUMN documents.curriculum_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.documents.curriculum_id IS 'Curriculum unit associated with this document';


--
-- Name: eval_item_reviews; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.eval_item_reviews (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    eval_item_id uuid NOT NULL,
    reviewer_id uuid NOT NULL,
    verdict public.review_verdict NOT NULL,
    reasons text[] NOT NULL,
    comments text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: COLUMN eval_item_reviews.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.eval_item_reviews.updated_at IS 'When the review was last updated';


--
-- Name: eval_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.eval_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    eval_id uuid NOT NULL,
    prompt text NOT NULL,
    options text[] NOT NULL,
    correct_idx integer NOT NULL,
    hint text,
    explanation text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: COLUMN eval_items.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.eval_items.created_at IS 'When the eval item was created';


--
-- Name: COLUMN eval_items.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.eval_items.updated_at IS 'When the eval item was last updated';


--
-- Name: evals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.evals (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    description text,
    status text DEFAULT 'draft'::text NOT NULL,
    difficulty text,
    instructions text,
    rubric jsonb,
    subject_id uuid,
    user_id uuid NOT NULL,
    published_at timestamp with time zone,
    archived_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.goose_db_version ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: prompt_templates; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.prompt_templates (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    key text NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    title text NOT NULL,
    description text,
    template text NOT NULL,
    metadata jsonb,
    created_by text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: schema_templates; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_templates (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    schema_type text NOT NULL,
    version integer NOT NULL,
    schema_json jsonb NOT NULL,
    subject_id uuid,
    curriculum_id uuid,
    is_active boolean DEFAULT true NOT NULL,
    created_by uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    locked_at timestamp with time zone
);


--
-- Name: COLUMN schema_templates.schema_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.schema_templates.schema_type IS 'Schema purpose (eval_generation, intent_extraction, etc.)';


--
-- Name: COLUMN schema_templates.schema_json; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.schema_templates.schema_json IS 'JSON schema defining expected AI output';


--
-- Name: COLUMN schema_templates.locked_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.schema_templates.locked_at IS 'Timestamp when version becomes immutable';


--
-- Name: subjects; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subjects (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    description text,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: test_attempts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.test_attempts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    eval_id uuid NOT NULL,
    score integer DEFAULT 0 NOT NULL,
    total integer NOT NULL,
    percentage real,
    total_time integer,
    feedback jsonb,
    summary text,
    started_at timestamp with time zone DEFAULT now() NOT NULL,
    completed_at timestamp with time zone
);


--
-- Name: user_answers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_answers (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    attempt_id uuid NOT NULL,
    eval_item_id uuid NOT NULL,
    selected_idx integer NOT NULL,
    is_correct boolean NOT NULL,
    time_spent integer,
    hints_used integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: COLUMN user_answers.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_answers.updated_at IS 'When the user answer was last updated';


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    is_admin boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    is_learner boolean DEFAULT true NOT NULL,
    is_teacher boolean DEFAULT false NOT NULL
);


--
-- Name: COLUMN users.is_admin; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.is_admin IS 'True if the user has administrative privileges';


--
-- Name: COLUMN users.is_learner; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.is_learner IS 'True if the user is a learner/student';


--
-- Name: COLUMN users.is_teacher; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.is_teacher IS 'True if the user is a teacher/instructor';


--
-- Name: artifacts artifacts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.artifacts
    ADD CONSTRAINT artifacts_pkey PRIMARY KEY (id);


--
-- Name: curricula curricula_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.curricula
    ADD CONSTRAINT curricula_pkey PRIMARY KEY (id);


--
-- Name: documents documents_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_pkey PRIMARY KEY (id);


--
-- Name: eval_item_reviews eval_item_reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eval_item_reviews
    ADD CONSTRAINT eval_item_reviews_pkey PRIMARY KEY (id);


--
-- Name: eval_items eval_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eval_items
    ADD CONSTRAINT eval_items_pkey PRIMARY KEY (id);


--
-- Name: evals evals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.evals
    ADD CONSTRAINT evals_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: prompt_templates prompt_templates_key_version_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.prompt_templates
    ADD CONSTRAINT prompt_templates_key_version_key UNIQUE (key, version);


--
-- Name: prompt_templates prompt_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.prompt_templates
    ADD CONSTRAINT prompt_templates_pkey PRIMARY KEY (id);


--
-- Name: schema_templates schema_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_templates
    ADD CONSTRAINT schema_templates_pkey PRIMARY KEY (id);


--
-- Name: schema_templates schema_templates_schema_type_version_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_templates
    ADD CONSTRAINT schema_templates_schema_type_version_key UNIQUE (schema_type, version);


--
-- Name: subjects subjects_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_pkey PRIMARY KEY (id);


--
-- Name: test_attempts test_attempts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.test_attempts
    ADD CONSTRAINT test_attempts_pkey PRIMARY KEY (id);


--
-- Name: user_answers user_answers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_answers
    ADD CONSTRAINT user_answers_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_artifacts_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_artifacts_type ON public.artifacts USING btree (type);


--
-- Name: idx_curricula_parent; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_curricula_parent ON public.curricula USING btree (parent_id);


--
-- Name: idx_curricula_subject; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_curricula_subject ON public.curricula USING btree (subject_id);


--
-- Name: idx_documents_curriculum; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_documents_curriculum ON public.documents USING btree (curriculum_id);


--
-- Name: idx_documents_subject; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_documents_subject ON public.documents USING btree (subject_id);


--
-- Name: idx_documents_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_documents_user ON public.documents USING btree (user_id);


--
-- Name: idx_eval_item_reviews_item; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_eval_item_reviews_item ON public.eval_item_reviews USING btree (eval_item_id);


--
-- Name: idx_eval_item_reviews_reviewer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_eval_item_reviews_reviewer ON public.eval_item_reviews USING btree (reviewer_id);


--
-- Name: idx_eval_item_reviews_updated_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_eval_item_reviews_updated_at ON public.eval_item_reviews USING btree (updated_at);


--
-- Name: idx_eval_items_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_eval_items_created_at ON public.eval_items USING btree (created_at);


--
-- Name: idx_eval_items_eval; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_eval_items_eval ON public.eval_items USING btree (eval_id);


--
-- Name: idx_eval_items_updated_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_eval_items_updated_at ON public.eval_items USING btree (updated_at);


--
-- Name: idx_evals_subject; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_evals_subject ON public.evals USING btree (subject_id);


--
-- Name: idx_evals_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_evals_user ON public.evals USING btree (user_id);


--
-- Name: idx_prompt_templates_key_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_prompt_templates_key_active ON public.prompt_templates USING btree (key, is_active);


--
-- Name: idx_schema_templates_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_schema_templates_active ON public.schema_templates USING btree (is_active);


--
-- Name: idx_schema_templates_curriculum; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_schema_templates_curriculum ON public.schema_templates USING btree (curriculum_id);


--
-- Name: idx_schema_templates_subject; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_schema_templates_subject ON public.schema_templates USING btree (subject_id);


--
-- Name: idx_schema_templates_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_schema_templates_type ON public.schema_templates USING btree (schema_type);


--
-- Name: idx_subjects_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subjects_user ON public.subjects USING btree (user_id);


--
-- Name: idx_test_attempts_eval; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_test_attempts_eval ON public.test_attempts USING btree (eval_id);


--
-- Name: idx_test_attempts_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_test_attempts_user ON public.test_attempts USING btree (user_id);


--
-- Name: idx_user_answers_attempt; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_answers_attempt ON public.user_answers USING btree (attempt_id);


--
-- Name: idx_user_answers_eval_item; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_answers_eval_item ON public.user_answers USING btree (eval_item_id);


--
-- Name: idx_user_answers_updated_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_answers_updated_at ON public.user_answers USING btree (updated_at);


--
-- Name: idx_users_is_admin; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_is_admin ON public.users USING btree (is_admin);


--
-- Name: idx_users_is_learner; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_is_learner ON public.users USING btree (is_learner);


--
-- Name: idx_users_is_teacher; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_is_teacher ON public.users USING btree (is_teacher);


--
-- Name: curricula update_curricula_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_curricula_updated_at BEFORE UPDATE ON public.curricula FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: documents update_documents_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_documents_updated_at BEFORE UPDATE ON public.documents FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: eval_item_reviews update_eval_item_reviews_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_eval_item_reviews_updated_at BEFORE UPDATE ON public.eval_item_reviews FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: eval_items update_eval_items_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_eval_items_updated_at BEFORE UPDATE ON public.eval_items FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: evals update_evals_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_evals_updated_at BEFORE UPDATE ON public.evals FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: prompt_templates update_prompt_templates_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_prompt_templates_updated_at BEFORE UPDATE ON public.prompt_templates FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: subjects update_subjects_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_subjects_updated_at BEFORE UPDATE ON public.subjects FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: test_attempts update_test_attempts_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_test_attempts_updated_at BEFORE UPDATE ON public.test_attempts FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: user_answers update_user_answers_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_user_answers_updated_at BEFORE UPDATE ON public.user_answers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: artifacts artifacts_attempt_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.artifacts
    ADD CONSTRAINT artifacts_attempt_id_fkey FOREIGN KEY (attempt_id) REFERENCES public.test_attempts(id) ON DELETE CASCADE;


--
-- Name: artifacts artifacts_eval_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.artifacts
    ADD CONSTRAINT artifacts_eval_id_fkey FOREIGN KEY (eval_id) REFERENCES public.evals(id) ON DELETE CASCADE;


--
-- Name: artifacts artifacts_eval_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.artifacts
    ADD CONSTRAINT artifacts_eval_item_id_fkey FOREIGN KEY (eval_item_id) REFERENCES public.eval_items(id) ON DELETE CASCADE;


--
-- Name: artifacts artifacts_prompt_template_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.artifacts
    ADD CONSTRAINT artifacts_prompt_template_id_fkey FOREIGN KEY (prompt_template_id) REFERENCES public.prompt_templates(id) ON DELETE SET NULL;


--
-- Name: artifacts artifacts_schema_template_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.artifacts
    ADD CONSTRAINT artifacts_schema_template_id_fkey FOREIGN KEY (schema_template_id) REFERENCES public.schema_templates(id) ON DELETE SET NULL;


--
-- Name: artifacts artifacts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.artifacts
    ADD CONSTRAINT artifacts_user_id_fkey FOREIGN KEY (reviewer_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: curricula curricula_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.curricula
    ADD CONSTRAINT curricula_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.curricula(id) ON DELETE SET NULL;


--
-- Name: curricula curricula_subject_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.curricula
    ADD CONSTRAINT curricula_subject_id_fkey FOREIGN KEY (subject_id) REFERENCES public.subjects(id) ON DELETE CASCADE;


--
-- Name: documents documents_curriculum_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_curriculum_id_fkey FOREIGN KEY (curriculum_id) REFERENCES public.curricula(id) ON DELETE SET NULL;


--
-- Name: documents documents_subject_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_subject_id_fkey FOREIGN KEY (subject_id) REFERENCES public.subjects(id) ON DELETE SET NULL;


--
-- Name: documents documents_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: eval_item_reviews eval_item_reviews_eval_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eval_item_reviews
    ADD CONSTRAINT eval_item_reviews_eval_item_id_fkey FOREIGN KEY (eval_item_id) REFERENCES public.eval_items(id) ON DELETE CASCADE;


--
-- Name: eval_item_reviews eval_item_reviews_reviewer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eval_item_reviews
    ADD CONSTRAINT eval_item_reviews_reviewer_id_fkey FOREIGN KEY (reviewer_id) REFERENCES public.users(id);


--
-- Name: eval_items eval_items_eval_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eval_items
    ADD CONSTRAINT eval_items_eval_id_fkey FOREIGN KEY (eval_id) REFERENCES public.evals(id) ON DELETE CASCADE;


--
-- Name: evals evals_subject_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.evals
    ADD CONSTRAINT evals_subject_id_fkey FOREIGN KEY (subject_id) REFERENCES public.subjects(id) ON DELETE SET NULL;


--
-- Name: evals evals_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.evals
    ADD CONSTRAINT evals_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: schema_templates schema_templates_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_templates
    ADD CONSTRAINT schema_templates_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: schema_templates schema_templates_curriculum_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_templates
    ADD CONSTRAINT schema_templates_curriculum_id_fkey FOREIGN KEY (curriculum_id) REFERENCES public.curricula(id) ON DELETE SET NULL;


--
-- Name: schema_templates schema_templates_subject_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_templates
    ADD CONSTRAINT schema_templates_subject_id_fkey FOREIGN KEY (subject_id) REFERENCES public.subjects(id) ON DELETE SET NULL;


--
-- Name: subjects subjects_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: test_attempts test_attempts_eval_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.test_attempts
    ADD CONSTRAINT test_attempts_eval_id_fkey FOREIGN KEY (eval_id) REFERENCES public.evals(id);


--
-- Name: test_attempts test_attempts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.test_attempts
    ADD CONSTRAINT test_attempts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_answers user_answers_attempt_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_answers
    ADD CONSTRAINT user_answers_attempt_id_fkey FOREIGN KEY (attempt_id) REFERENCES public.test_attempts(id) ON DELETE CASCADE;


--
-- Name: user_answers user_answers_eval_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_answers
    ADD CONSTRAINT user_answers_eval_item_id_fkey FOREIGN KEY (eval_item_id) REFERENCES public.eval_items(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: -
--

GRANT ALL ON SCHEMA public TO cloudsqlsuperuser;


--
-- PostgreSQL database dump complete
--

\unrestrict 5mHDx9R1dhUHgq4vpJKaGakPr2rbul4dxm6OA6DEQVPwTyLf8cbeSNYrqRc0wL3

