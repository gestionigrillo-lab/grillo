package com.gestionigrillo.securevault.remotecommands

import android.content.Context
import androidx.work.BackoffPolicy
import androidx.work.Constraints
import androidx.work.CoroutineWorker
import androidx.work.ExistingPeriodicWorkPolicy
import androidx.work.NetworkType
import androidx.work.PeriodicWorkRequestBuilder
import androidx.work.WorkManager
import androidx.work.WorkerParameters
import com.gestionigrillo.securevault.SecureVaultApp
import com.gestionigrillo.securevault.domain.model.CommandType
import com.gestionigrillo.securevault.wipe.VaultWiper
import java.util.concurrent.TimeUnit

class CommandPollingWorker(
    context: Context,
    params: WorkerParameters
) : CoroutineWorker(context, params) {

    override suspend fun doWork(): Result {
        val app = applicationContext as SecureVaultApp
        val device = app.deviceRepository.getEnrolledDevice() ?: return Result.failure()

        val result = app.deviceRepository.fetchPendingCommands(device.id, device.token)

        result.fold(
            onSuccess = { commands ->
                for (cmd in commands) {
                    val success = executeCommand(app, cmd.type)
                    app.deviceRepository.acknowledgeCommand(
                        commandId = cmd.id,
                        token = device.token,
                        success = success,
                        message = if (success) "executed" else "failed"
                    )
                }
            },
            onFailure = { return Result.retry() }
        )

        return Result.success()
    }

    private suspend fun executeCommand(app: SecureVaultApp, type: CommandType): Boolean {
        return when (type) {
            CommandType.LOCK -> {
                // Lock the vault by clearing session state
                VaultWiper.lockVault(applicationContext)
                true
            }
            CommandType.WIPE -> {
                VaultWiper.wipeVault(applicationContext, app.vaultRepository, app.keyStoreManager)
                true
            }
            CommandType.REVOKE -> {
                VaultWiper.wipeVault(applicationContext, app.vaultRepository, app.keyStoreManager)
                app.deviceRepository.clearDevice()
                HeartbeatWorker.cancel(applicationContext)
                cancel(applicationContext)
                true
            }
        }
    }

    companion object {
        private const val WORK_NAME = "command_polling_worker"

        fun schedule(context: Context) {
            val constraints = Constraints.Builder()
                .setRequiredNetworkType(NetworkType.CONNECTED)
                .build()

            val request = PeriodicWorkRequestBuilder<CommandPollingWorker>(15, TimeUnit.MINUTES)
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
