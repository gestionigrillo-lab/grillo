package com.gestionigrillo.securevault.presentation.vaultauth

import android.content.Intent
import android.os.Bundle
import android.view.WindowManager
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.biometric.BiometricManager
import androidx.biometric.BiometricPrompt
import androidx.core.content.ContextCompat
import androidx.lifecycle.ViewModelProvider
import com.gestionigrillo.securevault.R
import com.gestionigrillo.securevault.databinding.ActivityVaultAuthBinding
import com.gestionigrillo.securevault.presentation.enrollment.EnrollmentActivity
import com.gestionigrillo.securevault.presentation.main.MainActivity

class VaultAuthActivity : AppCompatActivity() {

    private lateinit var binding: ActivityVaultAuthBinding
    private lateinit var viewModel: VaultAuthViewModel

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Prevent screenshots on sensitive screens
        window.setFlags(
            WindowManager.LayoutParams.FLAG_SECURE,
            WindowManager.LayoutParams.FLAG_SECURE
        )

        binding = ActivityVaultAuthBinding.inflate(layoutInflater)
        setContentView(binding.root)

        viewModel = ViewModelProvider(this)[VaultAuthViewModel::class.java]

        observeState()
        setupListeners()
    }

    private fun observeState() {
        viewModel.authState.observe(this) { state ->
            when (state) {
                is VaultAuthState.NeedsEnrollment -> {
                    startActivity(Intent(this, EnrollmentActivity::class.java))
                    finish()
                }
                is VaultAuthState.NeedsSetup -> {
                    binding.pinInput.isEnabled = true
                    binding.btnConfirm.text = getString(R.string.confirm_pin)
                }
                is VaultAuthState.Locked -> {
                    binding.pinInput.isEnabled = true
                    binding.btnConfirm.text = getString(R.string.enter_pin)
                    tryBiometric()
                }
                is VaultAuthState.Unlocked -> {
                    startActivity(Intent(this, MainActivity::class.java))
                    finish()
                }
                is VaultAuthState.Error -> {
                    Toast.makeText(this, state.message, Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun setupListeners() {
        binding.btnConfirm.setOnClickListener {
            val pin = binding.pinInput.text.toString()
            if (pin.length < 6) {
                Toast.makeText(this, "PIN deve essere di almeno 6 cifre", Toast.LENGTH_SHORT).show()
                return@setOnClickListener
            }
            viewModel.submitPin(pin)
        }

        binding.btnBiometric.setOnClickListener {
            tryBiometric()
        }
    }

    private fun tryBiometric() {
        val biometricManager = BiometricManager.from(this)
        if (biometricManager.canAuthenticate(BiometricManager.Authenticators.BIOMETRIC_STRONG) !=
            BiometricManager.BIOMETRIC_SUCCESS
        ) {
            return
        }

        val prompt = BiometricPrompt(this, ContextCompat.getMainExecutor(this),
            object : BiometricPrompt.AuthenticationCallback() {
                override fun onAuthenticationSucceeded(result: BiometricPrompt.AuthenticationResult) {
                    viewModel.onBiometricSuccess()
                }

                override fun onAuthenticationError(errorCode: Int, errString: CharSequence) {
                    // User cancelled or biometric not available
                }
            })

        val promptInfo = BiometricPrompt.PromptInfo.Builder()
            .setTitle(getString(R.string.biometric_title))
            .setSubtitle(getString(R.string.biometric_subtitle))
            .setNegativeButtonText(getString(R.string.enter_pin))
            .build()

        prompt.authenticate(promptInfo)
    }
}
