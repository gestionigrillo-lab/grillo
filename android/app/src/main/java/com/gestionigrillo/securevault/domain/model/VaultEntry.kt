package com.gestionigrillo.securevault.domain.model

data class VaultEntry(
    val id: Long = 0,
    val title: String,
    val category: String,
    val encryptedData: ByteArray,
    val iv: ByteArray,
    val createdAt: Long,
    val updatedAt: Long
) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is VaultEntry) return false
        return id == other.id
    }

    override fun hashCode(): Int = id.hashCode()
}
