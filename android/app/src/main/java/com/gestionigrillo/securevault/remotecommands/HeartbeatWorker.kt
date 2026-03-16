package com.gestionigrillo.securevault.remotecommands

import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.os.BatteryManager
import androidx.work.BackoffPolicy
import androidx.work.Constraints
import androidx.work.CoroutineWorker
import androidx.work.ExistingPeriodicWorkPolicy
import androidx.work.NetworkType
import androidx.work.PeriodicWorkRequestBuilder
import androidx.work.WorkManager
import androidx.work.WorkerParameters
import com.gestionigrillo.securevault.BuildConfig
import com.gestionigrillo.securevault.SecureVaultApp
import com.gestionigrillo.securevault.compliance.ComplianceChecker
import java.util.concurrent.TimeUnit

class HeartbeatWorker(
    context: Context,
    params: WorkerParameters
) : CoroutineWorker(context, params) {

    override suspend fun doWork(): Result {
        val app = applicationContext as SecureVaultApp
        val device = app.deviceRepository.getEnrolledDevice() ?: return Result.failure()

        val batteryPct = getBatteryLevel()
        val isRooted = ComplianceChecker.isRooted()

        val result = app.deviceRepository.sendHeartbeat(
            deviceId = device.id,
            token = device.token,
            appVersion = BuildConfig.VERSION_NAME,
            batteryPct = batteryPct,
            isRooted = isRooted
        )

        return if (result.isSuccess) Result.success() else Result.retry()
    }

    private fun getBatteryLevel(): Int {
        val intent = applicationContext.registerReceiver(null, IntentFilter(Intent.ACTION_BATTERY_CHANGED))
        val level = intent?.getIntExtra(BatteryManager.EXTRA_LEVEL, -1) ?: -1
        val scale = intent?.getIntExtra(BatteryManager.EXTRA_SCALE, -1) ?: -1
        return if (level >= 0 && scale > 0) (level * 100 / scale) else -1
    }

    companion object {
        private const val WORK_NAME = "heartbeat_worker"

        fun schedule(context: Context) {
            val constraints = Constraints.Builder()
                .setRequiredNetworkType(NetworkType.CONNECTED)
                .build()

            val request = PeriodicWorkRequestBuilder<HeartbeatWorker>(15, TimeUnit.MINUTES)
                .setConstraints(constraints)
                .setBackoffCriteria(BackoffPolicy.EXPONENTIAL, 1, TimeUnit.MINUTES)
                .build()

            WorkManager.getInstance(context)
                .enqueueUniquePeriodicWork(WORK_NAME, ExistingPeriodicWorkPolicy.KEEP, request)
        }

        fun cancel(context: Context) {
            WorkManager.getInstance(context).cancelUniqueWork(WORK_NAME)
        }
    }
}
