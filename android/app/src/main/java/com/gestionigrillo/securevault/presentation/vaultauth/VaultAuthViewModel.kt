package com.gestionigrillo.securevault.presentation.vaultauth

import android.app.Application
import android.content.Context
import android.util.Base64
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.viewModelScope
import com.gestionigrillo.securevault.SecureVaultApp
import com.gestionigrillo.securevault.vaultcrypto.PinHasher
import kotlinx.coroutines.launch

sealed class VaultAuthState {
    data object NeedsEnrollment : VaultAuthState()
    data object NeedsSetup : VaultAuthState()
    data object Locked : VaultAuthState()
    data object Unlocked : VaultAuthState()
    data class Error(val message: String) : VaultAuthState()
}

class VaultAuthViewModel(application: Application) : AndroidViewModel(application) {

    private val app = application as SecureVaultApp
    private val prefs = application.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)

    private val _authState = MutableLiveData<VaultAuthState>()
    val authState: LiveData<VaultAuthState> = _authState

    init {
        checkState()
    }

    private fun checkState() {
        viewModelScope.launch {
            val device = app.deviceRepository.getEnrolledDevice()
            if (device == null) {
                _authState.value = VaultAuthState.NeedsEnrollment
                return@launch
            }

            if (!hasPinSet()) {
                _authState.value = VaultAuthState.NeedsSetup
            } else {
                _authState.value = VaultAuthState.Locked
            }
        }
    }

    fun submitPin(pin: String) {
        viewModelScope.launch {
            if (!hasPinSet()) {
                setupPin(pin)
            } else {
                verifyPin(pin)
            }
        }
    }

    private fun setupPin(pin: String) {
        val hashed = PinHasher.hash(pin)
        prefs.edit()
            .putString(KEY_PIN_HASH, Base64.encodeToString(hashed.hash, Base64.NO_WRAP))
            .putString(KEY_PIN_SALT, Base64.encodeToString(hashed.salt, Base64.NO_WRAP))
            .apply()

        // Ensure vault key is created
        app.keyStoreManager.getOrCreateVaultKey()
        _authState.value = VaultAuthState.Unlocked
    }

    private fun verifyPin(pin: String) {
        val hashB64 = prefs.getString(KEY_PIN_HASH, null) ?: return
        val saltB64 = prefs.getString(KEY_PIN_SALT, null) ?: return

        val stored = PinHasher.HashedPin(
            hash = Base64.decode(hashB64, Base64.NO_WRAP),
            salt = Base64.decode(saltB64, Base64.NO_WRAP)
        )

        if (PinHasher.verify(pin, stored)) {
            _authState.value = VaultAuthState.Unlocked
        } else {
            _authState.value = VaultAuthState.Error("PIN errato")
        }
    }

    fun onBiometricSuccess() {
        if (hasPinSet()) {
            _authState.value = VaultAuthState.Unlocked
        }
    }

    private fun hasPinSet(): Boolean = prefs.contains(KEY_PIN_HASH)

    companion object {
        private const val PREFS_NAME = "vault_auth_prefs"
        private const val KEY_PIN_HASH = "pin_hash"
        private const val KEY_PIN_SALT = "pin_salt"
    }
}
