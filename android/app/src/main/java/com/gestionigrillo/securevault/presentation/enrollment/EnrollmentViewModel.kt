package com.gestionigrillo.securevault.presentation.enrollment

import android.app.Application
import android.os.Build
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.viewModelScope
import com.gestionigrillo.securevault.BuildConfig
import com.gestionigrillo.securevault.SecureVaultApp
import kotlinx.coroutines.launch

sealed class EnrollmentState {
    data object Idle : EnrollmentState()
    data object Loading : EnrollmentState()
    data class Success(val deviceId: String) : EnrollmentState()
    data class Error(val message: String) : EnrollmentState()
}

class EnrollmentViewModel(application: Application) : AndroidViewModel(application) {

    private val app = application as SecureVaultApp

    private val _enrollState = MutableLiveData<EnrollmentState>(EnrollmentState.Idle)
    val enrollState: LiveData<EnrollmentState> = _enrollState

    fun enroll() {
        _enrollState.value = EnrollmentState.Loading

        viewModelScope.launch {
            val result = app.deviceRepository.enroll(
                deviceModel = Build.MODEL,
                androidVersion = Build.VERSION.RELEASE,
                appVersion = BuildConfig.VERSION_NAME
            )

            result.fold(
                onSuccess = { device ->
                    app.deviceRepository.saveDevice(device)
                    // Initialize master key on enrollment
                    app.keyStoreManager.getOrCreateMasterKey()
                    _enrollState.value = EnrollmentState.Success(device.id)
                },
                onFailure = { e ->
                    _enrollState.value = EnrollmentState.Error(e.message ?: "Enrollment failed")
                }
            )
        }
    }
}
