import datetime
from sqlalchemy import String, Integer, ForeignKey, DateTime, Text, func
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.models.base import Base, TimestampMixin


class Condominio(Base, TimestampMixin):
    __tablename__ = "condomini"

    id: Mapped[int] = mapped_column(primary_key=True)
    nome: Mapped[str] = mapped_column(String(255), nullable=False)
    indirizzo: Mapped[str] = mapped_column(String(500), nullable=False)
    citta: Mapped[str] = mapped_column(String(100), nullable=False, default="")
    cap: Mapped[str] = mapped_column(String(10), nullable=False, default="")
    codice_fiscale: Mapped[str] = mapped_column(String(20), nullable=True)
    num_unita: Mapped[int] = mapped_column(Integer, nullable=False, default=0)
    telefono: Mapped[str] = mapped_column(String(30), nullable=True)
    email: Mapped[str] = mapped_column(String(255), nullable=True)
    note: Mapped[str] = mapped_column(Text, nullable=True)

    contatti: Mapped[list["Contatto"]] = relationship(back_populates="condominio", cascade="all, delete-orphan")
    tickets: Mapped[list["Ticket"]] = relationship(back_populates="condominio")

    def __repr__(self):
        return f"<Condominio {self.id}: {self.nome}>"


class Contatto(Base, TimestampMixin):
    """Contatti associati a un condominio: portinere, condomino, amministratore delegato, ecc."""
    __tablename__ = "contatti"

    id: Mapped[int] = mapped_column(primary_key=True)
    condominio_id: Mapped[int] = mapped_column(ForeignKey("condomini.id"), nullable=False)
    nome: Mapped[str] = mapped_column(String(255), nullable=False)
    ruolo: Mapped[str] = mapped_column(String(100), nullable=False)  # portinere, condomino, tecnico, ecc.
    telefono: Mapped[str] = mapped_column(String(30), nullable=True)
    email: Mapped[str] = mapped_column(String(255), nullable=True)
    interno: Mapped[str] = mapped_column(String(20), nullable=True)
    note: Mapped[str] = mapped_column(Text, nullable=True)

    condominio: Mapped["Condominio"] = relationship(back_populates="contatti")

    def __repr__(self):
        return f"<Contatto {self.id}: {self.nome} ({self.ruolo})>"
