package com.gestionigrillo.securevault.presentation.enrollment

import android.content.Intent
import android.os.Build
import android.os.Bundle
import android.view.WindowManager
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.ViewModelProvider
import com.gestionigrillo.securevault.databinding.ActivityEnrollmentBinding
import com.gestionigrillo.securevault.presentation.vaultauth.VaultAuthActivity

class EnrollmentActivity : AppCompatActivity() {

    private lateinit var binding: ActivityEnrollmentBinding
    private lateinit var viewModel: EnrollmentViewModel

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        window.setFlags(
            WindowManager.LayoutParams.FLAG_SECURE,
            WindowManager.LayoutParams.FLAG_SECURE
        )

        binding = ActivityEnrollmentBinding.inflate(layoutInflater)
        setContentView(binding.root)

        viewModel = ViewModelProvider(this)[EnrollmentViewModel::class.java]

        binding.deviceModelText.text = Build.MODEL
        binding.androidVersionText.text = Build.VERSION.RELEASE

        observeState()

        binding.btnEnroll.setOnClickListener {
            viewModel.enroll()
        }
    }

    private fun observeState() {
        viewModel.enrollState.observe(this) { state ->
            when (state) {
                is EnrollmentState.Idle -> {
                    binding.btnEnroll.isEnabled = true
                }
                is EnrollmentState.Loading -> {
                    binding.btnEnroll.isEnabled = false
                }
                is EnrollmentState.Success -> {
                    Toast.makeText(this, "Dispositivo registrato: ${state.deviceId}", Toast.LENGTH_SHORT).show()
                    startActivity(Intent(this, VaultAuthActivity::class.java))
                    finish()
                }
                is EnrollmentState.Error -> {
                    binding.btnEnroll.isEnabled = true
                    Toast.makeText(this, state.message, Toast.LENGTH_LONG).show()
                }
            }
        }
    }
}
