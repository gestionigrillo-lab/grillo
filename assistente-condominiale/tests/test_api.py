"""Test di base per le API dell'assistente condominiale."""

import pytest
from unittest.mock import AsyncMock, patch, MagicMock


def test_system_prompt_loaded():
    """Verifica che il system prompt sia caricato correttamente."""
    from app.core.system_prompt import SYSTEM_PROMPT, ASSISTANT_TOOLS

    assert "Assistente Condominiale Ufficio" in SYSTEM_PROMPT
    assert "Gestioni Grillo" in SYSTEM_PROMPT
    assert len(ASSISTANT_TOOLS) == 8

    tool_names = {t["name"] for t in ASSISTANT_TOOLS}
    assert "cerca_condominio" in tool_names
    assert "crea_ticket" in tool_names
    assert "cerca_fornitore" in tool_names
    assert "notifica_fornitore" in tool_names
    assert "escalation_umano" in tool_names


def test_urgency_classification_in_prompt():
    """Verifica che il prompt contenga la classificazione delle urgenze."""
    from app.core.system_prompt import SYSTEM_PROMPT

    assert "URGENTE IMMEDIATO" in SYSTEM_PROMPT
    assert "URGENZA MEDIA" in SYSTEM_PROMPT
    assert "ORDINARIO" in SYSTEM_PROMPT
    assert "fuga d'acqua" in SYSTEM_PROMPT.lower() or "fuga d\u2019acqua" in SYSTEM_PROMPT.lower()
    assert "ascensore" in SYSTEM_PROMPT.lower()
    assert "odore di gas" in SYSTEM_PROMPT.lower()


def test_tool_schemas_valid():
    """Verifica che gli schemi dei tool siano validi."""
    from app.core.system_prompt import ASSISTANT_TOOLS

    for tool in ASSISTANT_TOOLS:
        assert "name" in tool
        assert "description" in tool
        assert "input_schema" in tool
        schema = tool["input_schema"]
        assert schema["type"] == "object"
        assert "properties" in schema
        assert "required" in schema


def test_config_defaults():
    """Verifica i valori di default della configurazione."""
    from app.config import Settings

    s = Settings(anthropic_api_key="test-key")
    assert s.app_port == 8000
    assert s.claude_model == "claude-sonnet-4-20250514"
    assert s.frontend_url == "http://localhost:5173"


def test_ticket_numero_generation():
    """Verifica la generazione del numero ticket."""
    from app.services.gestionale import _gen_ticket_numero

    numero = _gen_ticket_numero()
    assert numero.startswith("TK-")
    parts = numero.split("-")
    assert len(parts) == 3
    assert len(parts[1]) == 8  # YYYYMMDD
    assert len(parts[2]) == 4  # random seq


def test_tool_handlers_mapping():
    """Verifica che tutti i tool abbiano un handler registrato."""
    from app.core.system_prompt import ASSISTANT_TOOLS
    from app.services.gestionale import TOOL_HANDLERS

    for tool in ASSISTANT_TOOLS:
        assert tool["name"] in TOOL_HANDLERS, f"Handler mancante per {tool['name']}"
