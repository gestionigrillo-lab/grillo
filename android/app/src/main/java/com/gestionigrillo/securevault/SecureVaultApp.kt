package com.gestionigrillo.securevault

import android.app.Application
import androidx.work.Configuration
import com.gestionigrillo.securevault.data.database.VaultDatabase
import com.gestionigrillo.securevault.data.network.ApiClient
import com.gestionigrillo.securevault.data.repository.DeviceRepositoryImpl
import com.gestionigrillo.securevault.data.repository.VaultRepositoryImpl
import com.gestionigrillo.securevault.domain.repository.DeviceRepository
import com.gestionigrillo.securevault.domain.repository.VaultRepository
import com.gestionigrillo.securevault.vaultcrypto.KeyStoreManager
import com.gestionigrillo.securevault.vaultcrypto.VaultCryptoEngine

class SecureVaultApp : Application(), Configuration.Provider {

    lateinit var keyStoreManager: KeyStoreManager
        private set
    lateinit var cryptoEngine: VaultCryptoEngine
        private set
    lateinit var database: VaultDatabase
        private set
    lateinit var deviceRepository: DeviceRepository
        private set
    lateinit var vaultRepository: VaultRepository
        private set

    override fun onCreate() {
        super.onCreate()
        instance = this

        keyStoreManager = KeyStoreManager()
        cryptoEngine = VaultCryptoEngine(keyStoreManager)
        database = VaultDatabase.getInstance(this)

        val apiService = ApiClient.create(this)
        deviceRepository = DeviceRepositoryImpl(apiService, database.deviceDao())
        vaultRepository = VaultRepositoryImpl(database.vaultEntryDao(), cryptoEngine)
    }

    override val workManagerConfiguration: Configuration
        get() = Configuration.Builder()
            .setMinimumLoggingLevel(android.util.Log.INFO)
            .build()

    companion object {
        lateinit var instance: SecureVaultApp
            private set
    }
}
