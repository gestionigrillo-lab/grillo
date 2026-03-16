"""API REST per la gestione dei ticket."""

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import BaseModel
from sqlalchemy import select, func
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.core.database import get_db
from app.models.ticket import Ticket, TicketEvent

router = APIRouter(prefix="/api/tickets", tags=["tickets"])


class TicketOut(BaseModel):
    id: int
    numero: str
    condominio_id: int
    fornitore_id: int | None
    categoria: str
    urgenza: str
    stato: str
    oggetto: str
    descrizione: str
    richiedente_nome: str | None
    richiedente_ruolo: str | None
    richiedente_telefono: str | None
    indirizzo_preciso: str | None
    interno: str | None
    fornitore_contattato: bool
    fornitore_esito: str | None
    escalation: bool
    escalation_motivo: str | None
    note: str | None
    created_at: str | None
    updated_at: str | None

    class Config:
        from_attributes = True


class TicketEventOut(BaseModel):
    id: int
    tipo: str
    descrizione: str
    autore: str
    created_at: str | None

    class Config:
        from_attributes = True


class TicketStats(BaseModel):
    totali: int
    aperti: int
    in_lavorazione: int
    in_attesa: int
    risolti: int
    chiusi: int
    escalation: int
    urgenti: int
    medi: int
    ordinari: int


@router.get("/", response_model=list[TicketOut])
async def list_tickets(
    stato: str | None = None,
    urgenza: str | None = None,
    condominio_id: int | None = None,
    q: str | None = None,
    limit: int = Query(default=50, le=200),
    offset: int = 0,
    db: AsyncSession = Depends(get_db),
):
    stmt = select(Ticket).order_by(Ticket.created_at.desc())
    if stato:
        stmt = stmt.where(Ticket.stato == stato)
    if urgenza:
        stmt = stmt.where(Ticket.urgenza == urgenza)
    if condominio_id:
        stmt = stmt.where(Ticket.condominio_id == condominio_id)
    if q:
        pattern = f"%{q}%"
        stmt = stmt.where(
            Ticket.oggetto.ilike(pattern) | Ticket.numero.ilike(pattern) | Ticket.descrizione.ilike(pattern)
        )
    stmt = stmt.offset(offset).limit(limit)
    result = await db.execute(stmt)
    return result.scalars().all()


@router.get("/stats", response_model=TicketStats)
async def get_stats(db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Ticket))
    tickets = result.scalars().all()
    return TicketStats(
        totali=len(tickets),
        aperti=sum(1 for t in tickets if t.stato == "aperto"),
        in_lavorazione=sum(1 for t in tickets if t.stato == "in_lavorazione"),
        in_attesa=sum(1 for t in tickets if t.stato == "in_attesa"),
        risolti=sum(1 for t in tickets if t.stato == "risolto"),
        chiusi=sum(1 for t in tickets if t.stato == "chiuso"),
        escalation=sum(1 for t in tickets if t.stato == "escalation"),
        urgenti=sum(1 for t in tickets if t.urgenza == "urgente"),
        medi=sum(1 for t in tickets if t.urgenza == "media"),
        ordinari=sum(1 for t in tickets if t.urgenza == "ordinaria"),
    )


@router.get("/{ticket_numero}", response_model=TicketOut)
async def get_ticket(ticket_numero: str, db: AsyncSession = Depends(get_db)):
    stmt = select(Ticket).where(Ticket.numero == ticket_numero)
    result = await db.execute(stmt)
    ticket = result.scalar_one_or_none()
    if not ticket:
        raise HTTPException(status_code=404, detail="Ticket non trovato")
    return ticket


@router.get("/{ticket_numero}/eventi", response_model=list[TicketEventOut])
async def get_ticket_events(ticket_numero: str, db: AsyncSession = Depends(get_db)):
    stmt = (
        select(TicketEvent)
        .join(Ticket)
        .where(Ticket.numero == ticket_numero)
        .order_by(TicketEvent.created_at.desc())
    )
    result = await db.execute(stmt)
    return result.scalars().all()


@router.patch("/{ticket_numero}/stato")
async def update_ticket_stato(ticket_numero: str, stato: str, nota: str | None = None, db: AsyncSession = Depends(get_db)):
    stmt = select(Ticket).where(Ticket.numero == ticket_numero)
    result = await db.execute(stmt)
    ticket = result.scalar_one_or_none()
    if not ticket:
        raise HTTPException(status_code=404, detail="Ticket non trovato")

    old_stato = ticket.stato
    ticket.stato = stato
    evento = TicketEvent(
        ticket=ticket,
        tipo="aggiornato",
        descrizione=f"Stato cambiato da '{old_stato}' a '{stato}'" + (f". Nota: {nota}" if nota else ""),
        autore="operatore",
    )
    db.add(evento)
    await db.commit()
    return {"successo": True, "messaggio": f"Ticket {ticket_numero} aggiornato a '{stato}'."}
