"""API REST per la gestione dei fornitori."""

from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db
from app.models.fornitore import Fornitore

router = APIRouter(prefix="/api/fornitori", tags=["fornitori"])


class FornitoreCreate(BaseModel):
    nome: str
    tipo: str
    telefono: str
    telefono_reperibilita: str | None = None
    email: str | None = None
    piva: str | None = None
    indirizzo: str | None = None
    zone_coperte: str | None = None
    attivo: bool = True
    reperibile_h24: bool = False
    note: str | None = None


class FornitoreOut(FornitoreCreate):
    id: int

    class Config:
        from_attributes = True


@router.get("/", response_model=list[FornitoreOut])
async def list_fornitori(tipo: str | None = None, q: str | None = None, db: AsyncSession = Depends(get_db)):
    stmt = select(Fornitore).order_by(Fornitore.nome)
    if tipo:
        stmt = stmt.where(Fornitore.tipo == tipo)
    if q:
        pattern = f"%{q}%"
        stmt = stmt.where(Fornitore.nome.ilike(pattern))
    result = await db.execute(stmt)
    return result.scalars().all()


@router.get("/{fornitore_id}", response_model=FornitoreOut)
async def get_fornitore(fornitore_id: int, db: AsyncSession = Depends(get_db)):
    stmt = select(Fornitore).where(Fornitore.id == fornitore_id)
    result = await db.execute(stmt)
    forn = result.scalar_one_or_none()
    if not forn:
        raise HTTPException(status_code=404, detail="Fornitore non trovato")
    return forn


@router.post("/", response_model=FornitoreOut, status_code=201)
async def create_fornitore(data: FornitoreCreate, db: AsyncSession = Depends(get_db)):
    forn = Fornitore(**data.model_dump())
    db.add(forn)
    await db.commit()
    await db.refresh(forn)
    return forn


@router.put("/{fornitore_id}", response_model=FornitoreOut)
async def update_fornitore(fornitore_id: int, data: FornitoreCreate, db: AsyncSession = Depends(get_db)):
    stmt = select(Fornitore).where(Fornitore.id == fornitore_id)
    result = await db.execute(stmt)
    forn = result.scalar_one_or_none()
    if not forn:
        raise HTTPException(status_code=404, detail="Fornitore non trovato")
    for k, v in data.model_dump().items():
        setattr(forn, k, v)
    await db.commit()
    await db.refresh(forn)
    return forn


@router.delete("/{fornitore_id}", status_code=204)
async def delete_fornitore(fornitore_id: int, db: AsyncSession = Depends(get_db)):
    stmt = select(Fornitore).where(Fornitore.id == fornitore_id)
    result = await db.execute(stmt)
    forn = result.scalar_one_or_none()
    if not forn:
        raise HTTPException(status_code=404, detail="Fornitore non trovato")
    await db.delete(forn)
    await db.commit()
