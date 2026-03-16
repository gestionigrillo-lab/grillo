package com.gestionigrillo.securevault.domain.model

data class Device(
    val id: String,
    val token: String,
    val deviceModel: String,
    val androidVersion: String,
    val appVersion: String,
    val status: DeviceStatus,
    val enrolledAt: Long
)

enum class DeviceStatus {
    ACTIVE, LOCKED, REVOKED, WIPED;

    companion object {
        fun fromString(value: String): DeviceStatus =
            entries.firstOrNull { it.name.equals(value, ignoreCase = true) } ?: ACTIVE
    }
}
