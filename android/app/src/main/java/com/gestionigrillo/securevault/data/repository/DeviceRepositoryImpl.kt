package com.gestionigrillo.securevault.data.repository

import com.gestionigrillo.securevault.data.database.dao.DeviceDao
import com.gestionigrillo.securevault.data.database.entity.DeviceEntity
import com.gestionigrillo.securevault.data.network.ApiService
import com.gestionigrillo.securevault.data.network.dto.AckRequest
import com.gestionigrillo.securevault.data.network.dto.EnrollRequest
import com.gestionigrillo.securevault.data.network.dto.HeartbeatRequest
import com.gestionigrillo.securevault.domain.model.CommandType
import com.gestionigrillo.securevault.domain.model.Device
import com.gestionigrillo.securevault.domain.model.DeviceStatus
import com.gestionigrillo.securevault.domain.model.RemoteCommand
import com.gestionigrillo.securevault.domain.repository.DeviceRepository

class DeviceRepositoryImpl(
    private val api: ApiService,
    private val deviceDao: DeviceDao
) : DeviceRepository {

    override suspend fun enroll(
        deviceModel: String,
        androidVersion: String,
        appVersion: String
    ): Result<Device> = runCatching {
        val response = api.enroll(EnrollRequest(deviceModel, androidVersion, appVersion))
        if (!response.isSuccessful) {
            throw Exception("Enrollment failed: ${response.code()}")
        }
        val body = response.body() ?: throw Exception("Empty response")
        Device(
            id = body.deviceId,
            token = body.token,
            deviceModel = deviceModel,
            androidVersion = androidVersion,
            appVersion = appVersion,
            status = DeviceStatus.ACTIVE,
            enrolledAt = System.currentTimeMillis()
        )
    }

    override suspend fun sendHeartbeat(
        deviceId: String,
        token: String,
        appVersion: String,
        batteryPct: Int,
        isRooted: Boolean
    ): Result<String> = runCatching {
        val response = api.heartbeat(
            "Bearer $token",
            HeartbeatRequest(deviceId, appVersion, batteryPct, isRooted)
        )
        if (!response.isSuccessful) {
            throw Exception("Heartbeat failed: ${response.code()}")
        }
        response.body()?.status ?: "unknown"
    }

    override suspend fun fetchPendingCommands(
        deviceId: String,
        token: String
    ): Result<List<RemoteCommand>> = runCatching {
        val response = api.getCommands("Bearer $token", deviceId)
        if (!response.isSuccessful) {
            throw Exception("Fetch commands failed: ${response.code()}")
        }
        val body = response.body() ?: throw Exception("Empty response")
        body.commands.map { dto ->
            RemoteCommand(
                id = dto.id,
                deviceId = dto.deviceId,
                type = CommandType.fromString(dto.type),
                payload = dto.payload ?: ""
            )
        }
    }

    override suspend fun acknowledgeCommand(
        commandId: String,
        token: String,
        success: Boolean,
        message: String
    ): Result<Unit> = runCatching {
        val device = getEnrolledDevice() ?: throw Exception("No enrolled device")
        val response = api.ackCommand(
            "Bearer $token",
            commandId,
            AckRequest(device.id, success, message)
        )
        if (!response.isSuccessful) {
            throw Exception("Ack failed: ${response.code()}")
        }
    }

    override suspend fun getEnrolledDevice(): Device? {
        return deviceDao.getDevice()?.toDomain()
    }

    override suspend fun saveDevice(device: Device) {
        deviceDao.insert(DeviceEntity.fromDomain(device))
    }

    override suspend fun clearDevice() {
        deviceDao.deleteAll()
    }
}
