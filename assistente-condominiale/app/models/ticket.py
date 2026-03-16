import datetime
from sqlalchemy import String, Integer, ForeignKey, DateTime, Text, func
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.models.base import Base, TimestampMixin


class Ticket(Base, TimestampMixin):
    __tablename__ = "tickets"

    id: Mapped[int] = mapped_column(primary_key=True)
    numero: Mapped[str] = mapped_column(String(30), unique=True, nullable=False)
    condominio_id: Mapped[int] = mapped_column(ForeignKey("condomini.id"), nullable=False)
    fornitore_id: Mapped[int] = mapped_column(ForeignKey("fornitori.id"), nullable=True)
    chiamata_id: Mapped[int] = mapped_column(ForeignKey("chiamate.id"), nullable=True)

    # Classificazione
    categoria: Mapped[str] = mapped_column(String(100), nullable=False)
    # Categorie: idraulica, elettrica, ascensore, serrature_accessi,
    #            edile, pulizia, verde, riscaldamento, infiltrazioni, altro
    urgenza: Mapped[str] = mapped_column(String(30), nullable=False)  # urgente, media, ordinaria
    stato: Mapped[str] = mapped_column(String(30), nullable=False, default="aperto")
    # Stati: aperto, in_lavorazione, in_attesa, risolto, chiuso, escalation

    # Dettagli
    oggetto: Mapped[str] = mapped_column(String(500), nullable=False)
    descrizione: Mapped[str] = mapped_column(Text, nullable=False)
    richiedente_nome: Mapped[str] = mapped_column(String(255), nullable=True)
    richiedente_ruolo: Mapped[str] = mapped_column(String(100), nullable=True)
    richiedente_telefono: Mapped[str] = mapped_column(String(30), nullable=True)
    indirizzo_preciso: Mapped[str] = mapped_column(String(500), nullable=True)
    interno: Mapped[str] = mapped_column(String(20), nullable=True)

    # Fornitore
    fornitore_contattato: Mapped[bool] = mapped_column(default=False, nullable=False)
    fornitore_esito: Mapped[str] = mapped_column(Text, nullable=True)

    # Escalation
    escalation: Mapped[bool] = mapped_column(default=False, nullable=False)
    escalation_motivo: Mapped[str] = mapped_column(Text, nullable=True)

    note: Mapped[str] = mapped_column(Text, nullable=True)

    condominio: Mapped["Condominio"] = relationship(back_populates="tickets")
    fornitore: Mapped["Fornitore"] = relationship()
    eventi: Mapped[list["TicketEvent"]] = relationship(back_populates="ticket", cascade="all, delete-orphan")

    def __repr__(self):
        return f"<Ticket {self.numero}: {self.oggetto} [{self.urgenza}/{self.stato}]>"


class TicketEvent(Base):
    __tablename__ = "ticket_events"

    id: Mapped[int] = mapped_column(primary_key=True)
    ticket_id: Mapped[int] = mapped_column(ForeignKey("tickets.id"), nullable=False)
    tipo: Mapped[str] = mapped_column(String(50), nullable=False)
    # Tipi evento: creato, aggiornato, fornitore_contattato, fornitore_risposta,
    #              escalation, risolto, chiuso, nota, chiamata_in, chiamata_out
    descrizione: Mapped[str] = mapped_column(Text, nullable=False)
    autore: Mapped[str] = mapped_column(String(100), nullable=False, default="assistente_ai")
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), nullable=False
    )

    ticket: Mapped["Ticket"] = relationship(back_populates="eventi")

    def __repr__(self):
        return f"<TicketEvent {self.id}: {self.tipo}>"
