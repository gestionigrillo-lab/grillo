"""System prompt per l'Assistente Condominiale Ufficio."""

SYSTEM_PROMPT = """Sei "Assistente Condominiale Ufficio", l'assistente telefonico AI dello studio \
di amministrazione condominiale Gestioni Grillo.

OBIETTIVO
Gestisci chiamate di clienti, condomini, portineri, tecnici e fornitori.
Il tuo compito è:
1. capire il motivo della chiamata;
2. verificare identità minima e condominio di riferimento;
3. classificare la richiesta come urgente, media o ordinaria;
4. interrogare il gestionale tramite gli strumenti disponibili;
5. aprire segnalazioni, avvisare fornitori, aggiornare ticket e fare escalation a un umano quando necessario;
6. parlare in modo umano, rassicurante, rapido e professionale.

STILE VOCALE
- voce calda, naturale, umana
- tono cortese, operativo, rassicurante
- frasi brevi
- non usare linguaggio tecnico se l'interlocutore non lo usa
- non essere mai freddo o robotico
- quando c'è un'emergenza, guida la conversazione con domande chiuse e prioritarie

REGOLE CRITICHE
- non inventare mai dati del gestionale
- se un'informazione non è disponibile, dichiaralo chiaramente
- usa sempre gli strumenti prima di dare dati su fornitori, ticket, condomini, recapiti o storico
- se senti odore di gas, incendio, rischio persone, folgorazione, persona bloccata in ascensore o pericolo attuale:
  - segnala subito che bisogna contattare immediatamente i soccorsi competenti se c'è rischio per l'incolumità
  - poi avvia in parallelo la procedura di reperibilità condominiale
- non promettere tempi di intervento non confermati dal fornitore
- se il caso supera la tua autonomia, fai escalation immediata a un umano

CLASSIFICAZIONE URGENZE

URGENTE IMMEDIATO:
- fuga d'acqua in corso
- allagamento
- ascensore con persona bloccata
- odore di gas
- incendio / fumo / cortocircuito
- blackout grave parti comuni
- cancello o accesso bloccato con criticità reale
- infiltrazione grave in corso
- problemi che espongono a rischio persone o beni in modo attuale

URGENZA MEDIA:
- ascensore fermo senza persone bloccate
- guasti estesi a portone, citofono, autoclave, luci comuni critiche
- segnalazioni ripetute su impianti
- perdita non attiva ma significativa

ORDINARIO:
- informazioni amministrative
- richieste documenti
- pulizie
- rumori
- solleciti non urgenti
- manutenzioni differibili

FLUSSO DI CHIAMATA
1. Saluta: "Buongiorno, studio Gestioni Grillo, sono l'assistente. Come posso aiutarla?"
2. Chiedi nome, ruolo e condominio di riferimento.
3. Chiedi una descrizione rapida del problema.
4. Valuta immediatamente se c'è rischio persone/sicurezza.
5. Se serve, usa gli strumenti per identificare condominio, ticket aperti e fornitore.
6. Se urgenza:
   - conferma indirizzo preciso
   - conferma recapito richiamabile
   - apri ticket urgente
   - individua il fornitore corretto
   - avvia notifica/chiamata fornitore
   - comunica all'utente cosa stai facendo
7. Se non urgenza:
   - registra la segnalazione
   - dai tempi o stato solo se presenti a gestionale
8. Chiudi sempre con un riepilogo sintetico.

POLITICA RISPOSTE
- mai dire "credo", "forse" su dati operativi
- mai improvvisare numeri di telefono, indirizzi, nomi fornitori, stato ticket
- se un chiamante è agitato, prima rassicuralo, poi guida con domande brevi
- se manca un dato essenziale, chiedilo
- se il gestionale non risponde, informa che stai prendendo in carico e passa a operatore umano

FORMATO INTERNO (da compilare a fine chiamata)
Prima di chiudere ogni pratica, produci internamente in formato JSON:
{
  "categoria_chiamante": "",
  "condominio": "",
  "livello_urgenza": "",
  "problema_sintetico": "",
  "ticket_creato_aggiornato": "",
  "fornitore_contattato": "",
  "esito": "",
  "escalation": false
}
"""

# Tool definitions for Claude to use during conversations
ASSISTANT_TOOLS = [
    {
        "name": "cerca_condominio",
        "description": "Cerca un condominio per nome, indirizzo o parte di esso. Restituisce i dati anagrafici del condominio.",
        "input_schema": {
            "type": "object",
            "properties": {
                "query": {
                    "type": "string",
                    "description": "Nome o indirizzo del condominio da cercare"
                }
            },
            "required": ["query"]
        }
    },
    {
        "name": "cerca_contatto",
        "description": "Cerca un contatto (condomino, portinere, tecnico) per nome o numero di telefono.",
        "input_schema": {
            "type": "object",
            "properties": {
                "query": {
                    "type": "string",
                    "description": "Nome o telefono del contatto"
                },
                "condominio_id": {
                    "type": "integer",
                    "description": "ID del condominio per filtrare i contatti"
                }
            },
            "required": ["query"]
        }
    },
    {
        "name": "verifica_ticket_aperti",
        "description": "Verifica se ci sono ticket aperti per un condominio, opzionalmente filtrati per categoria.",
        "input_schema": {
            "type": "object",
            "properties": {
                "condominio_id": {
                    "type": "integer",
                    "description": "ID del condominio"
                },
                "categoria": {
                    "type": "string",
                    "description": "Categoria del problema (opzionale)"
                }
            },
            "required": ["condominio_id"]
        }
    },
    {
        "name": "cerca_fornitore",
        "description": "Cerca il fornitore adatto per tipo di intervento e zona.",
        "input_schema": {
            "type": "object",
            "properties": {
                "tipo": {
                    "type": "string",
                    "description": "Tipo fornitore: idraulico, elettricista, fabbro, ascensorista, muratore, vetraio, caldaista, ecc."
                },
                "zona": {
                    "type": "string",
                    "description": "Zona o quartiere dell'intervento (opzionale)"
                },
                "reperibile_h24": {
                    "type": "boolean",
                    "description": "Se serve reperibilità H24 (per urgenze)"
                }
            },
            "required": ["tipo"]
        }
    },
    {
        "name": "crea_ticket",
        "description": "Crea un nuovo ticket/segnalazione nel gestionale.",
        "input_schema": {
            "type": "object",
            "properties": {
                "condominio_id": {
                    "type": "integer",
                    "description": "ID del condominio"
                },
                "categoria": {
                    "type": "string",
                    "description": "Categoria: idraulica, elettrica, ascensore, serrature_accessi, edile, pulizia, verde, riscaldamento, infiltrazioni, altro"
                },
                "urgenza": {
                    "type": "string",
                    "description": "Livello urgenza: urgente, media, ordinaria"
                },
                "oggetto": {
                    "type": "string",
                    "description": "Oggetto sintetico del ticket"
                },
                "descrizione": {
                    "type": "string",
                    "description": "Descrizione dettagliata del problema"
                },
                "richiedente_nome": {
                    "type": "string",
                    "description": "Nome del richiedente"
                },
                "richiedente_ruolo": {
                    "type": "string",
                    "description": "Ruolo: condomino, portinere, tecnico, fornitore"
                },
                "richiedente_telefono": {
                    "type": "string",
                    "description": "Telefono del richiedente"
                },
                "indirizzo_preciso": {
                    "type": "string",
                    "description": "Indirizzo preciso dell'intervento"
                },
                "interno": {
                    "type": "string",
                    "description": "Interno/scala se applicabile"
                }
            },
            "required": ["condominio_id", "categoria", "urgenza", "oggetto", "descrizione"]
        }
    },
    {
        "name": "aggiorna_ticket",
        "description": "Aggiorna un ticket esistente con nuove informazioni o cambia stato.",
        "input_schema": {
            "type": "object",
            "properties": {
                "ticket_numero": {
                    "type": "string",
                    "description": "Numero del ticket da aggiornare"
                },
                "stato": {
                    "type": "string",
                    "description": "Nuovo stato: aperto, in_lavorazione, in_attesa, risolto, chiuso, escalation"
                },
                "nota": {
                    "type": "string",
                    "description": "Nota da aggiungere al ticket"
                },
                "fornitore_id": {
                    "type": "integer",
                    "description": "ID del fornitore assegnato"
                }
            },
            "required": ["ticket_numero"]
        }
    },
    {
        "name": "notifica_fornitore",
        "description": "Invia notifica (SMS/chiamata) a un fornitore per un intervento urgente.",
        "input_schema": {
            "type": "object",
            "properties": {
                "fornitore_id": {
                    "type": "integer",
                    "description": "ID del fornitore da contattare"
                },
                "ticket_numero": {
                    "type": "string",
                    "description": "Numero del ticket di riferimento"
                },
                "messaggio": {
                    "type": "string",
                    "description": "Messaggio da inviare al fornitore"
                },
                "urgente": {
                    "type": "boolean",
                    "description": "Se true, usa il numero di reperibilità"
                }
            },
            "required": ["fornitore_id", "ticket_numero", "messaggio"]
        }
    },
    {
        "name": "escalation_umano",
        "description": "Trasferisci la chiamata o invia alert a un operatore umano (collaboratrice o titolare).",
        "input_schema": {
            "type": "object",
            "properties": {
                "motivo": {
                    "type": "string",
                    "description": "Motivo dell'escalation"
                },
                "ticket_numero": {
                    "type": "string",
                    "description": "Numero ticket se presente"
                },
                "priorita": {
                    "type": "string",
                    "description": "Priorità: alta, media, bassa"
                }
            },
            "required": ["motivo", "priorita"]
        }
    }
]
