package com.gestionigrillo.securevault.vaultcrypto

import javax.crypto.Cipher
import javax.crypto.spec.GCMParameterSpec

/**
 * AES-256-GCM encryption/decryption using Android Keystore backed keys.
 */
class VaultCryptoEngine(private val keyStoreManager: KeyStoreManager) {

    data class EncryptedPayload(
        val ciphertext: ByteArray,
        val iv: ByteArray
    ) {
        override fun equals(other: Any?): Boolean {
            if (this === other) return true
            if (other !is EncryptedPayload) return false
            return ciphertext.contentEquals(other.ciphertext) && iv.contentEquals(other.iv)
        }
        override fun hashCode(): Int = ciphertext.contentHashCode() * 31 + iv.contentHashCode()
    }

    fun encrypt(plaintext: ByteArray): EncryptedPayload {
        val key = keyStoreManager.getOrCreateVaultKey()
        val cipher = Cipher.getInstance(TRANSFORMATION)
        cipher.init(Cipher.ENCRYPT_MODE, key)

        val ciphertext = cipher.doFinal(plaintext)
        return EncryptedPayload(ciphertext = ciphertext, iv = cipher.iv)
    }

    fun decrypt(ciphertext: ByteArray, iv: ByteArray): ByteArray {
        val key = keyStoreManager.getOrCreateVaultKey()
        val cipher = Cipher.getInstance(TRANSFORMATION)
        val spec = GCMParameterSpec(GCM_TAG_LENGTH, iv)
        cipher.init(Cipher.DECRYPT_MODE, key, spec)
        return cipher.doFinal(ciphertext)
    }

    fun encryptFile(plainBytes: ByteArray): EncryptedPayload = encrypt(plainBytes)

    fun decryptFile(ciphertext: ByteArray, iv: ByteArray): ByteArray = decrypt(ciphertext, iv)

    companion object {
        private const val TRANSFORMATION = "AES/GCM/NoPadding"
        private const val GCM_TAG_LENGTH = 128
    }
}
