package com.gestionigrillo.securevault.data.network

import com.gestionigrillo.securevault.data.network.dto.AckRequest
import com.gestionigrillo.securevault.data.network.dto.CommandsResponse
import com.gestionigrillo.securevault.data.network.dto.EnrollRequest
import com.gestionigrillo.securevault.data.network.dto.EnrollResponse
import com.gestionigrillo.securevault.data.network.dto.HeartbeatRequest
import com.gestionigrillo.securevault.data.network.dto.HeartbeatResponse
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.POST
import retrofit2.http.Path

interface ApiService {

    @POST("device/enroll")
    suspend fun enroll(@Body request: EnrollRequest): Response<EnrollResponse>

    @POST("device/heartbeat")
    suspend fun heartbeat(
        @Header("Authorization") token: String,
        @Body request: HeartbeatRequest
    ): Response<HeartbeatResponse>

    @GET("device/{deviceId}/commands")
    suspend fun getCommands(
        @Header("Authorization") token: String,
        @Path("deviceId") deviceId: String
    ): Response<CommandsResponse>

    @POST("device/commands/{commandId}/ack")
    suspend fun ackCommand(
        @Header("Authorization") token: String,
        @Path("commandId") commandId: String,
        @Body request: AckRequest
    ): Response<Unit>
}
