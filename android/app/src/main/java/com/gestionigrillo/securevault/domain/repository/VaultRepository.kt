package com.gestionigrillo.securevault.domain.repository

import com.gestionigrillo.securevault.domain.model.VaultEntry

interface VaultRepository {
    suspend fun saveEntry(title: String, category: String, plainData: ByteArray): Result<VaultEntry>
    suspend fun getEntry(id: Long): Result<VaultEntry>
    suspend fun getAllEntries(): Result<List<VaultEntry>>
    suspend fun decryptEntry(entry: VaultEntry): Result<ByteArray>
    suspend fun deleteEntry(id: Long): Result<Unit>
    suspend fun wipeAll(): Result<Unit>
}
