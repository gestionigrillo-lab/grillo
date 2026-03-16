package com.gestionigrillo.securevault.data.database

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase
import com.gestionigrillo.securevault.data.database.dao.DeviceDao
import com.gestionigrillo.securevault.data.database.dao.VaultEntryDao
import com.gestionigrillo.securevault.data.database.entity.DeviceEntity
import com.gestionigrillo.securevault.data.database.entity.VaultEntryEntity

@Database(
    entities = [DeviceEntity::class, VaultEntryEntity::class],
    version = 1,
    exportSchema = false
)
abstract class VaultDatabase : RoomDatabase() {

    abstract fun deviceDao(): DeviceDao
    abstract fun vaultEntryDao(): VaultEntryDao

    companion object {
        @Volatile
        private var INSTANCE: VaultDatabase? = null

        fun getInstance(context: Context): VaultDatabase {
            return INSTANCE ?: synchronized(this) {
                INSTANCE ?: Room.databaseBuilder(
                    context.applicationContext,
                    VaultDatabase::class.java,
                    "secure_vault.db"
                )
                    .fallbackToDestructiveMigration()
                    .build()
                    .also { INSTANCE = it }
            }
        }
    }
}
