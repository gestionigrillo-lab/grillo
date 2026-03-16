package com.gestionigrillo.securevault.domain.repository

import com.gestionigrillo.securevault.domain.model.Device
import com.gestionigrillo.securevault.domain.model.RemoteCommand

interface DeviceRepository {
    suspend fun enroll(deviceModel: String, androidVersion: String, appVersion: String): Result<Device>
    suspend fun sendHeartbeat(deviceId: String, token: String, appVersion: String, batteryPct: Int, isRooted: Boolean): Result<String>
    suspend fun fetchPendingCommands(deviceId: String, token: String): Result<List<RemoteCommand>>
    suspend fun acknowledgeCommand(commandId: String, token: String, success: Boolean, message: String): Result<Unit>
    suspend fun getEnrolledDevice(): Device?
    suspend fun saveDevice(device: Device)
    suspend fun clearDevice()
}
