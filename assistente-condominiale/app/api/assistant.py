"""API per l'interazione diretta con l'assistente AI (chat testuale / debug)."""

from fastapi import APIRouter, Depends
from pydantic import BaseModel
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db
from app.services.assistant import process_message, get_session_history, close_session

router = APIRouter(prefix="/api/assistant", tags=["assistant"])


class ChatRequest(BaseModel):
    session_id: str
    message: str
    caller_phone: str | None = None


class ChatResponse(BaseModel):
    session_id: str
    response: str


@router.post("/chat", response_model=ChatResponse)
async def chat(req: ChatRequest, db: AsyncSession = Depends(get_db)):
    """Endpoint per interagire con l'assistente in modalità testuale (debug/demo)."""
    caller_info = {"phone": req.caller_phone} if req.caller_phone else None
    response_text = await process_message(db, req.session_id, req.message, caller_info)
    return ChatResponse(session_id=req.session_id, response=response_text)


@router.get("/history/{session_id}")
async def history(session_id: str):
    """Restituisce la cronologia di una sessione."""
    hist = get_session_history(session_id)
    return {"session_id": session_id, "messages": len(hist)}


@router.post("/close/{session_id}")
async def close(session_id: str):
    """Chiude una sessione di conversazione."""
    result = close_session(session_id)
    return {"closed": result is not None}
