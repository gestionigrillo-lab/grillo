package com.gestionigrillo.securevault.vaultcrypto

import java.security.SecureRandom
import javax.crypto.SecretKeyFactory
import javax.crypto.spec.PBEKeySpec

/**
 * PBKDF2 hashing for vault PIN verification.
 */
object PinHasher {

    private const val ITERATIONS = 120_000
    private const val KEY_LENGTH = 256
    private const val SALT_LENGTH = 32
    private const val ALGORITHM = "PBKDF2WithHmacSHA256"

    data class HashedPin(val hash: ByteArray, val salt: ByteArray) {
        override fun equals(other: Any?): Boolean {
            if (this === other) return true
            if (other !is HashedPin) return false
            return hash.contentEquals(other.hash) && salt.contentEquals(other.salt)
        }
        override fun hashCode(): Int = hash.contentHashCode()
    }

    fun hash(pin: String): HashedPin {
        val salt = ByteArray(SALT_LENGTH).also { SecureRandom().nextBytes(it) }
        val hash = deriveKey(pin, salt)
        return HashedPin(hash, salt)
    }

    fun verify(pin: String, hashedPin: HashedPin): Boolean {
        val derived = deriveKey(pin, hashedPin.salt)
        return derived.contentEquals(hashedPin.hash)
    }

    private fun deriveKey(pin: String, salt: ByteArray): ByteArray {
        val spec = PBEKeySpec(pin.toCharArray(), salt, ITERATIONS, KEY_LENGTH)
        val factory = SecretKeyFactory.getInstance(ALGORITHM)
        return factory.generateSecret(spec).encoded
    }
}
