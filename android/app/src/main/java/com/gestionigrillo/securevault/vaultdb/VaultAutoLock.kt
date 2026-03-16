package com.gestionigrillo.securevault.vaultdb

import android.app.Activity
import android.app.Application
import android.os.Bundle
import android.os.Handler
import android.os.Looper

/**
 * Auto-locks the vault after a period of inactivity.
 * Register in Application.onCreate() to track activity lifecycle.
 */
class VaultAutoLock(
    private val timeoutMs: Long = DEFAULT_TIMEOUT_MS,
    private val onLock: () -> Unit
) : Application.ActivityLifecycleCallbacks {

    private val handler = Handler(Looper.getMainLooper())
    private val lockRunnable = Runnable { onLock() }
    private var activeActivities = 0

    override fun onActivityResumed(activity: Activity) {
        activeActivities++
        handler.removeCallbacks(lockRunnable)
    }

    override fun onActivityPaused(activity: Activity) {
        activeActivities--
        if (activeActivities == 0) {
            handler.postDelayed(lockRunnable, timeoutMs)
        }
    }

    override fun onActivityCreated(activity: Activity, savedInstanceState: Bundle?) {}
    override fun onActivityStarted(activity: Activity) {}
    override fun onActivityStopped(activity: Activity) {}
    override fun onActivitySaveInstanceState(activity: Activity, outState: Bundle) {}
    override fun onActivityDestroyed(activity: Activity) {}

    companion object {
        private const val DEFAULT_TIMEOUT_MS = 5 * 60 * 1000L // 5 minutes
    }
}
