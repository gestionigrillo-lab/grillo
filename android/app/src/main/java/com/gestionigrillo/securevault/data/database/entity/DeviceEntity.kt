package com.gestionigrillo.securevault.data.database.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey
import com.gestionigrillo.securevault.domain.model.Device
import com.gestionigrillo.securevault.domain.model.DeviceStatus

@Entity(tableName = "enrolled_device")
data class DeviceEntity(
    @PrimaryKey val id: String,
    val token: String,
    @ColumnInfo(name = "device_model") val deviceModel: String,
    @ColumnInfo(name = "android_version") val androidVersion: String,
    @ColumnInfo(name = "app_version") val appVersion: String,
    val status: String,
    @ColumnInfo(name = "enrolled_at") val enrolledAt: Long
) {
    fun toDomain(): Device = Device(
        id = id,
        token = token,
        deviceModel = deviceModel,
        androidVersion = androidVersion,
        appVersion = appVersion,
        status = DeviceStatus.fromString(status),
        enrolledAt = enrolledAt
    )

    companion object {
        fun fromDomain(device: Device): DeviceEntity = DeviceEntity(
            id = device.id,
            token = device.token,
            deviceModel = device.deviceModel,
            androidVersion = device.androidVersion,
            appVersion = device.appVersion,
            status = device.status.name,
            enrolledAt = device.enrolledAt
        )
    }
}
