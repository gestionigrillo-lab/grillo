-- Dati di esempio per demo/sviluppo

-- Condomini
INSERT INTO condomini (nome, indirizzo, citta, cap, num_unita, telefono, email) VALUES
('Condominio Villa Rosa', 'Via Roma 15', 'Milano', '20121', 24, '02-1234567', 'villarosa@email.it'),
('Condominio Sole', 'Via Garibaldi 42', 'Milano', '20124', 18, '02-7654321', 'sole@email.it'),
('Residenza Parco Verde', 'Viale dei Tigli 8', 'Milano', '20133', 36, '02-9876543', 'parcoverde@email.it'),
('Condominio Belvedere', 'Corso Buenos Aires 120', 'Milano', '20129', 30, '02-1112233', 'belvedere@email.it'),
('Condominio San Marco', 'Via Torino 55', 'Milano', '20123', 12, '02-4455667', 'sanmarco@email.it')
ON CONFLICT DO NOTHING;

-- Contatti
INSERT INTO contatti (condominio_id, nome, ruolo, telefono, email, interno) VALUES
(1, 'Mario Rossi', 'portinere', '333-1111111', 'mario.rossi@email.it', NULL),
(1, 'Anna Bianchi', 'condomino', '333-2222222', 'anna.bianchi@email.it', 'Scala A int. 3'),
(1, 'Giuseppe Verdi', 'condomino', '333-3333333', 'g.verdi@email.it', 'Scala B int. 7'),
(2, 'Franco Neri', 'portinere', '333-4444444', 'f.neri@email.it', NULL),
(2, 'Lucia Gialli', 'condomino', '333-5555555', 'l.gialli@email.it', 'Int. 12'),
(3, 'Roberto Blu', 'condomino', '333-6666666', 'r.blu@email.it', 'Scala C int. 1'),
(3, 'Carla Marroni', 'portinere', '333-7777777', 'c.marroni@email.it', NULL),
(4, 'Paolo Viola', 'condomino', '333-8888888', 'p.viola@email.it', 'Int. 5'),
(5, 'Sara Arancio', 'condomino', '333-9999999', 's.arancio@email.it', 'Int. 2')
ON CONFLICT DO NOTHING;

-- Fornitori
INSERT INTO fornitori (nome, tipo, telefono, telefono_reperibilita, email, piva, reperibile_h24, zone_coperte) VALUES
('Idraulica Rapida SRL', 'idraulico', '02-5551234', '335-1001001', 'info@idraulicarapida.it', '12345678901', TRUE, 'Milano centro, zona 1, zona 2, zona 3'),
('Elettro Service', 'elettricista', '02-5554321', '335-2002002', 'info@elettroservice.it', '12345678902', TRUE, 'Milano, hinterland nord'),
('Fabbro Express', 'fabbro', '02-5559876', '335-3003003', 'info@fabbroexpress.it', '12345678903', TRUE, 'Milano e provincia'),
('Ascensori Milano SPA', 'ascensorista', '02-5556789', '335-4004004', 'urgenze@ascensorimi.it', '12345678904', TRUE, 'Milano e provincia'),
('Edilstrutture SRL', 'muratore', '02-5551111', NULL, 'info@edilstrutture.it', '12345678905', FALSE, 'Milano sud'),
('Vetri & Serramenti', 'vetraio', '02-5552222', NULL, 'info@vetrieserr.it', '12345678906', FALSE, 'Milano'),
('Caldaie Sicure', 'caldaista', '02-5553333', '335-7007007', 'info@caldaiesicure.it', '12345678907', TRUE, 'Milano, Monza, Brianza'),
('Puliservice Coop', 'impresa_pulizie', '02-5554444', NULL, 'info@puliservice.it', '12345678908', FALSE, 'Milano e provincia'),
('Green Garden', 'giardiniere', '02-5555555', NULL, 'info@greengarden.it', '12345678909', FALSE, 'Milano nord')
ON CONFLICT DO NOTHING;

-- Ticket di esempio
INSERT INTO tickets (numero, condominio_id, categoria, urgenza, stato, oggetto, descrizione, richiedente_nome, richiedente_ruolo) VALUES
('TK-20260315-0001', 1, 'idraulica', 'media', 'in_lavorazione', 'Perdita acqua bagno comune piano terra', 'Segnalata perdita acqua dal soffitto del bagno comune al piano terra. Probabile origine dal bagno dell''appartamento soprastante.', 'Mario Rossi', 'portinere'),
('TK-20260315-0002', 2, 'elettrica', 'ordinaria', 'aperto', 'Luce scale piano 3 non funzionante', 'Le luci del terzo piano scala B non si accendono da 3 giorni.', 'Franco Neri', 'portinere'),
('TK-20260314-0003', 3, 'ascensore', 'urgente', 'in_lavorazione', 'Ascensore bloccato al piano 5', 'Ascensore bloccato al quinto piano senza persone all''interno. Display errore E04.', 'Carla Marroni', 'portinere')
ON CONFLICT DO NOTHING;

-- Eventi ticket
INSERT INTO ticket_events (ticket_id, tipo, descrizione, autore) VALUES
(1, 'creato', 'Ticket creato via chiamata portinere. Urgenza media.', 'assistente_ai'),
(1, 'fornitore_contattato', 'Contattato Idraulica Rapida SRL per sopralluogo.', 'assistente_ai'),
(2, 'creato', 'Ticket creato via chiamata portinere. Manutenzione ordinaria.', 'assistente_ai'),
(3, 'creato', 'Ticket creato via chiamata portinere. Urgente: ascensore bloccato.', 'assistente_ai'),
(3, 'fornitore_contattato', 'Contattato Ascensori Milano SPA su numero reperibilità.', 'assistente_ai')
ON CONFLICT DO NOTHING;
