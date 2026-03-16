package com.gestionigrillo.securevault.domain.model

data class ComplianceStatus(
    val isRooted: Boolean,
    val isScreenLockEnabled: Boolean,
    val isBiometricAvailable: Boolean,
    val isDeviceEncrypted: Boolean,
    val isDebugMode: Boolean
) {
    val isCompliant: Boolean
        get() = !isRooted && isScreenLockEnabled && isDeviceEncrypted && !isDebugMode
}
