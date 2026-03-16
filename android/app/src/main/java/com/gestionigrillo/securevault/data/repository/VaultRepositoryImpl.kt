package com.gestionigrillo.securevault.data.repository

import com.gestionigrillo.securevault.data.database.dao.VaultEntryDao
import com.gestionigrillo.securevault.data.database.entity.VaultEntryEntity
import com.gestionigrillo.securevault.domain.model.VaultEntry
import com.gestionigrillo.securevault.domain.repository.VaultRepository
import com.gestionigrillo.securevault.vaultcrypto.VaultCryptoEngine

class VaultRepositoryImpl(
    private val dao: VaultEntryDao,
    private val crypto: VaultCryptoEngine
) : VaultRepository {

    override suspend fun saveEntry(
        title: String,
        category: String,
        plainData: ByteArray
    ): Result<VaultEntry> = runCatching {
        val encrypted = crypto.encrypt(plainData)
        val now = System.currentTimeMillis()
        val entity = VaultEntryEntity(
            title = title,
            category = category,
            encryptedData = encrypted.ciphertext,
            iv = encrypted.iv,
            createdAt = now,
            updatedAt = now
        )
        val id = dao.insert(entity)
        entity.copy(id = id).toDomain()
    }

    override suspend fun getEntry(id: Long): Result<VaultEntry> = runCatching {
        dao.getById(id)?.toDomain() ?: throw Exception("Entry not found")
    }

    override suspend fun getAllEntries(): Result<List<VaultEntry>> = runCatching {
        dao.getAll().map { it.toDomain() }
    }

    override suspend fun decryptEntry(entry: VaultEntry): Result<ByteArray> = runCatching {
        crypto.decrypt(entry.encryptedData, entry.iv)
    }

    override suspend fun deleteEntry(id: Long): Result<Unit> = runCatching {
        dao.deleteById(id)
    }

    override suspend fun wipeAll(): Result<Unit> = runCatching {
        dao.deleteAll()
    }
}
