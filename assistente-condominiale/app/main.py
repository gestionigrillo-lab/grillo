"""FastAPI application principale - Assistente Condominiale Ufficio."""

import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.core.database import engine
from app.models.base import Base
from app.api.condomini import router as condomini_router
from app.api.tickets import router as tickets_router
from app.api.fornitori import router as fornitori_router
from app.api.chiamate import router as chiamate_router
from app.api.assistant import router as assistant_router
from app.webhooks.twilio import router as twilio_router

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Crea le tabelle al primo avvio (dev mode)."""
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    logger.info("Database tables created/verified.")
    yield
    await engine.dispose()


app = FastAPI(
    title="Assistente Condominiale Ufficio",
    description=(
        "API backend per l'assistente telefonico AI dello studio "
        "di amministrazione condominiale Gestioni Grillo."
    ),
    version="1.0.0",
    lifespan=lifespan,
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=[settings.frontend_url, "http://localhost:5173", "http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Routers
app.include_router(condomini_router)
app.include_router(tickets_router)
app.include_router(fornitori_router)
app.include_router(chiamate_router)
app.include_router(assistant_router)
app.include_router(twilio_router)


@app.get("/")
async def root():
    return {
        "nome": "Assistente Condominiale Ufficio",
        "versione": "1.0.0",
        "studio": "Gestioni Grillo",
        "stato": "operativo",
    }


@app.get("/health")
async def health():
    return {"status": "ok"}
