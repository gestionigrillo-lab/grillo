package com.gestionigrillo.securevault.data.database.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.gestionigrillo.securevault.data.database.entity.DeviceEntity

@Dao
interface DeviceDao {

    @Query("SELECT * FROM enrolled_device LIMIT 1")
    suspend fun getDevice(): DeviceEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(device: DeviceEntity)

    @Query("DELETE FROM enrolled_device")
    suspend fun deleteAll()
}
