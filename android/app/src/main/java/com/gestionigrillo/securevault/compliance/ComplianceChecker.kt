package com.gestionigrillo.securevault.compliance

import android.app.KeyguardManager
import android.content.Context
import android.os.Build
import android.provider.Settings
import androidx.biometric.BiometricManager
import com.gestionigrillo.securevault.domain.model.ComplianceStatus
import java.io.File

/**
 * Checks device compliance for security policy enforcement.
 */
object ComplianceChecker {

    fun check(context: Context): ComplianceStatus {
        return ComplianceStatus(
            isRooted = isRooted(),
            isScreenLockEnabled = isScreenLockEnabled(context),
            isBiometricAvailable = isBiometricAvailable(context),
            isDeviceEncrypted = isDeviceEncrypted(),
            isDebugMode = isDebugMode(context)
        )
    }

    fun isRooted(): Boolean {
        val paths = listOf(
            "/system/app/Superuser.apk",
            "/sbin/su",
            "/system/bin/su",
            "/system/xbin/su",
            "/data/local/xbin/su",
            "/data/local/bin/su",
            "/system/sd/xbin/su"
        )
        return paths.any { File(it).exists() }
    }

    private fun isScreenLockEnabled(context: Context): Boolean {
        val keyguard = context.getSystemService(Context.KEYGUARD_SERVICE) as KeyguardManager
        return keyguard.isDeviceSecure
    }

    private fun isBiometricAvailable(context: Context): Boolean {
        val manager = BiometricManager.from(context)
        return manager.canAuthenticate(BiometricManager.Authenticators.BIOMETRIC_STRONG) ==
                BiometricManager.BIOMETRIC_SUCCESS
    }

    private fun isDeviceEncrypted(): Boolean {
        // Android 7+ is always encrypted with Direct Boot
        return Build.VERSION.SDK_INT >= Build.VERSION_CODES.N
    }

    private fun isDebugMode(context: Context): Boolean {
        return Settings.Secure.getInt(
            context.contentResolver,
            Settings.Global.ADB_ENABLED, 0
        ) != 0
    }
}
