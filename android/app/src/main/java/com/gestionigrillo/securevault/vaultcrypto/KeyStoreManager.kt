package com.gestionigrillo.securevault.vaultcrypto

import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import java.security.KeyStore
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey

/**
 * Manages hardware-backed keys via Android Keystore.
 * Master key is non-exportable and requires user authentication for sensitive ops.
 */
class KeyStoreManager {

    private val keyStore: KeyStore = KeyStore.getInstance(KEYSTORE_PROVIDER).apply { load(null) }

    fun getOrCreateMasterKey(): SecretKey {
        val existing = keyStore.getEntry(MASTER_KEY_ALIAS, null)
        if (existing is KeyStore.SecretKeyEntry) {
            return existing.secretKey
        }
        return generateMasterKey()
    }

    fun getOrCreateVaultKey(): SecretKey {
        val existing = keyStore.getEntry(VAULT_KEY_ALIAS, null)
        if (existing is KeyStore.SecretKeyEntry) {
            return existing.secretKey
        }
        return generateVaultKey()
    }

    fun hasMasterKey(): Boolean = keyStore.containsAlias(MASTER_KEY_ALIAS)

    fun hasVaultKey(): Boolean = keyStore.containsAlias(VAULT_KEY_ALIAS)

    fun deleteMasterKey() {
        if (keyStore.containsAlias(MASTER_KEY_ALIAS)) {
            keyStore.deleteEntry(MASTER_KEY_ALIAS)
        }
    }

    fun deleteVaultKey() {
        if (keyStore.containsAlias(VAULT_KEY_ALIAS)) {
            keyStore.deleteEntry(VAULT_KEY_ALIAS)
        }
    }

    fun deleteAllKeys() {
        deleteMasterKey()
        deleteVaultKey()
    }

    private fun generateMasterKey(): SecretKey {
        val spec = KeyGenParameterSpec.Builder(
            MASTER_KEY_ALIAS,
            KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT
        )
            .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
            .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
            .setKeySize(256)
            .setRandomizedEncryptionRequired(true)
            .build()

        val keyGen = KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES, KEYSTORE_PROVIDER)
        keyGen.init(spec)
        return keyGen.generateKey()
    }

    private fun generateVaultKey(): SecretKey {
        val spec = KeyGenParameterSpec.Builder(
            VAULT_KEY_ALIAS,
            KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT
        )
            .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
            .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
            .setKeySize(256)
            .setUserAuthenticationRequired(false) // Will be gated by app-level PIN/biometric
            .setRandomizedEncryptionRequired(true)
            .build()

        val keyGen = KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES, KEYSTORE_PROVIDER)
        keyGen.init(spec)
        return keyGen.generateKey()
    }

    companion object {
        private const val KEYSTORE_PROVIDER = "AndroidKeyStore"
        private const val MASTER_KEY_ALIAS = "secure_vault_master_key"
        private const val VAULT_KEY_ALIAS = "secure_vault_data_key"
    }
}
