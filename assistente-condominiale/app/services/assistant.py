"""Servizio principale dell'assistente AI.

Gestisce la conversazione con Claude, inclusa l'esecuzione dei tool
e il mantenimento dello stato della conversazione per ogni chiamata.
"""

import json
import logging
from typing import Any

import anthropic
from sqlalchemy.ext.asyncio import AsyncSession

from app.config import settings
from app.core.system_prompt import SYSTEM_PROMPT, ASSISTANT_TOOLS
from app.services.gestionale import TOOL_HANDLERS

logger = logging.getLogger(__name__)

# Sessioni attive: call_sid -> conversation history
_active_sessions: dict[str, list[dict]] = {}


def _get_client() -> anthropic.Anthropic:
    return anthropic.Anthropic(api_key=settings.anthropic_api_key)


async def _execute_tool(db: AsyncSession, tool_name: str, tool_input: dict) -> str:
    """Esegue un tool del gestionale e restituisce il risultato come stringa JSON."""
    handler = TOOL_HANDLERS.get(tool_name)
    if not handler:
        return json.dumps({"errore": f"Strumento '{tool_name}' non disponibile."})

    try:
        # Passa db come primo argomento, il resto come kwargs
        result = await handler(db, **tool_input)
        return json.dumps(result, ensure_ascii=False, default=str)
    except Exception as e:
        logger.error(f"Errore esecuzione tool {tool_name}: {e}")
        return json.dumps({"errore": f"Errore nell'esecuzione di '{tool_name}': {str(e)}"})


async def process_message(
    db: AsyncSession,
    session_id: str,
    user_message: str,
    caller_info: dict | None = None,
) -> str:
    """Processa un messaggio dell'utente e restituisce la risposta dell'assistente.

    Args:
        db: Sessione database
        session_id: ID univoco della sessione/chiamata
        user_message: Messaggio dell'utente (trascrizione vocale)
        caller_info: Info opzionali sul chiamante (numero, ecc.)

    Returns:
        Risposta testuale dell'assistente (da convertire in voce)
    """
    client = _get_client()

    # Inizializza o recupera cronologia conversazione
    if session_id not in _active_sessions:
        _active_sessions[session_id] = []
        # Se abbiamo info sul chiamante, aggiungiamo contesto
        if caller_info:
            context = f"[SISTEMA: Chiamata in entrata dal numero {caller_info.get('phone', 'sconosciuto')}]"
            _active_sessions[session_id].append({"role": "user", "content": context})

    messages = _active_sessions[session_id]
    messages.append({"role": "user", "content": user_message})

    # Loop di conversazione con tool use
    max_iterations = 10
    for _ in range(max_iterations):
        response = client.messages.create(
            model=settings.claude_model,
            max_tokens=1024,
            system=SYSTEM_PROMPT,
            tools=ASSISTANT_TOOLS,
            messages=messages,
        )

        # Raccogli tutti i blocchi della risposta
        assistant_content = response.content
        messages.append({"role": "assistant", "content": assistant_content})

        # Se stop_reason è "end_turn", restituisci il testo
        if response.stop_reason == "end_turn":
            text_parts = [
                block.text for block in assistant_content
                if hasattr(block, "text")
            ]
            return " ".join(text_parts)

        # Se stop_reason è "tool_use", esegui i tool
        if response.stop_reason == "tool_use":
            tool_results = []
            for block in assistant_content:
                if block.type == "tool_use":
                    logger.info(f"Esecuzione tool: {block.name} con input: {block.input}")
                    result_str = await _execute_tool(db, block.name, block.input)
                    tool_results.append({
                        "type": "tool_result",
                        "tool_use_id": block.id,
                        "content": result_str,
                    })

            messages.append({"role": "user", "content": tool_results})
            continue

        # Caso inatteso
        text_parts = [
            block.text for block in assistant_content
            if hasattr(block, "text")
        ]
        return " ".join(text_parts) if text_parts else "Mi scusi, c'è stato un problema. Può ripetere?"

    return "Mi scusi, sto avendo difficoltà a elaborare la richiesta. La metto in contatto con un operatore."


def get_session_history(session_id: str) -> list[dict]:
    """Restituisce la cronologia di una sessione."""
    return _active_sessions.get(session_id, [])


def close_session(session_id: str) -> dict | None:
    """Chiude una sessione e restituisce lo storico."""
    return _active_sessions.pop(session_id, None)


def get_active_sessions() -> list[str]:
    """Restituisce gli ID delle sessioni attive."""
    return list(_active_sessions.keys())
