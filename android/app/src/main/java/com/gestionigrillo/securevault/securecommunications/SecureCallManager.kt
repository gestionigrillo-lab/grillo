package com.gestionigrillo.securevault.securecommunications

/**
 * Future module: Encrypted voice and video calls.
 * Will use WebRTC with SRTP for media encryption and custom signaling.
 */
interface SecureCallManager {
    suspend fun initiateCall(peerDeviceId: String, callType: CallType): Result<CallSession>
    suspend fun acceptCall(callId: String): Result<CallSession>
    suspend fun rejectCall(callId: String): Result<Unit>
    suspend fun endCall(callId: String): Result<Unit>
    fun setCallStateListener(listener: CallStateListener)
}

enum class CallType {
    VOICE, VIDEO
}

data class CallSession(
    val callId: String,
    val peerDeviceId: String,
    val callType: CallType,
    val state: CallState
)

enum class CallState {
    INITIATING, RINGING, CONNECTED, ENDED, FAILED
}

interface CallStateListener {
    fun onCallStateChanged(callId: String, state: CallState)
    fun onCallError(callId: String, error: String)
}
