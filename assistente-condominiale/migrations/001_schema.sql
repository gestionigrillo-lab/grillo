-- Schema iniziale: Assistente Condominiale Ufficio
-- Gestioni Grillo

CREATE TABLE IF NOT EXISTS condomini (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    indirizzo VARCHAR(500) NOT NULL,
    citta VARCHAR(100) NOT NULL DEFAULT '',
    cap VARCHAR(10) NOT NULL DEFAULT '',
    codice_fiscale VARCHAR(20),
    num_unita INTEGER NOT NULL DEFAULT 0,
    telefono VARCHAR(30),
    email VARCHAR(255),
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contatti (
    id SERIAL PRIMARY KEY,
    condominio_id INTEGER NOT NULL REFERENCES condomini(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    ruolo VARCHAR(100) NOT NULL,
    telefono VARCHAR(30),
    email VARCHAR(255),
    interno VARCHAR(20),
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fornitori (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    tipo VARCHAR(100) NOT NULL,
    telefono VARCHAR(30) NOT NULL,
    telefono_reperibilita VARCHAR(30),
    email VARCHAR(255),
    piva VARCHAR(20),
    indirizzo VARCHAR(500),
    zone_coperte TEXT,
    attivo BOOLEAN NOT NULL DEFAULT TRUE,
    reperibile_h24 BOOLEAN NOT NULL DEFAULT FALSE,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chiamate (
    id SERIAL PRIMARY KEY,
    direzione VARCHAR(10) NOT NULL,
    numero_chiamante VARCHAR(30),
    numero_chiamato VARCHAR(30),
    nome_chiamante VARCHAR(255),
    ruolo_chiamante VARCHAR(100),
    condominio_riferimento VARCHAR(255),
    condominio_id INTEGER,
    urgenza VARCHAR(30),
    categoria VARCHAR(100),
    problema_sintetico VARCHAR(500),
    ticket_creato VARCHAR(30),
    fornitore_contattato VARCHAR(255),
    escalation BOOLEAN NOT NULL DEFAULT FALSE,
    esito TEXT,
    trascrizione TEXT,
    durata_secondi INTEGER,
    inizio TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    fine TIMESTAMPTZ,
    twilio_call_sid VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    numero VARCHAR(30) UNIQUE NOT NULL,
    condominio_id INTEGER NOT NULL REFERENCES condomini(id),
    fornitore_id INTEGER REFERENCES fornitori(id),
    chiamata_id INTEGER REFERENCES chiamate(id),
    categoria VARCHAR(100) NOT NULL,
    urgenza VARCHAR(30) NOT NULL,
    stato VARCHAR(30) NOT NULL DEFAULT 'aperto',
    oggetto VARCHAR(500) NOT NULL,
    descrizione TEXT NOT NULL,
    richiedente_nome VARCHAR(255),
    richiedente_ruolo VARCHAR(100),
    richiedente_telefono VARCHAR(30),
    indirizzo_preciso VARCHAR(500),
    interno VARCHAR(20),
    fornitore_contattato BOOLEAN NOT NULL DEFAULT FALSE,
    fornitore_esito TEXT,
    escalation BOOLEAN NOT NULL DEFAULT FALSE,
    escalation_motivo TEXT,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ticket_events (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    tipo VARCHAR(50) NOT NULL,
    descrizione TEXT NOT NULL,
    autore VARCHAR(100) NOT NULL DEFAULT 'assistente_ai',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indici
CREATE INDEX IF NOT EXISTS idx_condomini_nome ON condomini(nome);
CREATE INDEX IF NOT EXISTS idx_contatti_condominio ON contatti(condominio_id);
CREATE INDEX IF NOT EXISTS idx_tickets_stato ON tickets(stato);
CREATE INDEX IF NOT EXISTS idx_tickets_urgenza ON tickets(urgenza);
CREATE INDEX IF NOT EXISTS idx_tickets_condominio ON tickets(condominio_id);
CREATE INDEX IF NOT EXISTS idx_ticket_events_ticket ON ticket_events(ticket_id);
CREATE INDEX IF NOT EXISTS idx_chiamate_inizio ON chiamate(inizio);
