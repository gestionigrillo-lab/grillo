"""API REST per il registro chiamate."""

from fastapi import APIRouter, Depends, Query
from pydantic import BaseModel
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db
from app.models.chiamata import Chiamata

router = APIRouter(prefix="/api/chiamate", tags=["chiamate"])


class ChiamataOut(BaseModel):
    id: int
    direzione: str
    numero_chiamante: str | None
    numero_chiamato: str | None
    nome_chiamante: str | None
    ruolo_chiamante: str | None
    condominio_riferimento: str | None
    condominio_id: int | None
    urgenza: str | None
    categoria: str | None
    problema_sintetico: str | None
    ticket_creato: str | None
    fornitore_contattato: str | None
    escalation: bool
    esito: str | None
    durata_secondi: int | None
    inizio: str | None
    fine: str | None

    class Config:
        from_attributes = True


@router.get("/", response_model=list[ChiamataOut])
async def list_chiamate(
    limit: int = Query(default=50, le=200),
    offset: int = 0,
    db: AsyncSession = Depends(get_db),
):
    stmt = select(Chiamata).order_by(Chiamata.inizio.desc()).offset(offset).limit(limit)
    result = await db.execute(stmt)
    return result.scalars().all()


@router.get("/attive")
async def get_active_sessions():
    from app.services.assistant import get_active_sessions
    sessions = get_active_sessions()
    return {"sessioni_attive": len(sessions), "ids": sessions}
