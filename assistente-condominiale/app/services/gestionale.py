"""Servizio gestionale: interfaccia tra l'assistente AI e il database.

Ogni funzione corrisponde a uno strumento (tool) che Claude può chiamare
durante una conversazione telefonica.
"""

import datetime
import json
from sqlalchemy import select, or_, func
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.models.condominio import Condominio, Contatto
from app.models.ticket import Ticket, TicketEvent
from app.models.fornitore import Fornitore
from app.models.chiamata import Chiamata


def _gen_ticket_numero() -> str:
    """Genera numero ticket: TK-YYYYMMDD-XXXX."""
    now = datetime.datetime.now()
    import random
    seq = random.randint(1000, 9999)
    return f"TK-{now.strftime('%Y%m%d')}-{seq}"


async def cerca_condominio(db: AsyncSession, query: str) -> dict:
    """Cerca condomini per nome o indirizzo (ILIKE)."""
    pattern = f"%{query}%"
    stmt = select(Condominio).where(
        or_(
            Condominio.nome.ilike(pattern),
            Condominio.indirizzo.ilike(pattern),
            Condominio.citta.ilike(pattern),
        )
    ).limit(5)
    result = await db.execute(stmt)
    condomini = result.scalars().all()

    if not condomini:
        return {"trovati": 0, "messaggio": f"Nessun condominio trovato per '{query}'."}

    return {
        "trovati": len(condomini),
        "condomini": [
            {
                "id": c.id,
                "nome": c.nome,
                "indirizzo": c.indirizzo,
                "citta": c.citta,
                "telefono": c.telefono,
                "email": c.email,
                "num_unita": c.num_unita,
            }
            for c in condomini
        ],
    }


async def cerca_contatto(db: AsyncSession, query: str, condominio_id: int | None = None) -> dict:
    """Cerca contatti per nome o telefono, opzionalmente filtrati per condominio."""
    pattern = f"%{query}%"
    stmt = select(Contatto).where(
        or_(
            Contatto.nome.ilike(pattern),
            Contatto.telefono.ilike(pattern),
        )
    )
    if condominio_id:
        stmt = stmt.where(Contatto.condominio_id == condominio_id)
    stmt = stmt.limit(10)

    result = await db.execute(stmt)
    contatti = result.scalars().all()

    if not contatti:
        return {"trovati": 0, "messaggio": f"Nessun contatto trovato per '{query}'."}

    return {
        "trovati": len(contatti),
        "contatti": [
            {
                "id": c.id,
                "nome": c.nome,
                "ruolo": c.ruolo,
                "telefono": c.telefono,
                "email": c.email,
                "interno": c.interno,
                "condominio_id": c.condominio_id,
            }
            for c in contatti
        ],
    }


async def verifica_ticket_aperti(
    db: AsyncSession, condominio_id: int, categoria: str | None = None
) -> dict:
    """Verifica ticket aperti per un condominio."""
    stmt = (
        select(Ticket)
        .where(
            Ticket.condominio_id == condominio_id,
            Ticket.stato.in_(["aperto", "in_lavorazione", "in_attesa"]),
        )
        .order_by(Ticket.created_at.desc())
    )
    if categoria:
        stmt = stmt.where(Ticket.categoria == categoria)
    stmt = stmt.limit(10)

    result = await db.execute(stmt)
    tickets = result.scalars().all()

    return {
        "trovati": len(tickets),
        "tickets": [
            {
                "numero": t.numero,
                "oggetto": t.oggetto,
                "categoria": t.categoria,
                "urgenza": t.urgenza,
                "stato": t.stato,
                "data_apertura": t.created_at.isoformat() if t.created_at else None,
                "fornitore_contattato": t.fornitore_contattato,
            }
            for t in tickets
        ],
    }


async def cerca_fornitore(
    db: AsyncSession, tipo: str, zona: str | None = None, reperibile_h24: bool = False
) -> dict:
    """Cerca fornitore per tipo e zona."""
    stmt = select(Fornitore).where(
        Fornitore.tipo == tipo,
        Fornitore.attivo == True,
    )
    if reperibile_h24:
        stmt = stmt.where(Fornitore.reperibile_h24 == True)
    stmt = stmt.limit(5)

    result = await db.execute(stmt)
    fornitori = result.scalars().all()

    if zona and fornitori:
        # Filtro soft sulle zone coperte
        filtered = [
            f for f in fornitori
            if f.zone_coperte and zona.lower() in f.zone_coperte.lower()
        ]
        if filtered:
            fornitori = filtered

    if not fornitori:
        return {
            "trovati": 0,
            "messaggio": f"Nessun fornitore di tipo '{tipo}' disponibile" +
                         (" con reperibilità H24" if reperibile_h24 else "") + "."
        }

    return {
        "trovati": len(fornitori),
        "fornitori": [
            {
                "id": f.id,
                "nome": f.nome,
                "tipo": f.tipo,
                "telefono": f.telefono,
                "telefono_reperibilita": f.telefono_reperibilita,
                "email": f.email,
                "reperibile_h24": f.reperibile_h24,
                "zone_coperte": f.zone_coperte,
            }
            for f in fornitori
        ],
    }


async def crea_ticket(db: AsyncSession, **kwargs) -> dict:
    """Crea un nuovo ticket nel gestionale."""
    numero = _gen_ticket_numero()
    ticket = Ticket(
        numero=numero,
        condominio_id=kwargs["condominio_id"],
        categoria=kwargs["categoria"],
        urgenza=kwargs["urgenza"],
        stato="aperto",
        oggetto=kwargs["oggetto"],
        descrizione=kwargs["descrizione"],
        richiedente_nome=kwargs.get("richiedente_nome"),
        richiedente_ruolo=kwargs.get("richiedente_ruolo"),
        richiedente_telefono=kwargs.get("richiedente_telefono"),
        indirizzo_preciso=kwargs.get("indirizzo_preciso"),
        interno=kwargs.get("interno"),
        chiamata_id=kwargs.get("chiamata_id"),
    )
    db.add(ticket)

    evento = TicketEvent(
        ticket=ticket,
        tipo="creato",
        descrizione=f"Ticket creato via assistente telefonico. Urgenza: {kwargs['urgenza']}. {kwargs['oggetto']}",
        autore="assistente_ai",
    )
    db.add(evento)
    await db.commit()
    await db.refresh(ticket)

    return {
        "successo": True,
        "numero_ticket": ticket.numero,
        "messaggio": f"Ticket {ticket.numero} creato con successo. Urgenza: {ticket.urgenza}.",
    }


async def aggiorna_ticket(
    db: AsyncSession,
    ticket_numero: str,
    stato: str | None = None,
    nota: str | None = None,
    fornitore_id: int | None = None,
) -> dict:
    """Aggiorna un ticket esistente."""
    stmt = select(Ticket).where(Ticket.numero == ticket_numero)
    result = await db.execute(stmt)
    ticket = result.scalar_one_or_none()

    if not ticket:
        return {"successo": False, "messaggio": f"Ticket {ticket_numero} non trovato."}

    changes = []
    if stato:
        ticket.stato = stato
        changes.append(f"Stato aggiornato a '{stato}'")
    if fornitore_id:
        ticket.fornitore_id = fornitore_id
        ticket.fornitore_contattato = True
        changes.append(f"Fornitore assegnato (ID: {fornitore_id})")
    if nota:
        evento = TicketEvent(
            ticket=ticket,
            tipo="nota" if not stato else "aggiornato",
            descrizione=nota,
            autore="assistente_ai",
        )
        db.add(evento)
        changes.append("Nota aggiunta")

    await db.commit()

    return {
        "successo": True,
        "messaggio": f"Ticket {ticket_numero} aggiornato: {', '.join(changes)}.",
    }


async def notifica_fornitore(
    db: AsyncSession,
    fornitore_id: int,
    ticket_numero: str,
    messaggio: str,
    urgente: bool = False,
) -> dict:
    """Simula invio notifica a fornitore. In produzione integra Twilio SMS/Call."""
    stmt = select(Fornitore).where(Fornitore.id == fornitore_id)
    result = await db.execute(stmt)
    fornitore = result.scalar_one_or_none()

    if not fornitore:
        return {"successo": False, "messaggio": f"Fornitore ID {fornitore_id} non trovato."}

    telefono = fornitore.telefono_reperibilita if urgente and fornitore.telefono_reperibilita else fornitore.telefono

    # Registra evento sul ticket
    stmt2 = select(Ticket).where(Ticket.numero == ticket_numero)
    result2 = await db.execute(stmt2)
    ticket = result2.scalar_one_or_none()
    if ticket:
        evento = TicketEvent(
            ticket=ticket,
            tipo="fornitore_contattato",
            descrizione=f"Notifica inviata a {fornitore.nome} ({telefono}): {messaggio}",
            autore="assistente_ai",
        )
        db.add(evento)
        ticket.fornitore_id = fornitore_id
        ticket.fornitore_contattato = True
        await db.commit()

    return {
        "successo": True,
        "fornitore": fornitore.nome,
        "telefono_usato": telefono,
        "messaggio": f"Notifica inviata a {fornitore.nome} al numero {telefono}.",
    }


async def escalation_umano(
    db: AsyncSession,
    motivo: str,
    priorita: str,
    ticket_numero: str | None = None,
) -> dict:
    """Registra escalation a operatore umano."""
    if ticket_numero:
        stmt = select(Ticket).where(Ticket.numero == ticket_numero)
        result = await db.execute(stmt)
        ticket = result.scalar_one_or_none()
        if ticket:
            ticket.escalation = True
            ticket.escalation_motivo = motivo
            ticket.stato = "escalation"
            evento = TicketEvent(
                ticket=ticket,
                tipo="escalation",
                descrizione=f"Escalation a operatore umano. Motivo: {motivo}. Priorità: {priorita}.",
                autore="assistente_ai",
            )
            db.add(evento)
            await db.commit()

    return {
        "successo": True,
        "messaggio": f"Escalation registrata con priorità {priorita}. Un operatore verrà notificato.",
    }


# Mappatura nome tool -> funzione
TOOL_HANDLERS = {
    "cerca_condominio": cerca_condominio,
    "cerca_contatto": cerca_contatto,
    "verifica_ticket_aperti": verifica_ticket_aperti,
    "cerca_fornitore": cerca_fornitore,
    "crea_ticket": crea_ticket,
    "aggiorna_ticket": aggiorna_ticket,
    "notifica_fornitore": notifica_fornitore,
    "escalation_umano": escalation_umano,
}
