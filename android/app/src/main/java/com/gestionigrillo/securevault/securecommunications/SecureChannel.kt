package com.gestionigrillo.securevault.securecommunications

/**
 * Future module: Secure Communications.
 * Interfaces prepared for encrypted messaging, calls, and video calls.
 * Implementation will use end-to-end encryption with Signal Protocol or similar.
 */

interface SecureChannel {
    suspend fun initSession(peerDeviceId: String): Result<SessionInfo>
    suspend fun sendMessage(sessionId: String, plaintext: ByteArray): Result<String>
    suspend fun receiveMessage(sessionId: String, ciphertext: ByteArray): Result<ByteArray>
    suspend fun closeSession(sessionId: String): Result<Unit>
}

data class SessionInfo(
    val sessionId: String,
    val peerDeviceId: String,
    val established: Boolean
)
