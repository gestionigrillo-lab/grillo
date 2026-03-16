package com.gestionigrillo.securevault.data.network.dto

import com.google.gson.annotations.SerializedName

data class EnrollRequest(
    @SerializedName("device_model") val deviceModel: String,
    @SerializedName("android_version") val androidVersion: String,
    @SerializedName("app_version") val appVersion: String,
    @SerializedName("push_token") val pushToken: String = ""
)

data class EnrollResponse(
    @SerializedName("device_id") val deviceId: String,
    val token: String
)

data class HeartbeatRequest(
    @SerializedName("device_id") val deviceId: String,
    @SerializedName("app_version") val appVersion: String,
    @SerializedName("battery_pct") val batteryPct: Int,
    @SerializedName("is_rooted") val isRooted: Boolean
)

data class HeartbeatResponse(
    val status: String,
    @SerializedName("pending_commands") val pendingCommands: List<CommandDto>?
)

data class CommandDto(
    val id: String,
    @SerializedName("device_id") val deviceId: String,
    val type: String,
    val payload: String?
)

data class CommandsResponse(
    val commands: List<CommandDto>
)

data class AckRequest(
    @SerializedName("device_id") val deviceId: String,
    val success: Boolean,
    val message: String = ""
)
