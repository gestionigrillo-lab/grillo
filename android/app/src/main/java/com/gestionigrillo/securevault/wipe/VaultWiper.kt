package com.gestionigrillo.securevault.wipe

import android.content.Context
import com.gestionigrillo.securevault.domain.repository.VaultRepository
import com.gestionigrillo.securevault.vaultcrypto.KeyStoreManager

/**
 * Handles vault data destruction for lock, wipe, and revoke commands.
 */
object VaultWiper {

    /**
     * Lock the vault - clears session state forcing re-authentication.
     */
    fun lockVault(context: Context) {
        // Clear auth session prefs
        context.getSharedPreferences("vault_session", Context.MODE_PRIVATE)
            .edit()
            .clear()
            .apply()
    }

    /**
     * Full wipe - destroys all vault data and cryptographic keys.
     */
    suspend fun wipeVault(
        context: Context,
        vaultRepository: VaultRepository,
        keyStoreManager: KeyStoreManager
    ) {
        // 1. Wipe all encrypted entries from database
        vaultRepository.wipeAll()

        // 2. Destroy all cryptographic keys
        keyStoreManager.deleteAllKeys()

        // 3. Clear all shared preferences
        clearAllPreferences(context)

        // 4. Clear app cache
        context.cacheDir.deleteRecursively()
    }

    private fun clearAllPreferences(context: Context) {
        val prefsToWipe = listOf(
            "vault_auth_prefs",
            "vault_session",
            "device_prefs"
        )
        for (name in prefsToWipe) {
            context.getSharedPreferences(name, Context.MODE_PRIVATE)
                .edit()
                .clear()
                .apply()
        }
    }
}
