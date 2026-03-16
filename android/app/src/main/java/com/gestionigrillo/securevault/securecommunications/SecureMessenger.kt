package com.gestionigrillo.securevault.securecommunications

/**
 * Future module: Encrypted messaging.
 * Will support text, media, and file attachments with E2E encryption.
 */
interface SecureMessenger {
    suspend fun sendTextMessage(recipientId: String, text: String): Result<MessageReceipt>
    suspend fun sendMediaMessage(recipientId: String, media: ByteArray, mimeType: String): Result<MessageReceipt>
    suspend fun getMessages(conversationId: String, limit: Int, offset: Int): Result<List<SecureMessage>>
    suspend fun markAsRead(messageId: String): Result<Unit>
}

data class SecureMessage(
    val id: String,
    val senderId: String,
    val recipientId: String,
    val content: ByteArray,
    val mimeType: String,
    val timestamp: Long,
    val isRead: Boolean
) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is SecureMessage) return false
        return id == other.id
    }
    override fun hashCode(): Int = id.hashCode()
}

data class MessageReceipt(
    val messageId: String,
    val timestamp: Long,
    val delivered: Boolean
)
