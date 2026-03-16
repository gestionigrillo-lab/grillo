from app.models.base import Base
from app.models.condominio import Condominio, Contatto
from app.models.ticket import Ticket, TicketEvent
from app.models.fornitore import Fornitore
from app.models.chiamata import Chiamata

__all__ = [
    "Base",
    "Condominio",
    "Contatto",
    "Ticket",
    "TicketEvent",
    "Fornitore",
    "Chiamata",
]
