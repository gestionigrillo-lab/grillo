package com.gestionigrillo.securevault.data.database.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey
import com.gestionigrillo.securevault.domain.model.VaultEntry

@Entity(tableName = "vault_entries")
data class VaultEntryEntity(
    @PrimaryKey(autoGenerate = true) val id: Long = 0,
    val title: String,
    val category: String,
    @ColumnInfo(name = "encrypted_data") val encryptedData: ByteArray,
    val iv: ByteArray,
    @ColumnInfo(name = "created_at") val createdAt: Long,
    @ColumnInfo(name = "updated_at") val updatedAt: Long
) {
    fun toDomain(): VaultEntry = VaultEntry(
        id = id,
        title = title,
        category = category,
        encryptedData = encryptedData,
        iv = iv,
        createdAt = createdAt,
        updatedAt = updatedAt
    )

    companion object {
        fun fromDomain(entry: VaultEntry): VaultEntryEntity = VaultEntryEntity(
            id = entry.id,
            title = entry.title,
            category = entry.category,
            encryptedData = entry.encryptedData,
            iv = entry.iv,
            createdAt = entry.createdAt,
            updatedAt = entry.updatedAt
        )
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is VaultEntryEntity) return false
        return id == other.id
    }

    override fun hashCode(): Int = id.hashCode()
}
