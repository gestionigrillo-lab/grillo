from sqlalchemy import String, Text, Boolean
from sqlalchemy.orm import Mapped, mapped_column
from app.models.base import Base, TimestampMixin


class Fornitore(Base, TimestampMixin):
    __tablename__ = "fornitori"

    id: Mapped[int] = mapped_column(primary_key=True)
    nome: Mapped[str] = mapped_column(String(255), nullable=False)
    tipo: Mapped[str] = mapped_column(String(100), nullable=False)
    # Tipi: idraulico, elettricista, fabbro, ascensorista, muratore,
    #        vetraio, serramentista, impresa_pulizie, giardiniere,
    #        disinfestazione, antennista, caldaista, generico
    telefono: Mapped[str] = mapped_column(String(30), nullable=False)
    telefono_reperibilita: Mapped[str] = mapped_column(String(30), nullable=True)
    email: Mapped[str] = mapped_column(String(255), nullable=True)
    piva: Mapped[str] = mapped_column(String(20), nullable=True)
    indirizzo: Mapped[str] = mapped_column(String(500), nullable=True)
    zone_coperte: Mapped[str] = mapped_column(Text, nullable=True)  # JSON array di zone/quartieri
    attivo: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    reperibile_h24: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    note: Mapped[str] = mapped_column(Text, nullable=True)

    def __repr__(self):
        return f"<Fornitore {self.id}: {self.nome} ({self.tipo})>"
