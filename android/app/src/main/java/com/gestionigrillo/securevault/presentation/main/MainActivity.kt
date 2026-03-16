package com.gestionigrillo.securevault.presentation.main

import android.os.Bundle
import android.view.WindowManager
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.ViewModelProvider
import com.gestionigrillo.securevault.R
import com.gestionigrillo.securevault.databinding.ActivityMainBinding

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var viewModel: MainViewModel

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Prevent screenshots
        window.setFlags(
            WindowManager.LayoutParams.FLAG_SECURE,
            WindowManager.LayoutParams.FLAG_SECURE
        )

        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        viewModel = ViewModelProvider(this)[MainViewModel::class.java]

        observeState()
        setupListeners()
    }

    private fun observeState() {
        viewModel.vaultEntries.observe(this) { entries ->
            binding.entryCount.text = "Elementi nel vault: ${entries.size}"
        }

        viewModel.deviceStatus.observe(this) { status ->
            binding.deviceStatus.text = "Stato: $status"
        }
    }

    private fun setupListeners() {
        binding.btnWipeVault.setOnClickListener {
            AlertDialog.Builder(this)
                .setTitle("Wipe Vault")
                .setMessage(getString(R.string.wipe_confirm))
                .setPositiveButton("Conferma") { _, _ ->
                    viewModel.wipeVault()
                }
                .setNegativeButton("Annulla", null)
                .show()
        }

        binding.btnLockVault.setOnClickListener {
            viewModel.lockVault()
            finish()
        }
    }
}
