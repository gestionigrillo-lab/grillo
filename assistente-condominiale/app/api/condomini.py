"""API REST per la gestione dei condomini."""

from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.core.database import get_db
from app.models.condominio import Condominio, Contatto

router = APIRouter(prefix="/api/condomini", tags=["condomini"])


class CondominioCreate(BaseModel):
    nome: str
    indirizzo: str
    citta: str = ""
    cap: str = ""
    codice_fiscale: str | None = None
    num_unita: int = 0
    telefono: str | None = None
    email: str | None = None
    note: str | None = None


class CondominioOut(CondominioCreate):
    id: int

    class Config:
        from_attributes = True


class ContattoCreate(BaseModel):
    condominio_id: int
    nome: str
    ruolo: str
    telefono: str | None = None
    email: str | None = None
    interno: str | None = None
    note: str | None = None


class ContattoOut(ContattoCreate):
    id: int

    class Config:
        from_attributes = True


@router.get("/", response_model=list[CondominioOut])
async def list_condomini(q: str | None = None, db: AsyncSession = Depends(get_db)):
    stmt = select(Condominio).order_by(Condominio.nome)
    if q:
        pattern = f"%{q}%"
        stmt = stmt.where(Condominio.nome.ilike(pattern) | Condominio.indirizzo.ilike(pattern))
    result = await db.execute(stmt)
    return result.scalars().all()


@router.get("/{condominio_id}", response_model=CondominioOut)
async def get_condominio(condominio_id: int, db: AsyncSession = Depends(get_db)):
    stmt = select(Condominio).where(Condominio.id == condominio_id)
    result = await db.execute(stmt)
    cond = result.scalar_one_or_none()
    if not cond:
        raise HTTPException(status_code=404, detail="Condominio non trovato")
    return cond


@router.post("/", response_model=CondominioOut, status_code=201)
async def create_condominio(data: CondominioCreate, db: AsyncSession = Depends(get_db)):
    cond = Condominio(**data.model_dump())
    db.add(cond)
    await db.commit()
    await db.refresh(cond)
    return cond


@router.put("/{condominio_id}", response_model=CondominioOut)
async def update_condominio(condominio_id: int, data: CondominioCreate, db: AsyncSession = Depends(get_db)):
    stmt = select(Condominio).where(Condominio.id == condominio_id)
    result = await db.execute(stmt)
    cond = result.scalar_one_or_none()
    if not cond:
        raise HTTPException(status_code=404, detail="Condominio non trovato")
    for k, v in data.model_dump().items():
        setattr(cond, k, v)
    await db.commit()
    await db.refresh(cond)
    return cond


@router.delete("/{condominio_id}", status_code=204)
async def delete_condominio(condominio_id: int, db: AsyncSession = Depends(get_db)):
    stmt = select(Condominio).where(Condominio.id == condominio_id)
    result = await db.execute(stmt)
    cond = result.scalar_one_or_none()
    if not cond:
        raise HTTPException(status_code=404, detail="Condominio non trovato")
    await db.delete(cond)
    await db.commit()


# --- Contatti ---

@router.get("/{condominio_id}/contatti", response_model=list[ContattoOut])
async def list_contatti(condominio_id: int, db: AsyncSession = Depends(get_db)):
    stmt = select(Contatto).where(Contatto.condominio_id == condominio_id).order_by(Contatto.nome)
    result = await db.execute(stmt)
    return result.scalars().all()


@router.post("/{condominio_id}/contatti", response_model=ContattoOut, status_code=201)
async def create_contatto(condominio_id: int, data: ContattoCreate, db: AsyncSession = Depends(get_db)):
    data_dict = data.model_dump()
    data_dict["condominio_id"] = condominio_id
    contatto = Contatto(**data_dict)
    db.add(contatto)
    await db.commit()
    await db.refresh(contatto)
    return contatto
