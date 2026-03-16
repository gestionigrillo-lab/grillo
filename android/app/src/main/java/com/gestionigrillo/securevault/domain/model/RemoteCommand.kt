package com.gestionigrillo.securevault.domain.model

data class RemoteCommand(
    val id: String,
    val deviceId: String,
    val type: CommandType,
    val payload: String
)

enum class CommandType {
    LOCK, WIPE, REVOKE;

    companion object {
        fun fromString(value: String): CommandType =
            entries.firstOrNull { it.name.equals(value, ignoreCase = true) } ?: LOCK
    }
}
