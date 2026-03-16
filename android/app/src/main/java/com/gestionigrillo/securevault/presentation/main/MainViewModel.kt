package com.gestionigrillo.securevault.presentation.main

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.viewModelScope
import com.gestionigrillo.securevault.SecureVaultApp
import com.gestionigrillo.securevault.domain.model.VaultEntry
import kotlinx.coroutines.launch

class MainViewModel(application: Application) : AndroidViewModel(application) {

    private val app = application as SecureVaultApp

    private val _vaultEntries = MutableLiveData<List<VaultEntry>>(emptyList())
    val vaultEntries: LiveData<List<VaultEntry>> = _vaultEntries

    private val _deviceStatus = MutableLiveData<String>()
    val deviceStatus: LiveData<String> = _deviceStatus

    init {
        loadEntries()
        loadDeviceStatus()
    }

    private fun loadEntries() {
        viewModelScope.launch {
            app.vaultRepository.getAllEntries().fold(
                onSuccess = { _vaultEntries.value = it },
                onFailure = { _vaultEntries.value = emptyList() }
            )
        }
    }

    private fun loadDeviceStatus() {
        viewModelScope.launch {
            val device = app.deviceRepository.getEnrolledDevice()
            _deviceStatus.value = device?.status?.name ?: "NOT ENROLLED"
        }
    }

    fun wipeVault() {
        viewModelScope.launch {
            app.vaultRepository.wipeAll()
            app.keyStoreManager.deleteVaultKey()
            _vaultEntries.value = emptyList()
        }
    }

    fun lockVault() {
        // Simply closing the activity returns to VaultAuthActivity
    }
}
