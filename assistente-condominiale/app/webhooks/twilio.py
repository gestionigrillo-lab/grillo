"""Webhook Twilio per gestione chiamate in entrata.

Integra Twilio Voice con l'assistente AI per gestire il flusso:
1. Chiamata in entrata -> risposta con messaggio di benvenuto
2. Raccolta input vocale (speech-to-text)
3. Invio a Claude per elaborazione
4. Risposta vocale (text-to-speech)
5. Loop fino a fine chiamata
"""

import logging
from fastapi import APIRouter, Depends, Request, Form
from fastapi.responses import Response
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db
from app.services.assistant import process_message, close_session
from app.models.chiamata import Chiamata

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/webhooks/twilio", tags=["twilio"])

# TwiML templates
TWIML_GATHER = """<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say language="it-IT" voice="Google.it-IT-Wavenet-A">{message}</Say>
    <Gather input="speech" language="it-IT" speechTimeout="3"
            action="/webhooks/twilio/process" method="POST"
            speechModel="phone_call">
        <Say language="it-IT" voice="Google.it-IT-Wavenet-A">Sono in ascolto.</Say>
    </Gather>
    <Say language="it-IT" voice="Google.it-IT-Wavenet-A">Non ho sentito nulla. Se ha bisogno, richiami pure. Arrivederci.</Say>
</Response>"""

TWIML_REDIRECT = """<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say language="it-IT" voice="Google.it-IT-Wavenet-A">{message}</Say>
    <Gather input="speech" language="it-IT" speechTimeout="3"
            action="/webhooks/twilio/process" method="POST"
            speechModel="phone_call">
    </Gather>
    <Say language="it-IT" voice="Google.it-IT-Wavenet-A">Se non ha altro, la saluto. Buona giornata!</Say>
</Response>"""

TWIML_GOODBYE = """<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say language="it-IT" voice="Google.it-IT-Wavenet-A">{message}</Say>
    <Hangup/>
</Response>"""


@router.post("/incoming")
async def incoming_call(request: Request, db: AsyncSession = Depends(get_db)):
    """Gestisce la chiamata in entrata da Twilio."""
    form = await request.form()
    call_sid = form.get("CallSid", "unknown")
    caller = form.get("From", "sconosciuto")

    logger.info(f"Chiamata in entrata: {call_sid} da {caller}")

    # Registra la chiamata
    chiamata = Chiamata(
        direzione="in",
        numero_chiamante=caller,
        twilio_call_sid=call_sid,
    )
    db.add(chiamata)
    await db.commit()

    # Messaggio di benvenuto
    welcome = (
        "Buongiorno, studio Gestioni Grillo, sono l'assistente. "
        "Come posso aiutarla?"
    )

    twiml = TWIML_GATHER.format(message=welcome)
    return Response(content=twiml, media_type="application/xml")


@router.post("/process")
async def process_speech(request: Request, db: AsyncSession = Depends(get_db)):
    """Processa l'input vocale trascritto da Twilio e genera risposta."""
    form = await request.form()
    call_sid = form.get("CallSid", "unknown")
    speech_result = form.get("SpeechResult", "")
    caller = form.get("From", "sconosciuto")

    if not speech_result:
        twiml = TWIML_GATHER.format(message="Mi scusi, non ho capito. Può ripetere?")
        return Response(content=twiml, media_type="application/xml")

    logger.info(f"[{call_sid}] Trascrizione: {speech_result}")

    # Processa con l'assistente AI
    caller_info = {"phone": caller, "call_sid": call_sid}
    response_text = await process_message(db, call_sid, speech_result, caller_info)

    # Verifica se la conversazione è finita (saluti finali)
    goodbye_keywords = ["arrivederci", "buona giornata", "buona serata", "la saluto"]
    is_goodbye = any(kw in response_text.lower() for kw in goodbye_keywords)

    if is_goodbye:
        twiml = TWIML_GOODBYE.format(message=response_text)
        close_session(call_sid)
    else:
        twiml = TWIML_REDIRECT.format(message=response_text)

    return Response(content=twiml, media_type="application/xml")


@router.post("/status")
async def call_status(request: Request, db: AsyncSession = Depends(get_db)):
    """Callback per aggiornamenti di stato della chiamata Twilio."""
    form = await request.form()
    call_sid = form.get("CallSid", "")
    status = form.get("CallStatus", "")
    duration = form.get("CallDuration", "0")

    logger.info(f"[{call_sid}] Status: {status}, Durata: {duration}s")

    if status in ("completed", "failed", "no-answer", "busy"):
        close_session(call_sid)

    return Response(content="<Response/>", media_type="application/xml")
