import datetime
from sqlalchemy import String, Integer, DateTime, Text, func
from sqlalchemy.orm import Mapped, mapped_column
from app.models.base import Base, TimestampMixin


class Chiamata(Base, TimestampMixin):
    """Registro chiamate in entrata e in uscita."""
    __tablename__ = "chiamate"

    id: Mapped[int] = mapped_column(primary_key=True)
    direzione: Mapped[str] = mapped_column(String(10), nullable=False)  # in, out
    numero_chiamante: Mapped[str] = mapped_column(String(30), nullable=True)
    numero_chiamato: Mapped[str] = mapped_column(String(30), nullable=True)

    # Dati raccolti durante la chiamata
    nome_chiamante: Mapped[str] = mapped_column(String(255), nullable=True)
    ruolo_chiamante: Mapped[str] = mapped_column(String(100), nullable=True)
    condominio_riferimento: Mapped[str] = mapped_column(String(255), nullable=True)
    condominio_id: Mapped[int] = mapped_column(Integer, nullable=True)

    # Classificazione
    urgenza: Mapped[str] = mapped_column(String(30), nullable=True)
    categoria: Mapped[str] = mapped_column(String(100), nullable=True)
    problema_sintetico: Mapped[str] = mapped_column(String(500), nullable=True)

    # Esito
    ticket_creato: Mapped[str] = mapped_column(String(30), nullable=True)
    fornitore_contattato: Mapped[str] = mapped_column(String(255), nullable=True)
    escalation: Mapped[bool] = mapped_column(default=False, nullable=False)
    esito: Mapped[str] = mapped_column(Text, nullable=True)

    # Trascrizione
    trascrizione: Mapped[str] = mapped_column(Text, nullable=True)

    # Durata
    durata_secondi: Mapped[int] = mapped_column(Integer, nullable=True)
    inizio: Mapped[datetime.datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), nullable=False
    )
    fine: Mapped[datetime.datetime] = mapped_column(DateTime(timezone=True), nullable=True)

    # Twilio
    twilio_call_sid: Mapped[str] = mapped_column(String(50), nullable=True)

    def __repr__(self):
        return f"<Chiamata {self.id}: {self.direzione} {self.numero_chiamante}>"
